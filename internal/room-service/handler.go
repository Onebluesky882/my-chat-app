package room

import (
	"errors"

	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v3"
)

var (
	ErrRoomNotFound          = errors.New("room not found")
	ErrRoomFull              = errors.New("direct room is full")
	ErrInvalidRoom           = errors.New("invalid room_id")
	ErrDirectRequires2       = errors.New("direct room requires exactly 2 members")
	ErrGroupRequiresAtLeast2 = errors.New("group room requires at least 2 members")
)

func handleCreateRoom(roomSvc *Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req struct {
			Type    string   `json:"type"`
			Members []string `json:"members"`
		}
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if req.Type == "" || len(req.Members) == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "type and members required"})
		}

		memberIDs := make([]gocql.UUID, 0, len(req.Members))
		for _, member := range req.Members {
			uuid, err := gocql.ParseUUID(member)
			if err != nil {
				return c.Status(400).JSON(fiber.Map{"error": "invalid member UUID"})
			}
			memberIDs = append(memberIDs, uuid)
		}

		roomID, err := roomSvc.CreateRoom(req.Type, memberIDs)
		if err != nil {
			switch {
			case errors.Is(err, ErrDirectRequires2),
				errors.Is(err, ErrGroupRequiresAtLeast2):
				return c.Status(400).JSON(fiber.Map{"error": err.Error()})
			default:
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
		}

		return c.Status(201).JSON(fiber.Map{
			"room_id": roomID,
			"type":    req.Type,
		})
	}
}
func handleJoinRoom(roomSvc *Service) fiber.Handler {
	// userId require from auth who create
	return func(c fiber.Ctx) error {
		var req struct {
			RoomID string `json:"room_id"`
			UserID string `json:"user_id"`
		}
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if req.RoomID == "" || req.UserID == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "room_id and user_id required",
			})
		}

		err := roomSvc.JoinRoom(req.RoomID, req.UserID)
		if err != nil {

			switch {
			case errors.Is(err, ErrRoomNotFound):
				return c.Status(404).JSON(fiber.Map{"error": err.Error()})
			case errors.Is(err, ErrRoomFull):
				return c.Status(409).JSON(fiber.Map{"error": err.Error()})
			default:
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}

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
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if req.RoomID == "" || req.UserID == "" {
			return c.Status(400).JSON(fiber.Map{"error": "room_id and user_id required"})
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
