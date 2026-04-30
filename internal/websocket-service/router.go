package websocket

import (
	"github.com/Onebluesky882/my-chat-app/internal/chat-service"
	fastwsocket "github.com/fasthttp/websocket"
	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fasthttp"
)

var upgrader = fastwsocket.FastHTTPUpgrader{
	CheckOrigin: func(r *fasthttp.RequestCtx) bool {
		return true
	},
}

func WSRouter(app *fiber.App, wsSvc *Service, chatSvc *chat.Service) {
	group := app.Group("/ws")

	group.Get("/", func(c fiber.Ctx) error {
		roomID := c.Query("room_id")
		userID := c.Query("user_id")
		rID, err := gocql.ParseUUID(roomID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid room_id"})
		}
		uID, err := gocql.ParseUUID(userID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid user_id"})
		}
		return upgrader.Upgrade(c.RequestCtx(), HandleWS(wsSvc, rID, uID))
	})

	group.Post("/send", HandleSendToUser(wsSvc, chatSvc))
}
