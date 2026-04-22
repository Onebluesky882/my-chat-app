package websocket

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func WebSocketHandler(wsHandler func(*websocket.Conn)) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		wsHandler(c)
	})
}
