package cronjobs

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/robfig/cron/v3"
	"github.com/seriousm4x/upsnap/logger"
	"github.com/seriousm4x/upsnap/networking"
)

var (
	PingRunning         = false
	WakeShutdownRunning = false
	CronPing            = cron.New(cron.WithSeconds())
	CronWakeShutdown    = cron.New(cron.WithSeconds())
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
				isOnline, err := networking.PingDevice(dev)
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
				isOnline, err := networking.PingDevice(dev)
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
