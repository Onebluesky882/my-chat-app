package websocket

import "github.com/gofiber/fiber/v3"

func WSRouter(app *fiber.App, wsSvc *Service) {
	app.All("/ws", wsSvc.HandleWS)
	app.Get("/test-broadcast", func(c fiber.Ctx) error {

		wsSvc.Broadcast(
			"room-1",
			[]byte("hello realtime"),
		)
		return c.SendString("ok")

	})
}
