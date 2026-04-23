package room

import "github.com/gofiber/fiber/v3"

func handleJoinRoom(roomSvc *Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req struct {
			RoomID string `json:"room_id"`
			UserID string `json:"user_id"`
		}
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if req.RoomID == "" || req.UserID == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "room_id and user_id required",
			})
		}

		err := roomSvc.JoinRoom(req.RoomID, req.UserID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{
			"status": "joined",
		})
	}
}

func handleLeaveRoom(roomSvc *Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req struct {
			RoomID string `json:"room_id"`
			UserID string `json:"user_id"`
		}
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if req.RoomID == "" || req.UserID == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "room_id and user_id required",
			})
		}

		err := roomSvc.LeaveRoom(req.RoomID, req.UserID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{
			"status":  "left",
			"room_id": req.RoomID,
			"user_id": req.UserID,
		})
	}
}
