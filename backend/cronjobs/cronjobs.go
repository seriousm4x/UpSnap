package cronjobs

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/robfig/cron/v3"
	"github.com/seriousm4x/upsnap/iptracking"
	"github.com/seriousm4x/upsnap/logging"
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

func SetPingJobs(app core.App) {
	log := logging.Logger(app)
	// remove existing jobs
	for _, job := range CronPing.Entries() {
		CronPing.Remove(job.ID)
	}

	settingsPrivateRecords, err := app.FindAllRecords("settings_private")
	if err != nil {
		log.Error("Failed to load private settings", "error", err)
		return
	}

	CronPing.AddFunc(settingsPrivateRecords[0].GetString("interval"), func() {
		// skip cron if no realtime clients connected and lazy_ping is turned on
		realtimeClients := len(app.SubscriptionsBroker().Clients())
		if realtimeClients == 0 && settingsPrivateRecords[0].GetBool("lazy_ping") {
			return
		}

		devices, err := app.FindAllRecords("devices")
		if err != nil {
			log.Error("Failed to load devices for ping job", "error", err)
			return
		}

		// expand ports field
		expandFetchFunc := func(c *core.Collection, ids []string) ([]*core.Record, error) {
			return app.FindRecordsByIds(c.Id, ids, nil)
		}
		merr := app.ExpandRecords(devices, []string{"ports"}, expandFetchFunc)
		if len(merr) > 0 {
			log.Error("Failed to expand device ports", "errors", merr)
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
					log.Error("Failed to ping device", "device", d.GetString("name"), "error", err)
				}
				if isUp {
					if status == "online" {
						return
					}
					d.Set("status", "online")
					if err := app.Save(d); err != nil {
						log.Error("Failed to save device status", "device", d.GetString("name"), "error", err)
					} else {
						log.Info("Device turned on", "device", d.GetString("name"))
					}
				} else {
					if status == "offline" {
						return
					}
					d.Set("status", "offline")
					if err := app.Save(d); err != nil {
						log.Error("Failed to save device status", "device", d.GetString("name"), "error", err)
					} else {
						log.Info("Device turned off", "device", d.GetString("name"))
					}
				}
			}(device)

			// ping ports
			go func(d *core.Record) {
				ports, err := app.FindRecordsByIds("ports", d.GetStringSlice("ports"))
				if err != nil {
					log.Error("Failed to load device ports", "device", d.GetString("name"), "error", err)
					return
				}
				deviceIP, resolveErr := networking.ResolveToIPAddr(d.GetString("ip"))
				if resolveErr != nil {
					log.Error("Failed to resolve device IP", "device", d.GetString("name"), "error", resolveErr)
				}
				for _, port := range ports {
					isUp := false
					// No IP found means the device is down.
					if resolveErr == nil && deviceIP != "" {
						var checkErr error
						isUp, checkErr = networking.CheckPort(deviceIP, port.GetString("number"))
						if checkErr != nil {
							log.Error("Failed to check port", "device", d.GetString("name"), "port", port.GetString("number"), "error", checkErr)
						}
					}
					if isUp != port.GetBool("status") {
						port.Set("status", isUp)
						if err := app.Save(port); err != nil {
							log.Error("Failed to save port status", "device", d.GetString("name"), "port", port.GetString("number"), "error", err)
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
			// pause scans if no realtime clients connected and lazy_ping is
			// turned on; a catch-up sweep runs when the next client connects
			realtimeClients := len(app.SubscriptionsBroker().Clients())
			iptracking.PeriodicSweep(app, realtimeClients == 0 && settingsPrivateRecords[0].GetBool("lazy_ping"))
		}); err != nil {
			log.Error("Failed to add IP tracking cron job", "error", err)
		}
	}
}

func SetWakeShutdownJobs(app core.App) {
	log := logging.Logger(app)
	// remove existing jobs
	for _, job := range CronWakeShutdown.Entries() {
		CronWakeShutdown.Remove(job.ID)
	}

	devices, err := app.FindAllRecords("devices")
	if err != nil {
		log.Error("Failed to load devices for wake and shutdown jobs", "error", err)
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
					log.Error("Failed to load device for wake job", "device", dev.GetString("name"), "error", err)
					return
				}
				if d.GetString("status") == "pending" {
					return
				}
				isOnline, err := networking.PingDevice(d)
				if err != nil {
					log.Error("Failed to ping device before wake", "device", d.GetString("name"), "error", err)
					return
				}
				if isOnline {
					return
				}
				d.Set("status", "pending")
				if err := app.Save(d); err != nil {
					log.Error("Failed to save pending device status", "device", d.GetString("name"), "error", err)
					return
				}
				if err := networking.WakeDevice(app, d, "cron"); err != nil {
					log.Error("Wake job failed", "device", d.GetString("name"), "error", err)
					d.Set("status", "offline")
				} else {
					d.Set("status", "online")
				}
				if err := app.Save(d); err != nil {
					log.Error("Failed to save device status after wake", "device", d.GetString("name"), "error", err)
				}
			})
			if err != nil {
				log.Error("Failed to add wake cron job", "device", dev.GetString("name"), "error", err)
			}
		}

		if shutdown_cron_enabled && shutdown_cron != "" {
			_, err := CronWakeShutdown.AddFunc(shutdown_cron, func() {
				d, err := app.FindRecordById("devices", dev.Id)
				if err != nil {
					log.Error("Failed to load device for shutdown job", "device", dev.GetString("name"), "error", err)
					return
				}
				if d.GetString("status") == "pending" {
					return
				}
				isOnline, err := networking.PingDevice(d)
				if err != nil {
					log.Error("Failed to ping device before shutdown", "device", d.GetString("name"), "error", err)
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
					log.Error("Failed to save pending device status", "device", d.GetString("name"), "error", err)
				}
				if err := networking.ShutdownDevice(app, d, "cron"); err != nil {
					log.Error("Shutdown job failed", "device", d.GetString("name"), "error", err)
					d.Set("status", "online")
				} else {
					d.Set("status", "offline")
				}
				if err := app.Save(d); err != nil {
					log.Error("Failed to save device status after shutdown", "device", d.GetString("name"), "error", err)
				}
			})
			if err != nil {
				log.Error("Failed to add shutdown cron job", "device", dev.GetString("name"), "error", err)
			}
		}
	}
}

func StartWakeShutdown() {
	WakeShutdownRunning = true
	go CronWakeShutdown.Run()

}

func StopWakeShutdown(app core.App) {
	if WakeShutdownRunning {
		logging.Logger(app).Info("Stopping wake and shutdown cron jobs")
		CronWakeShutdown.Stop()
	}
	WakeShutdownRunning = false
}

func StartPing() {
	PingRunning = true
	go CronPing.Run()
}

func StopPing(app core.App) {
	if PingRunning {
		logging.Logger(app).Info("Stopping ping cron jobs")
		CronPing.Stop()
	}
	PingRunning = false
}

func StopAll(app core.App) {
	StopPing(app)
	StopWakeShutdown(app)
}
