package cronjobs

import (
	"net"
	"sync/atomic"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/robfig/cron/v3"
	"github.com/seriousm4x/upsnap/logger"
	"github.com/seriousm4x/upsnap/networking"
)

var (
	PingRunning         = false
	WakeShutdownRunning = false
	CronPing            = cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
	)))
	CronWakeShutdown = cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
	)))
)

func SetPingJobs(app *pocketbase.PocketBase) {
	// remove existing jobs
	for _, job := range CronPing.Entries() {
		CronPing.Remove(job.ID)
	}

	settingsPrivateRecords, err := app.FindAllRecords("settings_private")
	if err != nil {
		logger.Error.Println(err)
	}

	CronPing.AddFunc(settingsPrivateRecords[0].GetString("interval"), func() {
		// skip cron if no realtime clients connected and lazy_ping is turned on
		realtimeClients := len(app.SubscriptionsBroker().Clients())
		if realtimeClients == 0 && settingsPrivateRecords[0].GetBool("lazy_ping") {
			return
		}

		devices, err := app.FindAllRecords("devices")
		if err != nil {
			logger.Error.Println(err)
			return
		}

		// expand ports field
		expandFetchFunc := func(c *core.Collection, ids []string) ([]*core.Record, error) {
			return app.FindRecordsByIds(c.Id, ids, nil)
		}
		merr := app.ExpandRecords(devices, []string{"ports"}, expandFetchFunc)
		if len(merr) > 0 {
			return
		}

		for _, device := range devices {
			// ping device
			go func(d *core.Record) {
				status := d.GetString("status")
				if status == "pending" {
					return
				}
				isUp, err := networking.PingDevice(d)
				if err != nil {
					logger.Error.Println(err)
				}
				if isUp {
					if status == "online" {
						return
					}
					d.Set("status", "online")
					if err := app.Save(d); err != nil {
						logger.Error.Println("Failed to save record:", err)
					}
				} else {
					if status == "offline" {
						return
					}
					d.Set("status", "offline")
					if err := app.Save(d); err != nil {
						logger.Error.Println("Failed to save record:", err)
					}
				}
			}(device)

			// ping ports
			go func(d *core.Record) {
				ports, err := app.FindRecordsByIds("ports", d.GetStringSlice("ports"))
				if err != nil {
					logger.Error.Println(err)
				}
				for _, port := range ports {
					isUp, err := networking.CheckPort(d.GetString("ip"), port.GetString("number"))
					if err != nil {
						logger.Error.Println("Failed to check port:", err)
					}
					if isUp != port.GetBool("status") {
						port.Set("status", isUp)
						if err := app.Save(port); err != nil {
							logger.Error.Println("Failed to save record:", err)
						}
					}
				}
			}(device)
		}
	})

	// update ip addresses of opted-in devices from a periodic arp scan
	trackIpInterval := settingsPrivateRecords[0].GetString("track_ip_interval")
	if trackIpInterval != "" {
		if _, err := CronPing.AddFunc(trackIpInterval, func() {
			trackDeviceIPs(app)
		}); err != nil {
			logger.Error.Println("Failed to add ip tracking cronjob:", err)
		}
	}
}

var trackIpRunning atomic.Bool

// trackDeviceIPs scans the local subnets of devices with ip tracking enabled
// and updates their ip address if their mac address is found at a different one.
func trackDeviceIPs(app *pocketbase.PocketBase) {
	// skip if the previous scan is still running
	if !trackIpRunning.CompareAndSwap(false, true) {
		return
	}
	defer trackIpRunning.Store(false)

	devices, err := app.FindRecordsByFilter("devices", "track_ip = true", "", 0, 0)
	if err != nil {
		logger.Error.Println(err)
		return
	}

	// group devices by subnet, skipping subnets not directly attached to the host
	subnets := make(map[string][]*core.Record)
	for _, device := range devices {
		subnet, err := networking.DeviceSubnet(device.GetString("ip"), device.GetString("netmask"))
		if err != nil {
			logger.Error.Println("Ip tracking for", device.GetString("name")+":", err)
			continue
		}
		if ones, bits := subnet.Mask.Size(); ones == bits || ones < 16 {
			// a /32 has nothing to scan, anything larger than a /16 takes too long
			continue
		}
		if !networking.IsLocalSubnet(subnet) {
			continue
		}
		subnets[subnet.String()] = append(subnets[subnet.String()], device)
	}

	for cidr, subnetDevices := range subnets {
		scan, err := networking.NmapScan(cidr)
		if err != nil {
			logger.Error.Println("Ip tracking scan for", cidr+":", err)
			continue
		}

		// map mac addresses to ips
		macToIp := make(map[string]string)
		for _, host := range scan.Host {
			var hostIp, hostMac string
			for _, addr := range host.Address {
				if addr.Addrtype == "ipv4" {
					hostIp = addr.Addr
				} else if addr.Addrtype == "mac" {
					hostMac = addr.Addr
				}
			}
			if parsedMac, err := net.ParseMAC(hostMac); hostIp != "" && err == nil {
				macToIp[parsedMac.String()] = hostIp
			}
		}

		for _, device := range subnetDevices {
			parsedMac, err := net.ParseMAC(device.GetString("mac"))
			if err != nil {
				continue
			}
			newIp, ok := macToIp[parsedMac.String()]
			if !ok || newIp == device.GetString("ip") {
				continue
			}
			logger.Info.Println("Ip tracking: updating", device.GetString("name"), "from", device.GetString("ip"), "to", newIp)
			device.Set("ip", newIp)
			// only write the changed ip field to avoid clobbering concurrent
			// status updates from the ping and wake/shutdown cronjobs
			device.IgnoreUnchangedFields(true)
			if err := app.Save(device); err != nil {
				logger.Error.Println("Failed to save record:", err)
			}
		}
	}
}

