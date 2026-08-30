package pb

import (
	"fmt"
	"io/fs"
	"net"
	"os"
	"path"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/ghupdate"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/robfig/cron/v3"
	"github.com/seriousm4x/upsnap/cronjobs"
	"github.com/seriousm4x/upsnap/iptracking"
	"github.com/seriousm4x/upsnap/logger"
	_ "github.com/seriousm4x/upsnap/migrations"
)

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
	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir: dataDir,
	})
	app.RootCmd.Short = "UpSnap CLI"
	app.RootCmd.Version = Version

	// auto migrate db
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: true,
	})

	// GitHub selfupdate
	ghupdate.MustRegister(app, app.RootCmd, ghupdate.Config{
		Owner:             "seriousm4x",
		Repo:              "UpSnap",
		ArchiveExecutable: "upsnap",
	})

	// event hooks
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		logger.InitLogger(app)

		se.Router.GET("/{path...}", apis.Static(distDirFS, true))
		se.Router.GET("/api/upsnap/wake/{id}", HandlerWake).Bind(RequireUpSnapPermission())
		se.Router.GET("/api/upsnap/wakegroup/{id}", HandlerWakeGroup).Bind(RequireUpSnapPermission())
		se.Router.GET("/api/upsnap/sleep/{id}", HandlerSleep).Bind(RequireUpSnapPermission())
		se.Router.GET("/api/upsnap/reboot/{id}", HandlerReboot).Bind(RequireUpSnapPermission())
		se.Router.GET("/api/upsnap/shutdown/{id}", HandlerShutdown).Bind(RequireUpSnapPermission())
		se.Router.GET("/api/upsnap/shutdowngroup/{id}", HandlerShutdownGroup).Bind(RequireUpSnapPermission())
		se.Router.GET("/api/upsnap/scan", HandlerScan).Bind(RequireScanDevicesPermission())
		se.Router.POST("/api/upsnap/validate-cron", HandlerValidateCron)
		se.Router.GET("/api/upsnap/manifest.webmanifest", HandlerWebsiteManifest)

		if err := importSettings(app); err != nil {
			return err
		}

		if err := resetDeviceStates(app); err != nil {
			return err
		}

		cronjobs.SetPingJobs(app)
		cronjobs.StartPing()
		cronjobs.SetWakeShutdownJobs(app)
		cronjobs.StartWakeShutdown()

		// restart ping cronjobs or wake/shutdown cronjobs on model update
		// add event hook before starting server.
		// using this outside App.OnBeforeServe() would not work
		app.OnModelAfterUpdateSuccess("settings_private", "devices").BindFunc(func(e *core.ModelEvent) error {
			if e.Model.TableName() == "settings_private" {
				cronjobs.SetPingJobs(app)
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

				newIp := newRecord.GetString("ip")
				oldIp := oldRecord.GetString("ip")
				newPingCmd := newRecord.GetString("ping_cmd")
				oldPingCmd := oldRecord.GetString("ping_cmd")

				if newWakeCron != oldWakeCron ||
					newWakeCmd != oldWakeCmd ||
					newWakeCronEnabled != oldWakeCronEnabled ||
					newShutdownCron != oldShutdownCron ||
					newShutdownCronEnabled != oldShutdownCronEnabled ||
					newShutdownCmd != oldShutdownCmd ||
					newIp != oldIp ||
					newPingCmd != oldPingCmd {
					cronjobs.SetWakeShutdownJobs(app)
				}
			}
			return e.Next()
		})
		return se.Next()
	})

	app.OnModelAfterCreateSuccess("_superusers", "devices").BindFunc(func(e *core.ModelEvent) error {
		if e.Model.TableName() == "_superusers" {
			// when pocketbase creates it's default superuser, do not trigger setSetupCompleted()
			if e.Model.(*core.Record).Email() == core.DefaultInstallerEmail {
				return e.Next()
			}
			if err := setSetupCompleted(e.App); err != nil {
				logger.Error.Println(err)
				return err
			}
			return e.Next()
		} else if e.Model.TableName() == "devices" {
			// when a device is created, give the user all rights to the device he just created
			deviceRec := e.Model.(*core.Record)
			userId := deviceRec.GetString("created_by")

			var permissionRec *core.Record
			permissionRec, err := app.FindFirstRecordByFilter("permissions",
				fmt.Sprintf("user.id = '%s'", userId))
			if err != nil && err.Error() != "sql: no rows in result set" {
				logger.Error.Println(err)
				return err
			} else if permissionRec != nil {
				permissionRec.Set("read", append(permissionRec.GetStringSlice("read"), deviceRec.Id))
				permissionRec.Set("update", append(permissionRec.GetStringSlice("update"), deviceRec.Id))
				permissionRec.Set("delete", append(permissionRec.GetStringSlice("delete"), deviceRec.Id))
				permissionRec.Set("power", append(permissionRec.GetStringSlice("power"), deviceRec.Id))
				if err := app.Save(permissionRec); err != nil {
					logger.Error.Println(err)
					return err
				}
			}
		}
		return e.Next()
	})

	app.OnModelAfterDeleteSuccess("_superusers").BindFunc(func(e *core.ModelEvent) error {
		if err := setSetupCompleted(e.App); err != nil {
			logger.Error.Println(err)
			return err
		}
		return e.Next()
	})

	app.OnRecordValidate("devices").BindFunc(func(e *core.RecordEvent) error {
		ip := net.ParseIP(e.Record.GetString("netmask"))
		if ip == nil {
			logger.Warning.Println("device", e.Record.GetString("name"), "does not contain a valid netmask! Will be set to 255.255.255.255")
			e.Record.Set("netmask", "255.255.255.255")
			saveErr := e.App.Save(e.Record)
			if saveErr != nil {
				return saveErr
			}
		}
		return e.Next()
	})

	app.OnRealtimeConnectRequest().BindFunc(func(e *core.RealtimeConnectRequestEvent) error {
		// a client just became active: catch up on any tracking sweep that
		// was skipped while nobody was connected
		settings, err := e.App.FindFirstRecordByFilter("settings_private", "")
		if err != nil {
			logger.Error.Println(err)
		} else if settings.GetBool("lazy_ping") && settings.GetString("track_ip_interval") != "" {
			go iptracking.CatchUpSweep(e.App)
		}
		return e.Next()
	})

	app.OnRealtimeConnectRequest().BindFunc(func(e *core.RealtimeConnectRequestEvent) error {
		// PocketBase has no realtime disconnect hook. When a browser refreshes,
		// the old client may disappear before the new connection is registered.
		// Delay the offline transition and use the connecting client's id to
		// distinguish a refresh from a genuine disconnect.
		clientID := e.Client.Id()
		defer func() {
			if len(e.App.SubscriptionsBroker().Clients()) != 0 {
				return
			}
			time.AfterFunc(2*time.Second, func() {
				if _, err := e.App.SubscriptionsBroker().ClientById(clientID); err == nil {
					return
				}
				if len(e.App.SubscriptionsBroker().Clients()) != 0 {
					return
				}
				allDevices, err := app.FindAllRecords("devices", dbx.NewExp("shutdown_cmd != ''"))
				if err != nil {
					logger.Error.Println(err)
					return
				}
				for _, device := range allDevices {
					device.Set("status", "offline")
					if err := app.Save(device); err != nil {
						logger.Error.Println(err)
						return
					}
				}
			})
		}()
		return e.Next()
	})

	app.OnTerminate().BindFunc(func(e *core.TerminateEvent) error {
		cronjobs.StopAll()
		return e.Next()
	})

	// check for custom http listen address in env var, else use default
	httpListen := os.Getenv("UPSNAP_HTTP_LISTEN")
	if httpListen != "" {
		if err := app.Bootstrap(); err != nil {
			logger.Error.Fatalln(err)
		}
		if err := apis.Serve(app, apis.ServeConfig{
			HttpAddr:        httpListen,
			ShowStartBanner: true,
		}); err != nil {
			logger.Error.Fatalln(err)
		}
	} else {
		if err := app.Start(); err != nil {
			logger.Error.Fatalln(err)
		}
	}

}

