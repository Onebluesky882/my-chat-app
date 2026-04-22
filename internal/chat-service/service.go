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

func (s *Service) GetRecentMessage(ctx context.Context, roomID string) ([]Message, error) {
	key := "chat:" + roomID

	data, err := s.redis.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var messages []Message

	for _, item := range data {
		var msg Message
		if err := json.Unmarshal([]byte(item), &msg); err != nil {
			continue
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// Fallback ไป Scylla
func (s *Service) GetMessagesFromDB(roomID string, limit int) ([]Message, error) {
	iter := s.scylla.Query(
		`
		SELECT room_id , message_id , sender_id , content FROM messages WHERE room_id = ? LIMIT ?
		`, roomID, limit,
	).Iter()

	var messages []Message
	var m Message

	for iter.Scan(&m.RoomID, &m.MessageID, &m.SenderID, &m.Content) {
		messages = append(messages, m)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return messages, nil
}

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
