package cronjobs

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/robfig/cron/v3"
	"github.com/seriousm4x/upsnap/logger"
	"github.com/seriousm4x/upsnap/networking"
)

var PingRunning bool = false
var WakeShutdownRunning bool = false
var CronPing *cron.Cron
var CronWakeShutdown *cron.Cron

func RunPing(app *pocketbase.PocketBase) {
	PingRunning = true

	settingsPrivateRecords, err := app.Dao().FindRecordsByExpr("settings_private")
	if err != nil {
		logger.Error.Println(err)
	}

	// init cronjob
	CronPing = cron.New()
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
	CronPing.Run()
}

func RunWakeShutdown(app *pocketbase.PocketBase) {
	WakeShutdownRunning = true

	CronWakeShutdown = cron.New()
	devices, err := app.Dao().FindRecordsByExpr("devices")
	if err != nil {
		logger.Error.Println(err)
		return
	}
	for _, device := range devices {
		dev := device
		wake_cron := dev.GetString("wake_cron")
		wake_cron_enabled := dev.GetBool("wake_cron_enabled")
		shutdown_cron := dev.GetString("shutdown_cron")
		shutdown_cron_enabled := dev.GetBool("shutdown_cron_enabled")

		if wake_cron_enabled && wake_cron != "" {
			go func(d *models.Record) {
				_, err := CronWakeShutdown.AddFunc(wake_cron, func() {
					d.Set("status", "pending")
					if err := app.Dao().SaveRecord(d); err != nil {
						logger.Error.Println("Failed to save record:", err)
					}
					if err := networking.WakeDevice(d); err != nil {
						logger.Error.Println(err)
						d.Set("status", "offline")
						if err := app.Dao().SaveRecord(d); err != nil {
							logger.Error.Println("Failed to save record:", err)
						}
					} else {
						d.Set("status", "online")
						if err := app.Dao().SaveRecord(d); err != nil {
							logger.Error.Println("Failed to save record:", err)
						}
					}
				})
				if err != nil {
					logger.Error.Println(err)
				}
			}(dev)
		}

		if shutdown_cron_enabled && shutdown_cron != "" {
			go func(d *models.Record) {
				_, err := CronWakeShutdown.AddFunc(shutdown_cron, func() {
					d.Set("status", "pending")
					if err := app.Dao().SaveRecord(d); err != nil {
						logger.Error.Println("Failed to save record:", err)
					}
					if err := networking.ShutdownDevice(d); err != nil {
						logger.Error.Println(err)
						d.Set("status", "online")
						if err := app.Dao().SaveRecord(d); err != nil {
							logger.Error.Println("Failed to save record:", err)
						}
					} else {
						d.Set("status", "offline")
						if err := app.Dao().SaveRecord(d); err != nil {
							logger.Error.Println("Failed to save record:", err)
						}
					}
				})
				if err != nil {
					logger.Error.Println(err)
				}
			}(dev)
		}
	}
	CronWakeShutdown.Run()
}

func StopAll() {
	if PingRunning {
		logger.Info.Println("stopping ping cronjob")
		ctx := CronPing.Stop()
		<-ctx.Done()

	}
	if WakeShutdownRunning {
		logger.Info.Println("stopping wake/shutdown cronjob")
		ctx := CronWakeShutdown.Stop()
		<-ctx.Done()
	}
}
