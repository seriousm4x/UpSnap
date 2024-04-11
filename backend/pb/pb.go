package pb

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"

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
	// set data dir
	// use "./pb_data" if it's in the same dir as upsnap binary
	// else use os.UserConfigDir() / upsnap
	var dataDir string
	baseDir, err := os.Getwd()
	if err != nil {
		logger.Error.Fatalln(err)
	}
	pb_data := path.Join(baseDir, "pb_data")
	if _, err = os.Stat(pb_data); err == nil {
		dataDir = pb_data
	} else if os.IsNotExist(err) {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			logger.Error.Fatalln(err)
		}
		upsnap_data := path.Join(userConfigDir, "upsnap")
		if _, err = os.Stat(upsnap_data); err == nil {
			dataDir = upsnap_data
		} else if os.IsNotExist(err) {
			if err := os.MkdirAll(upsnap_data, 0700); err != nil {
				logger.Error.Fatalln(err)
			}
		}
	}

	// create app
	App = pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir: dataDir,
	})
	App.RootCmd.Short = "UpSnap CLI"
	App.RootCmd.Version = Version

	// auto migrate db
	migratecmd.MustRegister(App, App.RootCmd, migratecmd.Config{
		Automigrate: true,
	})

	// event hooks
	App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(distDirFS, true))

		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/upsnap/wake/:id",
			Handler: HandlerWake,
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(App),
				RequireUpSnapPermission(),
			},
		})

		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/upsnap/sleep/:id",
			Handler: HandlerSleep,
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(App),
				RequireUpSnapPermission(),
			},
		})

		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/upsnap/reboot/:id",
			Handler: HandlerReboot,
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(App),
				RequireUpSnapPermission(),
			},
		})

		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/upsnap/shutdown/:id",
			Handler: HandlerShutdown,
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(App),
				RequireUpSnapPermission(),
			},
		})

		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/upsnap/scan",
			Handler: HandlerScan,
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(App),
				apis.RequireAdminAuth(),
			},
		})

		if err := importSettings(); err != nil {
			return err
		}

		if err := resetDeviceStates(); err != nil {
			return err
		}

		cronjobs.SetPingJobs(App)
		cronjobs.StartPing()
		cronjobs.SetWakeShutdownJobs(App)
		cronjobs.StartWakeShutdown()

		// restart ping cronjobs or wake/shutdown cronjobs on model update
		// add event hook before starting server.
		// using this outside App.OnBeforeServe() would not work
		App.OnModelAfterUpdate("settings_private", "devices").Add(func(e *core.ModelEvent) error {
			if e.Model.TableName() == "settings_private" {
				cronjobs.SetPingJobs(App)
			} else if e.Model.TableName() == "devices" {
				// only restart wake/shutdown cronjobs if new model's cron changed
				record := e.Model.(*models.Record)
				newRecord := record.CleanCopy()
				oldRecord := record.OriginalCopy()

				newWakeCron := newRecord.GetString("wake_cron")
				newWakeCmd := newRecord.GetString("wake_cmd")
				newWakeCronEnabled := newRecord.GetBool("wake_cron_enabled")
				newShutdownCron := newRecord.GetString("shutdown_cron")
				newShutdownCmd := newRecord.GetString("shutdown_cmd")
				newShutdownCronEnabled := newRecord.GetBool("shutdown_cron_enabled")

				oldWakeCron := oldRecord.GetString("wake_cron")
				oldWakeCmd := oldRecord.GetString("wake_cmd")
				oldWakeCronEnabled := oldRecord.GetBool("wake_cron_enabled")
				oldShutdownCron := oldRecord.GetString("shutdown_cron")
				oldShutdownCmd := oldRecord.GetString("shutdown_cmd")
				oldShutdownCronEnabled := oldRecord.GetBool("shutdown_cron_enabled")

				if newWakeCron != oldWakeCron ||
					newWakeCmd != oldWakeCmd ||
					newWakeCronEnabled != oldWakeCronEnabled ||
					newShutdownCron != oldShutdownCron ||
					newShutdownCronEnabled != oldShutdownCronEnabled ||
					newShutdownCmd != oldShutdownCmd {
					cronjobs.SetWakeShutdownJobs(App)
				}
			}
			return nil
		})
		return nil
	})

	App.OnModelAfterCreate().Add(func(e *core.ModelEvent) error {
		if e.Model.TableName() == "_admins" {
			if err := setSetupCompleted(); err != nil {
				logger.Error.Println(err)
				return err
			}
			return nil
		} else if e.Model.TableName() == "devices" {
			// when a device is created, give the user all rights to the device he just created
			deviceRec := e.Model.(*models.Record)
			userId := deviceRec.GetString("created_by")

			var permissionRec *models.Record
			permissionRec, err := App.Dao().FindFirstRecordByFilter("permissions",
				fmt.Sprintf("user.id = '%s'", userId))
			if err != nil && err.Error() != "sql: no rows in result set" {
				logger.Error.Println(err)
				return err
			} else if permissionRec != nil {
				permissionRec.Set("read", append(permissionRec.GetStringSlice("read"), deviceRec.Id))
				permissionRec.Set("update", append(permissionRec.GetStringSlice("update"), deviceRec.Id))
				permissionRec.Set("delete", append(permissionRec.GetStringSlice("delete"), deviceRec.Id))
				permissionRec.Set("power", append(permissionRec.GetStringSlice("power"), deviceRec.Id))
				if err := App.Dao().SaveRecord(permissionRec); err != nil {
					logger.Error.Println(err)
					return err
				}
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
		}
		return nil
	})

	App.OnTerminate().Add(func(e *core.TerminateEvent) error {
		cronjobs.StopAll()
		return nil
	})

	if err := App.Start(); err != nil {
		logger.Error.Fatalln(err)
	}
}

func importSettings() error {
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

	settingsPrivate.Set("interval", interval)
	if scanRange := os.Getenv("UPSNAP_SCAN_RANGE"); scanRange != "" {
		settingsPrivate.Set("scan_range", scanRange)
	}

	if websiteTitle := os.Getenv("UPSNAP_WEBSITE_TITLE"); websiteTitle != "" {
		settingsPublic.Set("website_title", websiteTitle)
	}

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
