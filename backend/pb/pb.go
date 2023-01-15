package pb

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/seriousm4x/upsnap/backend/cronjobs"
	"github.com/seriousm4x/upsnap/backend/logger"
	_ "github.com/seriousm4x/upsnap/backend/migrations"
)

var App *pocketbase.PocketBase

func StartPocketBase() {
	App = pocketbase.New()

	// auto migrate db
	migratecmd.MustRegister(App, App.RootCmd, &migratecmd.Options{
		Automigrate: true,
	})

	// event hooks
	App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// set static website path
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), true))

		// add wake route to api
		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/upsnap/wake/:id",
			Handler: HandlerWake,
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(App),
			},
		})

		// add shutdown route to api
		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/upsnap/shutdown/:id",
			Handler: HandlerShutdown,
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(App),
			},
		})

		// add network scan route to api
		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/upsnap/scan",
			Handler: HandlerScan,
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(App),
			},
		})

		// reset device states and run ping cronjob
		if err := resetDeviceStates(); err != nil {
			return err
		}

		// run ping cronjob
		go cronjobs.RunCron(App)

		return nil
	})

	// refresh the device list on database events
	App.OnModelAfterCreate().Add(func(e *core.ModelEvent) error {
		refreshDeviceList()
		return nil
	})
	App.OnModelAfterUpdate().Add(func(e *core.ModelEvent) error {
		refreshDeviceList()
		return nil
	})
	App.OnModelAfterDelete().Add(func(e *core.ModelEvent) error {
		refreshDeviceList()
		return nil
	})

	// start pocketbase
	if err := App.Start(); err != nil {
		log.Fatal(err)
	}

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

func refreshDeviceList() {
	var err error
	cronjobs.Devices, err = App.Dao().FindRecordsByExpr("devices")
	if err != nil {
		logger.Error.Println(err)
	}
}
