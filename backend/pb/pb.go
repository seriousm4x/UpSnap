package pb

import (
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/seriousm4x/upsnap/backend/cronjobs"
	"github.com/seriousm4x/upsnap/backend/logger"
	_ "github.com/seriousm4x/upsnap/backend/migrations"
	"github.com/seriousm4x/upsnap/backend/networking"
)

var app *pocketbase.PocketBase

func StartPocketBase() {
	app = pocketbase.New()

	// auto migrate db
	migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
		Automigrate: true,
	})

	// event hooks
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// set static website path
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), true))

		// add wake route to api
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/upsnap/wake/:id",
			Handler: func(c echo.Context) error {
				record, err := app.Dao().FindFirstRecordByData("devices", "id", c.PathParam("id"))
				if err != nil {
					return apis.NewNotFoundError("The device does not exist.", err)
				}
				go func(*models.Record) {
					record.Set("status", "pending")
					app.Dao().SaveRecord(record)
					isOnline := networking.WakeDevice(record)
					if isOnline {
						record.Set("status", "online")
					} else {
						record.Set("status", "offline")
					}
					app.Dao().SaveRecord(record)
				}(record)
				return c.JSON(http.StatusOK, record)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})

		// add shutdown route to api
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/upsnap/shutdown/:id",
			Handler: func(c echo.Context) error {
				record, err := app.Dao().FindFirstRecordByData("devices", "id", c.PathParam("id"))
				if err != nil {
					return apis.NewNotFoundError("The device does not exist.", err)
				}
				shutdown_cmd := record.GetString("shutdown_cmd")
				if shutdown_cmd != "" {
					cmd := exec.Command(shutdown_cmd)
					if err := cmd.Run(); err != nil {
						logger.Error.Println(err)
						return apis.NewBadRequestError(err.Error(), record)
					}
				}
				return nil
			},
		})

		// add shutdown route to api
		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/upsnap/scan",
			Handler: networking.ScanNetwork,
		})

		// reset device states and run ping cronjob
		devices, err := app.Dao().FindRecordsByExpr("devices")
		if err != nil {
			return err
		}
		for _, device := range devices {
			device.Set("status", "offline")
			if err := app.Dao().SaveRecord(device); err != nil {
				return err
			}
		}
		cronjobs.Devices = devices
		go cronjobs.RunCron(app)
		return nil
	})

	// refresh the device list on database events
	app.OnModelAfterCreate().Add(func(e *core.ModelEvent) error {
		refreshDeviceList()
		return nil
	})
	app.OnModelAfterUpdate().Add(func(e *core.ModelEvent) error {
		refreshDeviceList()
		return nil
	})
	app.OnModelAfterDelete().Add(func(e *core.ModelEvent) error {
		refreshDeviceList()
		return nil
	})

	// start pocketbase
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

}

func refreshDeviceList() {
	var err error
	cronjobs.Devices, err = app.Dao().FindRecordsByExpr("devices")
	if err != nil {
		logger.Error.Println(err)
	}
}
