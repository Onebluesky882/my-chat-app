package chat

import (
	"context"
	"encoding/json"

	"github.com/gocql/gocql"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	scylla *gocql.Session
	redis  *redis.Client
}

type Message struct {
	RoomID    string     `json:"room_id"  `
	MessageID gocql.UUID `json:"message_id"  `
	SenderID  string     `json:"sender_id"  `
	Content   string     `json:"content" `
}

func New(s *gocql.Session, r *redis.Client) *Service {
	return &Service{s, r}
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
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	key := "chat:" + msg.RoomID

	err = s.redis.LPush(ctx, key, data).Err()
	if err != nil {
		return err
	}

	err = s.redis.LTrim(ctx, key, 0, 99).Err()
	if err != nil {
		return err
	}

	return nil
}
