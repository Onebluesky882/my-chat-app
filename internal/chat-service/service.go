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
}

func New(s *gocql.Session, r *redis.Client, roomSvc *room.Service) *Service {
	return &Service{
		scylla: s,
		redis:  r,
		room:   roomSvc,
	}
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

	// 2. cache to Redis
	data, _ := json.Marshal(msg)
	key := "chat:" + msg.RoomID
	s.redis.LPush(ctx, key, data)
	s.redis.LTrim(ctx, key, 0, 99)

	// ดึง participants จาก db
	participants, err := s.room.GetParticipants(msg.RoomID)
	if err != nil {
		return err
	}
	fmt.Println("participants:", participants)
	// 4. unread logic
	for _, userID := range participants {
		if userID == msg.SenderID {
			continue
		}
		unreadKey := "unread" + userID + ":" + msg.RoomID
		s.redis.Incr(ctx, unreadKey)
	}

	return nil
}

// get Messages
func (s *Service) GetMessages(ctx context.Context, roomID string) ([]Message, error) {
	// 1. try Redis
	msgs, err := s.GetRecentMessage(ctx, roomID)
	if err != nil {
		return nil, err
	}

	if len(msgs) > 0 {
		return msgs, nil
	}
	// 2. fallback Scylla
	msgs, err = s.GetMessagesFromDB(roomID, 50)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
