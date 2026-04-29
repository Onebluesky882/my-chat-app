package chat

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Onebluesky882/my-chat-app/internal/room-service"
	"github.com/gocql/gocql"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	scylla *gocql.Session
	redis  *redis.Client
	room   *room.Service
	ws     Broadcaster
}

type Broadcaster interface {
	Broadcast(roomID string, payload []byte)
}
type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type UnreadEvent struct {
	UserID string `json:"user_id"`
	RoomID string `json:"room_id"`
	Count  int64  `json:"count"`
}

func New(s *gocql.Session, r *redis.Client, roomSvc *room.Service) *Service {
	return &Service{
		scylla: s,
		redis:  r,
		room:   roomSvc,
	}

}
func (s *Service) SetBroadcaster(b Broadcaster) {
	s.ws = b
}

func (s *Service) Send(ctx context.Context, msg Message) error {
	msg.MessageID = gocql.TimeUUID()
	// 1. insert into ScyllaDB
	err := s.scylla.Query(
		`
		INSERT INTO messages (room_id, message_id, sender_id, content)
		VALUES (?, ?, ?, ?)
		`,
		msg.RoomID,
		msg.MessageID,
		msg.SenderID,
		msg.Content,
	).Exec()
	if err != nil {
		return err
	}

	data, err := json.Marshal(msg)

	if err != nil {
		return err
	}
	// 2. cache to Redis
	key := "chat:" + msg.RoomID
	s.redis.LPush(ctx, key, data)
	s.redis.LTrim(ctx, key, 0, 99)

	// DB done → broadcast message
	if s.ws != nil {
		s.ws.Broadcast(msg.RoomID, data)
	}

	// ดึง participants จาก db
	participants, err := s.room.GetParticipants(msg.RoomID)
	if err != nil {
		return err
	}

	// 4. unread logic
	for _, userID := range participants {
		if userID == msg.SenderID {
			continue
		}
		unreadKey := "unread:" + userID + ":" + msg.RoomID
		count, err := s.redis.Incr(ctx, unreadKey).Result()
		if err != nil {
			return err
		}

		fmt.Println("increase unread for:", userID, count)

		// 5. broadcast unread event
		event := Event{
			Type: "unread",
			Data: UnreadEvent{
				UserID: userID,
				RoomID: msg.RoomID,
				Count:  count,
			},
		}
		eventData, _ := json.Marshal(event)
		s.ws.Broadcast(msg.RoomID, eventData)
	}

	return nil
}

// get Messages
func (s *Service) GetMessages(ctx context.Context, roomID string, limit int) ([]Message, error) {
	// 1. try Redis (don't fail if Redis error)

	msgs, err := s.GetRecentMessage(ctx, roomID)
	if err == nil && len(msgs) > 0 {
		if len(msgs) > limit {
			msgs = msgs[:limit]
		}
		return msgs, nil
	}
	// 2. fallback Scylla
	msgs, err = s.GetMessagesFromDB(roomID, limit)
	if err != nil {
		return nil, err
	}

	// 3. cache back to Redis

	if len(msgs) > 0 {
		key := "chat:" + roomID
		for i := len(msgs) - 1; i >= 0; i-- {
			data, _ := json.Marshal(msgs[i])
			s.redis.LPush(ctx, key, data)
		}
		s.redis.LTrim(ctx, key, 0, 99)
	}
	return msgs, nil
}