func importSettings(app core.App) error {
	settingsPrivateRecords, err := app.FindAllRecords("settings_private")
	if err != nil {
		return err
	}
	settingsPrivateCollection, err := app.FindCollectionByNameOrId("settings_private")
	if err != nil {
		return err
	}
	settingsPrivate := core.NewRecord(settingsPrivateCollection)
	if len(settingsPrivateRecords) > 0 {
		settingsPrivate = settingsPrivateRecords[0]
	}

	settingsPublicRecords, err := app.FindAllRecords("settings_public")
	if err != nil {
		return err
	}
	settingsPublicCollection, err := app.FindCollectionByNameOrId("settings_public")
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
	defaultInterval := "*/3 * * * * *"
	interval := defaultInterval
	if settingsPrivate.GetString("interval") != "" {
		interval = settingsPrivate.GetString("interval")
	}
	if os.Getenv("UPSNAP_INTERVAL") != "" {
		interval = os.Getenv("UPSNAP_INTERVAL")
	}

	// validate interval before saving
	p := cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	if _, err := p.Parse(interval); err != nil {
		if interval == "*/3 * * * *" || interval == "@every 3s" {
			settingsPrivate.Set("interval", defaultInterval)
			if e := app.Save(settingsPrivate); e != nil {
				logger.Error.Println(e)
			}
		} else {
			logger.Error.Printf("'%s' ping interval is not valid.\n", interval)
			logger.Error.Println("Please go to '/settings/' and change your ping interval.")
			logger.Error.Println("Falling back to default interval: " + defaultInterval)
		}
		interval = defaultInterval
	}

	settingsPrivate.Set("interval", interval)
	if scanRange := os.Getenv("UPSNAP_SCAN_RANGE"); scanRange != "" {
		settingsPrivate.Set("scan_range", scanRange)
	}

	if websiteTitle := os.Getenv("UPSNAP_WEBSITE_TITLE"); websiteTitle != "" {
		settingsPublic.Set("website_title", websiteTitle)
	}

	if err := app.Save(settingsPrivate); err != nil {
		return err
	}
	if err := app.Save(settingsPublic); err != nil {
		return err
	}
	if err := setSetupCompleted(app); err != nil {
		logger.Error.Println(err)
		return err
	}

	return nil
}

func resetDeviceStates(app core.App) error {
	devices, err := app.FindAllRecords("devices")
	if err != nil {
		return err
	}
	for _, device := range devices {
		device.Set("status", "offline")
		if err := app.Save(device); err != nil {
			return err
		}
	}
	return nil
}

func setSetupCompleted(app core.App) error {
	allSuperusers, err := app.FindAllRecords(core.CollectionNameSuperusers)
	if err != nil {
		return err
	}
	realSuperusers := 0
	for _, r := range allSuperusers {
		// exclude the temporary PocketBase installer account from the count
		if r.Email() != core.DefaultInstallerEmail {
			realSuperusers++
		}
	}
	settingsPublicRecords, err := app.FindAllRecords("settings_public")
	if err != nil {
		return err
	}
	if realSuperusers > 0 {
		settingsPublicRecords[0].Set("setup_completed", true)
	} else {
		settingsPublicRecords[0].Set("setup_completed", false)
	}
	if err := app.Save(settingsPublicRecords[0]); err != nil {
		return err
	}
	return nil
}
