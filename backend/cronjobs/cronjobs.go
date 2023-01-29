package cronjobs

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/robfig/cron/v3"
	"github.com/seriousm4x/upsnap/backend/logger"
	"github.com/seriousm4x/upsnap/backend/networking"
)

var Devices []*models.Record
var CronPing *cron.Cron
var CronWakeShutdown *cron.Cron

func RunPing(app *pocketbase.PocketBase) {
	settingsRecords, err := app.Dao().FindRecordsByExpr("settings")
	if err != nil {
		logger.Error.Println(err)
	}

	// init cronjob
	CronPing = cron.New()
	CronPing.AddFunc(settingsRecords[0].GetString("interval"), func() {
		// skip cron if no realtime clients connected
		realtimeClients := len(app.SubscriptionsBroker().Clients())
		if realtimeClients == 0 {
			return
		}
		// expand ports field
		expandFetchFunc := func(c *models.Collection, ids []string) ([]*models.Record, error) {
			return app.Dao().FindRecordsByIds(c.Id, ids, nil)
		}
		merr := app.Dao().ExpandRecords(Devices, []string{"ports"}, expandFetchFunc)
		if len(merr) > 0 {
			return
		}

		for _, device := range Devices {
			// ping
			go func(device *models.Record) {
				oldStatus := device.Get("status")
				newStatus := networking.PingDevice(device)
				if newStatus {
					if oldStatus == "offline" || oldStatus == "" {
						device.Set("status", "online")
						app.Dao().SaveRecord(device)
					}
				} else {
					if oldStatus == "online" || oldStatus == "" {
						device.Set("status", "offline")
						app.Dao().SaveRecord(device)
					}
				}
			}(device)

			// scan ports
			go func(device *models.Record) {
				ports, err := app.Dao().FindRecordsByIds("ports", device.GetStringSlice("ports"))
				if err != nil {
					logger.Error.Println(err)
				}
				for _, port := range ports {
					isUp := networking.CheckPort(device.GetString("ip"), port.GetString("number"))
					if isUp != port.GetBool("status") {
						port.Set("status", isUp)
						app.Dao().SaveRecord(port)
						device.RefreshUpdated()
						app.Dao().SaveRecord(device)
					}
				}
			}(device)
		}
	})
	logger.Debug.Println("CronPing entries:", CronPing.Entries())
	CronPing.Run()
}

func RunWakeShutdown() {
	CronWakeShutdown = cron.New()
	for _, device := range Devices {
		wake_cron := device.GetString("wake_cron")
		wake_cron_enabled := device.GetBool("wake_cron_enabled")
		shutdown_cron := device.GetString("shutdown_cron")
		shutdown_cron_enabled := device.GetBool("shutdown_cron_enabled")

		if wake_cron_enabled && wake_cron != "" {
			CronWakeShutdown.AddFunc(wake_cron, func() {
				if err := networking.WakeDevice(device); err != nil {
					logger.Error.Println(err)
				}
			})
		}

		if shutdown_cron_enabled && shutdown_cron != "" {
			CronWakeShutdown.AddFunc(shutdown_cron, func() {
				if err := networking.ShutdownDevice(device); err != nil {
					logger.Error.Println(err)
				}
			})
		}
	}
	logger.Debug.Println("CronWakeShutdown entries:", CronWakeShutdown.Entries())
	CronWakeShutdown.Run()
}
