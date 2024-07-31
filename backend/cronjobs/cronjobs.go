package cronjobs

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/robfig/cron/v3"
	"github.com/seriousm4x/upsnap/logger"
	"github.com/seriousm4x/upsnap/networking"
)

var (
	PingRunning         = false
	WakeShutdownRunning = false
	CronPing            = cron.New()
	CronWakeShutdown    = cron.New()
)

func SetPingJobs(app *pocketbase.PocketBase) {
	// remove existing jobs
	for _, job := range CronPing.Entries() {
		CronPing.Remove(job.ID)
	}

	settingsPrivateRecords, err := app.Dao().FindRecordsByExpr("settings_private")
	if err != nil {
		logger.Error.Println(err)
	}

	CronPing.AddFunc(settingsPrivateRecords[0].GetString("interval"), func() {
		// skip cron if no realtime clients connected and lazy_ping is turned on
		realtimeClients := len(app.SubscriptionsBroker().Clients())
		if realtimeClients == 0 && settingsPrivateRecords[0].GetBool("lazy_ping") {
			return
		}

		devices, err := app.Dao().FindRecordsByExpr("devices")
		if err != nil {
			logger.Error.Println(err)
			return
		}

		// expand ports field
		expandFetchFunc := func(c *models.Collection, ids []string) ([]*models.Record, error) {
			return app.Dao().FindRecordsByIds(c.Id, ids, nil)
		}
		merr := app.Dao().ExpandRecords(devices, []string{"ports"}, expandFetchFunc)
		if len(merr) > 0 {
			return
		}

		for _, device := range devices {
			// ping device
			go func(d *models.Record) {
				status := d.GetString("status")
				if status == "pending" {
					return
				}
				if networking.PingDevice(d) {
					if status == "online" {
						return
					}
					d.Set("status", "online")
					if err := app.Dao().SaveRecord(d); err != nil {
						logger.Error.Println("Failed to save record:", err)
					}
				} else {
					if status == "offline" {
						return
					}
					d.Set("status", "offline")
					if err := app.Dao().SaveRecord(d); err != nil {
						logger.Error.Println("Failed to save record:", err)
					}
				}
			}(device)

			// ping ports
			go func(d *models.Record) {
				ports, err := app.Dao().FindRecordsByIds("ports", d.GetStringSlice("ports"))
				if err != nil {
					logger.Error.Println(err)
				}
				for _, port := range ports {
					isUp := networking.CheckPort(d.GetString("ip"), port.GetString("number"))
					if isUp != port.GetBool("status") {
						port.Set("status", isUp)
						if err := app.Dao().SaveRecord(port); err != nil {
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

	devices, err := app.Dao().FindRecordsByExpr("devices")
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
				logger.Debug.Printf("[CRON1 \"%s\"]: cron func started", dev.GetString("name"))
				dev, err = app.Dao().FindRecordById("devices", dev.Id)
				if err != nil {
					logger.Error.Println(err)
					return
				}
				logger.Debug.Printf("[CRON2 \"%s\"]: got record from db", dev.GetString("name"))

				status := dev.GetString("status")
				logger.Debug.Printf("[CRON3 \"%s\"]: status is %s", dev.GetString("name"), status)
				if status == "pending" || status == "online" {
					logger.Debug.Printf("[CRON3.5 \"%s\"]: skipping run because already pending", dev.GetString("name"))
					return
				}
				dev.Set("status", "pending")
				if err := app.Dao().SaveRecord(dev); err != nil {
					logger.Error.Println("Failed to save record:", err)
					return
				}
				logger.Debug.Printf("[CRON4 \"%s\"]: saved status pending", dev.GetString("name"))
				if err := networking.WakeDevice(dev); err != nil {
					logger.Error.Println(err)
					dev.Set("status", "offline")
				} else {
					dev.Set("status", "online")
				}
				logger.Debug.Printf("[CRON5 \"%s\"]: wake device done", dev.GetString("name"))
				if err := app.Dao().SaveRecord(dev); err != nil {
					logger.Error.Println("Failed to save record:", err)
				}
				logger.Debug.Printf("[CRON6 \"%s\"]: saved device", dev.GetString("name"))
			})
			if err != nil {
				logger.Error.Println(err)
			}
		}

		if shutdown_cron_enabled && shutdown_cron != "" {
			_, err := CronWakeShutdown.AddFunc(shutdown_cron, func() {
				dev, err = app.Dao().FindRecordById("devices", dev.Id)
				if err != nil {
					logger.Error.Println(err)
					return
				}

				status := dev.GetString("status")
				if status == "pending" || status == "offline" {
					return
				}
				dev.Set("status", "pending")
				if err := app.Dao().SaveRecord(dev); err != nil {
					logger.Error.Println("Failed to save record:", err)
				}
				if err := networking.ShutdownDevice(dev); err != nil {
					logger.Error.Println(err)
					dev.Set("status", "online")
				} else {
					dev.Set("status", "offline")
				}
				if err := app.Dao().SaveRecord(dev); err != nil {
					logger.Error.Println("Failed to save record:", err)
				}
			})
			if err != nil {
				logger.Error.Println(err)
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
		ctx := CronWakeShutdown.Stop()
		<-ctx.Done()
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
		ctx := CronPing.Stop()
		<-ctx.Done()
	}
	PingRunning = false
}

func StopAll() {
	StopPing()
	StopWakeShutdown()
}
