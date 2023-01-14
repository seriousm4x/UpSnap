package cronjobs

import (
	"fmt"
	"os"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/robfig/cron"
	"github.com/seriousm4x/upsnap/backend/logger"
	"github.com/seriousm4x/upsnap/backend/networking"
)

var Devices []*models.Record

func RunCron(app *pocketbase.PocketBase) {
	var interval time.Duration
	var err error

	// get ping interval from env var
	// no env: fallback to 5 seconds
	// below 1 sec: fallback to 1 second
	interval_default := 5 * time.Second
	interval_minimum := 1 * time.Second
	interval_env := os.Getenv("UPSNAP_INTERVAL")
	if interval_env == "" {
		interval = interval_default
	} else {
		interval, err = time.ParseDuration(interval_env)
		if err != nil {
			logger.Error.Println(err)
			interval = interval_default
		}
		if interval < interval_minimum {
			logger.Warning.Printf("Ping interval below %s is not recommended", interval_minimum)
			interval = interval_minimum
		}
	}
	logger.Debug.Println("Ping interval set to", interval)

	// init cronjob
	c := cron.New()
	c.AddFunc(fmt.Sprintf("@every %s", interval), func() {
		for _, device := range Devices {
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
		}
	})
	c.Run()
}
