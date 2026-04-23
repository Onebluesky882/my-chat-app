package chat

import (
	"strconv"

	"github.com/Onebluesky882/my-chat-app/internal/room-service"
	"github.com/gofiber/fiber/v3"
)

func handleGetMessage(chatSvc *Service, roomSvc *room.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		limitStr := c.Query("limit", "50")
		// todo change to jwt userId
		userID := c.Query("user_id")
		roomID := c.Query("room_id")
		if userID == "" || roomID == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "user_id and room_id required",
			})
		}
		ok, err := roomSvc.IsMember(roomID, userID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if !ok {
			return c.Status(403).JSON(fiber.Map{
				"error": "you are not in this room",
			})
		}
		limit, err := strconv.Atoi(limitStr)

		if err != nil {
			limit = 50
		}

		msgs, err := chatSvc.GetMessages(c.Context(), roomID, limit)
		if err != nil {

			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(msgs)
	}
}

func handleSendMessage(chatSvc *Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req SendRequest

		// parse JSON
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid request",
			})
		}
		// validate เบื้องต้น
		if req.RoomID == "" || req.SenderID == "" || req.Content == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "missing fields",
			})
		}

		// call service
		err := chatSvc.Send(c.Context(), Message{
			RoomID:   req.RoomID,
			SenderID: req.SenderID,
			Content:  req.Content,
		})

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	}
}
func handleGetUnread(chatSvc *Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := c.Query("user_id")
		roomID := c.Query("room_id")

		if userID == "" || roomID == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "user_id and room_id are required",
			})
		}

		count, err := chatSvc.GetUnread(c.Context(), userID, roomID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"user_id": userID,
			"room_id": roomID,
			"unread":  count,
		})
	}
}

func handleMarkAsRead(chatSvc *Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := c.Query("user_id")
		roomID := c.Query("room_id")
		// validate
		if userID == "" || roomID == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "user_id and room_id are required",
			})
		}

		err := chatSvc.MarkAsRead(c.Context(), userID, roomID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"status":  "read",
			"user_id": userID,
			"room_id": roomID,
		})
	}
}
