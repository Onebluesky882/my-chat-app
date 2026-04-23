package chat

import "github.com/gofiber/fiber/v3"

func ChatRouter(app *fiber.App, chatSvc *Service) {
	chatGroup := app.Group("/chat")
	chatGroup.Post("/send", SendMessage(chatSvc))
	chatGroup.Get("/unread", GetUnread(chatSvc))
	chatGroup.Post("/read", MarkAsRead(chatSvc))
}
