package chat

import (
	"context"
	"encoding/json"

	"github.com/gocql/gocql"
	"github.com/redis/go-redis/v9"
)

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

// Pagination
func (s *Service) GetMessagesWithCursor(roomID string, cursor gocql.UUID, limit int) ([]Message, error) {
	var iter *gocql.Iter
	if cursor == (gocql.UUID{}) {
		// first page
		iter = s.scylla.Query(
			`SELECT room_id, message_id , sender_id , content FROM messages WHERE room_id = ? LIMIT ?`,
			roomID, limit,
		).Iter()
	} else {
		iter = s.scylla.Query(
			`
			SELECT room_id, message_id, sender_id, content FROM messages WHERE room_id = ? AND message_id < ? LIMIT ?
			`,
			roomID, cursor, limit,
		).Iter()
	}

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

func (s *Service) GetMessagesWithCache(ctx context.Context, roomID string, cursor gocql.UUID, limit int) ([]Message, error) {

	//1. ถ้าเป็นหน้าแรก → ใช้ Redis
	if cursor == (gocql.UUID{}) {
		msgs, err := s.GetRecentMessage(ctx, roomID)
		if err != nil {
			return nil, err
		}
		if len(msgs) > 0 {
			if len(msgs) > limit {
				msgs = msgs[:limit]
			}
			return msgs, nil
		}
	}
	//  2. fallback → Scylla
	msgs, err := s.GetMessagesWithCursor(roomID, cursor, limit)
	if err != nil {
		return nil, err
	}

	if cursor == (gocql.UUID{}) {
		key := "chat:" + roomID
		for i := len(msgs) - 1; i >= 0; i-- {
			data, _ := json.Marshal(msgs[i])
			s.redis.LPush(ctx, key, data)
		}
		s.redis.LTrim(ctx, key, 0, 99)
	}
	return msgs, nil
}

func (s *Service) GetUnread(ctx context.Context, userID, roomID string) (int, error) {
	key := "unread:" + userID + ":" + roomID

	val, err := s.redis.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (s *Service) MarkAsRead(ctx context.Context, userID, roomID string) error {
	key := "unread:" + userID + ":" + roomID
	return s.redis.Del(ctx, key).Err()
}

