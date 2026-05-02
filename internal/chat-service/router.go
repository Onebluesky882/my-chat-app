package chat

import (
	"github.com/Onebluesky882/my-chat-app/internal/room-service"
	"github.com/gofiber/fiber/v3"
)

func ChatRouter(app *fiber.App, chatSvc *Service, roomSvc *room.Service) {
	chatGroup := app.Group("/chat")
	chatGroup.Post("/send", handleSendMessage(chatSvc))
	chatGroup.Get("/unread", handleGetUnread(chatSvc))
	chatGroup.Post("/read", handleMarkAsRead(chatSvc))
	chatGroup.Get("/message", handleGetMessage(chatSvc, roomSvc))
}

/*

POST /chat/send     → message ถูก save ลง DB ไหม
GET  /chat/message  → ดึงออกมาได้ไหม
GET  /chat/unread   → count ถูกต้องไหม
POST /chat/read     → reset count ได้ไหม

*/
