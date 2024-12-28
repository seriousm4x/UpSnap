package pb

import (
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
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
	App.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/{path...}", apis.Static(distDirFS, true))

		se.Router.GET("/api/upsnap/wake/{id}", HandlerWake).Bind(RequireUpSnapPermission())
		se.Router.GET("/api/upsnap/sleep/{id}", HandlerSleep).Bind(RequireUpSnapPermission())
		se.Router.GET("/api/upsnap/reboot/{id}", HandlerReboot).Bind(RequireUpSnapPermission())
		se.Router.GET("/api/upsnap/shutdown/{id}", HandlerShutdown).Bind(RequireUpSnapPermission())
		se.Router.GET("/api/upsnap/scan", HandlerScan).Bind(apis.RequireSuperuserAuth())

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
		App.OnModelAfterUpdateSuccess("settings_private", "devices").BindFunc(func(e *core.ModelEvent) error {
			if e.Model.TableName() == "settings_private" {
				cronjobs.SetPingJobs(App)
			} else if e.Model.TableName() == "devices" {
				// only restart wake/shutdown cronjobs if new model's cron changed
				record := e.Model.(*core.Record)
				newRecord := record.Fresh()
				oldRecord := record.Original()

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
		return se.Next()
	})

	App.OnModelAfterCreateSuccess().BindFunc(func(e *core.ModelEvent) error {
		if e.Model.TableName() == "_admins" {
			if err := setSetupCompleted(); err != nil {
				logger.Error.Println(err)
				return err
			}
			return nil
		} else if e.Model.TableName() == "devices" {
			// when a device is created, give the user all rights to the device he just created
			deviceRec := e.Model.(*core.Record)
			userId := deviceRec.GetString("created_by")

			var permissionRec *core.Record
			permissionRec, err := App.FindFirstRecordByFilter("permissions",
				fmt.Sprintf("user.id = '%s'", userId))
			if err != nil && err.Error() != "sql: no rows in result set" {
				logger.Error.Println(err)
				return err
			} else if permissionRec != nil {
				permissionRec.Set("read", append(permissionRec.GetStringSlice("read"), deviceRec.Id))
				permissionRec.Set("update", append(permissionRec.GetStringSlice("update"), deviceRec.Id))
				permissionRec.Set("delete", append(permissionRec.GetStringSlice("delete"), deviceRec.Id))
				permissionRec.Set("power", append(permissionRec.GetStringSlice("power"), deviceRec.Id))
				if err := App.Save(permissionRec); err != nil {
					logger.Error.Println(err)
					return err
				}
			}
		}
		return nil
	})
	App.OnModelAfterDeleteSuccess().BindFunc(func(e *core.ModelEvent) error {
		if e.Model.TableName() == "_admins" {
			if err := setSetupCompleted(); err != nil {
				logger.Error.Println(err)
				return err
			}
		}
		return nil
	})

	App.OnTerminate().BindFunc(func(e *core.TerminateEvent) error {
		cronjobs.StopAll()
		return nil
	})

	if err := App.Start(); err != nil {
		logger.Error.Fatalln(err)
	}
}

func importSettings() error {
	settingsPrivateRecords, err := App.FindAllRecords("settings_private")
	if err != nil {
		return err
	}
	settingsPrivateCollection, err := App.FindCollectionByNameOrId("settings_private")
	if err != nil {
		return err
	}
	settingsPrivate := core.NewRecord(settingsPrivateCollection)
	if len(settingsPrivateRecords) > 0 {
		settingsPrivate = settingsPrivateRecords[0]
	}

	settingsPublicRecords, err := App.FindAllRecords("settings_public")
	if err != nil {
		return err
	}
	settingsPublicCollection, err := App.FindCollectionByNameOrId("settings_public")
	if err != nil {
		return err
	}
	settingsPublic := core.NewRecord(settingsPublicCollection)
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

	if err := App.Save(settingsPrivate); err != nil {
		return err
	}
	if err := App.Save(settingsPublic); err != nil {
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
	devices, err := App.FindAllRecords("devices")
	if err != nil {
		return err
	}
	for _, device := range devices {
		device.Set("status", "offline")
		if err := App.Save(device); err != nil {
			return err
		}
	}
	return nil
}

func setSetupCompleted() error {
	totalAdmins, err := App.CountRecords(core.CollectionNameSuperusers)
	if err != nil {
		return err
	}
	settingsPublicRecords, err := App.FindAllRecords("settings_public")
	if err != nil {
		return err
	}
	if totalAdmins > 0 {
		settingsPublicRecords[0].Set("setup_completed", true)
	} else {
		settingsPublicRecords[0].Set("setup_completed", false)
	}
	if err := App.Save(settingsPublicRecords[0]); err != nil {
		return err
	}
	return nil
}
