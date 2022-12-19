package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"github.com/seriousm4x/UpSnap/controllers"
)

func main() {
	app := fiber.New()
	app.Use(recover.New())
	app.Get("/metrics", monitor.New(monitor.Config{Title: "UpSnap Metrics Page"}))

	devicesGroup := app.Group("/devices")
	devicesGroup.Get("/", controllers.GetDevices)
	devicesGroup.Post("/", controllers.CreateDevice)
	devicesGroup.Patch("/:id", controllers.PatchDevice)
	devicesGroup.Delete("/:id", controllers.DeleteDevice)

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)

			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}

	}))
	log.Fatal(app.Listen(":3000"))
}
