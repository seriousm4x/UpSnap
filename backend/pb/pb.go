package pb

import (
	"io/fs"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/seriousm4x/upsnap/cronjobs"
	"github.com/seriousm4x/upsnap/logger"
	_ "github.com/seriousm4x/upsnap/migrations"
)

var App *pocketbase.PocketBase
var Version = "(untracked)"

func StartPocketBase(distDirFS fs.FS) {
	App = pocketbase.New()
	App.RootCmd.Short = "UpSnap CLI"
	App.RootCmd.Version = Version

	// auto migrate db
	migratecmd.MustRegister(App, App.RootCmd, &migratecmd.Options{
		Automigrate: true,
	})

	// event hooks
	App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// set static website path
		e.Router.GET("/*", apis.StaticDirectoryHandler(distDirFS, true))

		// add wake route to api
		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/upsnap/wake/:id",
			Handler: HandlerWake,
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(App),
				apis.RequireAdminOrRecordAuth("users"),
			},
		})

		// add shutdown route to api
		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/upsnap/shutdown/:id",
			Handler: HandlerShutdown,
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(App),
				apis.RequireAdminOrRecordAuth("users"),
			},
		})

		// add network scan route to api
		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/upsnap/scan",
			Handler: HandlerScan,
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(App),
				apis.RequireAdminAuth(),
			},
		})

		// import environment and set settings
		if err := importSettings(); err != nil {
			return err
		}

		// reset device states and run ping cronjob
		if err := resetDeviceStates(); err != nil {
			return err
		}

		// run cronjobs
		go cronjobs.RunPing(App)
		go cronjobs.RunWakeShutdown(App)

		// restart ping cronjobs or wake/shutdown cronjobs on model update
		// add event hook before starting server.
		// using this outside App.OnBeforeServe() would not work
		App.OnModelAfterUpdate().Add(func(e *core.ModelEvent) error {
			if e.Model.TableName() == "settings_private" {
				for _, job := range cronjobs.CronPing.Entries() {
					cronjobs.CronPing.Remove(job.ID)
				}
				go cronjobs.RunPing(App)
			} else if e.Model.TableName() == "devices" {
				if err := refreshDeviceList(); err != nil {
					logger.Error.Println(err)
					return err
				}
				for _, job := range cronjobs.CronWakeShutdown.Entries() {
					cronjobs.CronWakeShutdown.Remove(job.ID)
				}
				go cronjobs.RunWakeShutdown(App)
			}
			return nil
		})
		return nil
	})

	// refresh the device list on database events
	App.OnModelAfterCreate().Add(func(e *core.ModelEvent) error {
		if e.Model.TableName() == "_admins" {
			if err := setSetupCompleted(); err != nil {
				logger.Error.Println(err)
				return err
			}
		} else {
			if err := refreshDeviceList(); err != nil {
				logger.Error.Println(err)
				return err
			}
		}
		return nil
	})
	App.OnModelAfterDelete().Add(func(e *core.ModelEvent) error {
		if e.Model.TableName() == "_admins" {
			if err := setSetupCompleted(); err != nil {
				logger.Error.Println(err)
				return err
			}
		} else {
			if err := refreshDeviceList(); err != nil {
				logger.Error.Println(err)
				return err
			}
		}
		return nil
	})

	// start pocketbase
	if err := App.Start(); err != nil {
		logger.Error.Fatalln(err)
	}
}

func importSettings() error {
	// get first settingsPrivate record
	settingsPrivateRecords, err := App.Dao().FindRecordsByExpr("settings_private")
	if err != nil {
		return err
	}
	settingsPrivateCollection, err := App.Dao().FindCollectionByNameOrId("settings_private")
	if err != nil {
		return err
	}
	settingsPrivate := models.NewRecord(settingsPrivateCollection)
	if len(settingsPrivateRecords) > 0 {
		settingsPrivate = settingsPrivateRecords[0]
	}

	// get first settingsPublic record
	settingsPublicRecords, err := App.Dao().FindRecordsByExpr("settings_public")
	if err != nil {
		return err
	}
	settingsPublicCollection, err := App.Dao().FindCollectionByNameOrId("settings_public")
	if err != nil {
		return err
	}
	settingsPublic := models.NewRecord(settingsPublicCollection)
	if len(settingsPublicRecords) > 0 {
		settingsPublic = settingsPublicRecords[0]
	}

	// set ping interval settings. priority:
	// 1st: env var
	// 2nd: database entry
	// 3rd: default values
	interval := "@every 3s"
	if settingsPrivate.GetString("interval") != "" {
		interval = settingsPrivate.GetString("interval")
	}
	if os.Getenv("UPSNAP_INTERVAL") != "" {
		interval = os.Getenv("UPSNAP_INTERVAL")
	}

	// set private settings
	settingsPrivate.Set("interval", interval)
	if scanRange := os.Getenv("UPSNAP_SCAN_RANGE"); scanRange != "" {
		settingsPrivate.Set("scan_range", scanRange)
	}

	// set public settings
	if websiteTitle := os.Getenv("UPSNAP_WEBSITE_TITLE"); websiteTitle != "" {
		settingsPublic.Set("website_title", websiteTitle)
	}

	// save records
	if err := App.Dao().SaveRecord(settingsPrivate); err != nil {
		return err
	}
	if err := App.Dao().SaveRecord(settingsPublic); err != nil {
		return err
	}
	if err := setSetupCompleted(); err != nil {
		logger.Error.Println(err)
		return err
	}

	logger.Info.Println("Ping interval set to", interval)
	return nil
}

func resetDeviceStates() error {
	devices, err := App.Dao().FindRecordsByExpr("devices")
	if err != nil {
		return err
	}
	for _, device := range devices {
		device.Set("status", "offline")
		if err := App.Dao().SaveRecord(device); err != nil {
			return err
		}
	}
	cronjobs.Devices = devices
	return nil
}

func refreshDeviceList() error {
	var err error
	if cronjobs.Devices, err = App.Dao().FindRecordsByExpr("devices"); err != nil {
		return err
	}
	return nil
}

func setSetupCompleted() error {
	totalAdmins, err := App.Dao().TotalAdmins()
	if err != nil {
		return err
	}
	settingsPublicRecords, err := App.Dao().FindRecordsByExpr("settings_public")
	if err != nil {
		return err
	}
	if totalAdmins > 0 {
		settingsPublicRecords[0].Set("setup_completed", true)
	} else {
		settingsPublicRecords[0].Set("setup_completed", false)
	}
	if err := App.Dao().SaveRecord(settingsPublicRecords[0]); err != nil {
		return err
	}
	return nil
}