func SetWakeShutdownJobs(app *pocketbase.PocketBase) {
	// remove existing jobs
	for _, job := range CronWakeShutdown.Entries() {
		CronWakeShutdown.Remove(job.ID)
	}

	devices, err := app.FindAllRecords("devices")
	if err != nil {
		logger.Error.Println(err)
		return
	}
	for _, dev := range devices {
		wake_cron := dev.GetString("wake_cron")
		wake_cron_enabled := dev.GetBool("wake_cron_enabled")
		shutdown_cron := dev.GetString("shutdown_cron")
		shutdown_cron_enabled := dev.GetBool("shutdown_cron_enabled")

		if wake_cron_enabled && wake_cron != "" {
			_, err := CronWakeShutdown.AddFunc(wake_cron, func() {
				d, err := app.FindRecordById("devices", dev.Id)
				if err != nil {
					logger.Error.Println(err)
					return
				}
				if d.GetString("status") == "pending" {
					return
				}
				isOnline, err := networking.PingDevice(d)
				if err != nil {
					logger.Error.Println(err)
					return
				}
				if isOnline {
					return
				}
				d.Set("status", "pending")
				if err := app.Save(d); err != nil {
					logger.Error.Println("Failed to save record:", err)
					return
				}
				if err := networking.WakeDevice(d); err != nil {
					logger.Error.Println(err)
					d.Set("status", "offline")
				} else {
					d.Set("status", "online")
				}
				if err := app.Save(d); err != nil {
					logger.Error.Println("Failed to save record:", err)
				}
			})
			if err != nil {
				logger.Error.Printf("device %s: %+v", dev.GetString("name"), err)
			}
		}

		if shutdown_cron_enabled && shutdown_cron != "" {
			_, err := CronWakeShutdown.AddFunc(shutdown_cron, func() {
				d, err := app.FindRecordById("devices", dev.Id)
				if err != nil {
					logger.Error.Println(err)
					return
				}
				if d.GetString("status") == "pending" {
					return
				}
				isOnline, err := networking.PingDevice(d)
				if err != nil {
					logger.Error.Println(err)
					return
				}
				if !isOnline {
					return
				}
				status := d.GetString("status")
				if status != "online" {
					return
				}
				d.Set("status", "pending")
				if err := app.Save(d); err != nil {
					logger.Error.Println("Failed to save record:", err)
				}
				if err := networking.ShutdownDevice(d); err != nil {
					logger.Error.Println(err)
					d.Set("status", "online")
				} else {
					d.Set("status", "offline")
				}
				if err := app.Save(d); err != nil {
					logger.Error.Println("Failed to save record:", err)
				}
			})
			if err != nil {
				logger.Error.Printf("device %s: %+v", dev.GetString("name"), err)
			}
		}
	}
}

func StartWakeShutdown() {
	WakeShutdownRunning = true
	go CronWakeShutdown.Run()

}

func StopWakeShutdown() {
	if WakeShutdownRunning {
		logger.Info.Println("Stopping wake/shutdown cronjob")
		CronWakeShutdown.Stop()
	}
	WakeShutdownRunning = false
}

func StartPing() {
	PingRunning = true
	go CronPing.Run()
}

func StopPing() {
	if PingRunning {
		logger.Info.Println("Stopping wake/shutdown cronjob")
		CronPing.Stop()
	}
	PingRunning = false
}

func StopAll() {
	StopPing()
	StopWakeShutdown()
}
