package pb

import (
	"fmt"
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

		go cronjobs.RunPing(App)
		go cronjobs.RunWakeShutdown(App)

		// restart ping cronjobs or wake/shutdown cronjobs on model update
		// add event hook before starting server.
		// using this outside App.OnBeforeServe() would not work
		App.OnModelAfterUpdate("settings_private", "devices").Add(func(e *core.ModelEvent) error {
			if e.Model.TableName() == "settings_private" {
				for _, job := range cronjobs.CronPing.Entries() {
					cronjobs.CronPing.Remove(job.ID)
				}
				go cronjobs.RunPing(App)
			} else if e.Model.TableName() == "devices" {
				for _, job := range cronjobs.CronWakeShutdown.Entries() {
					cronjobs.CronWakeShutdown.Remove(job.ID)
				}
				go cronjobs.RunWakeShutdown(App)
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

	App.OnTerminate().PreAdd(func(e *core.TerminateEvent) error {
		logger.Info.Println("Stopping cronjobs")
		ctx := cronjobs.CronPing.Stop()
		<-ctx.Done()
		ctx = cronjobs.CronWakeShutdown.Stop()
		<-ctx.Done()
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
		d := device
		d.Set("status", "offline")
		if err := App.Dao().SaveRecord(d); err != nil {
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
