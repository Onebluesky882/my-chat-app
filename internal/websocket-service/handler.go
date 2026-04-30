package websocket

import (
	"encoding/json"

	"github.com/Onebluesky882/my-chat-app/internal/chat-service"
	"github.com/fasthttp/websocket"
	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v3"
)

func HandleWS(wsSvc *Service, roomID, userID gocql.UUID) func(*websocket.Conn) {
	return func(conn *websocket.Conn) {
		// ✅ ไม่ต้อง conn.Close() — Connect → close(c.send) → writePump ปิดเอง
		wsSvc.Connect(conn, roomID, userID)
	}
}

func HandleSendToUser(wsSvc *Service, chatSvc *chat.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req struct {
			UserID  string `json:"user_id"`
			RoomID  string `json:"room_id"` // ✅ เพิ่ม room_id
			Message string `json:"message"`
		}
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		rID, err := gocql.ParseUUID(req.RoomID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid room_id",
			})
		}

		uID, err := gocql.ParseUUID(req.UserID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid user_id",
			})
		}

		// ✅ บันทึก DB ก่อน
		msg := chat.Message{
			RoomID:   rID,
			SenderID: uID,
			Content:  req.Message,
		}
		if err := chatSvc.Send(c.Context(), msg); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		// ✅ ส่งผ่าน WebSocket หลัง save สำเร็จ
		payload, _ := json.Marshal(map[string]any{
			"type": "message",
			"data": msg,
		})
		wsSvc.SendToUser(uID, payload)

		return c.JSON(fiber.Map{"status": "sent"})
	}
}
