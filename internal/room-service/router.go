package room

import "github.com/gofiber/fiber/v3"

func RoomRouter(app *fiber.App, roomSvc *Service) {
	group := app.Group("/room")
	group.Post("/join", handleJoinRoom(roomSvc))
	group.Post("/leave", handleLeaveRoom(roomSvc))
}
