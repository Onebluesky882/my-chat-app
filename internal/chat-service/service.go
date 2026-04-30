package chat

import (
	"context"
	"encoding/json"

	"github.com/Onebluesky882/my-chat-app/internal/room-service"
	"github.com/gocql/gocql"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	scylla *gocql.Session
	redis  *redis.Client
	room   *room.Service
	ws     RoomBroadcaster
	userWs UserBroadcaster
}

type RoomBroadcaster interface {
	Broadcast(roomID gocql.UUID, payload []byte)
}

type UserBroadcaster interface {
	BroadcastUser(userID gocql.UUID, payload []byte)
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

// setters
func (s *Service) SetUserBroadcaster(b UserBroadcaster) {
	s.userWs = b
}

func (s *Service) SetBroadcaster(b RoomBroadcaster) {
	s.ws = b
}

func (s *Service) Send(ctx context.Context, msg Message) error {

	msg.MessageID = gocql.TimeUUID()

	// 1. DB insert
	if err := s.scylla.Query(`
		INSERT INTO messages (room_id, message_id, sender_id, content)
		VALUES (?, ?, ?, ?)
	`,
		msg.RoomID,
		msg.MessageID,
		msg.SenderID,
		msg.Content,
	).Exec(); err != nil {
		return err
	}

	// 2. marshal
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// 3. cache redis
	key := "chat:" + msg.RoomID.String()

	if err := s.redis.LPush(ctx, key, data).Err(); err != nil {
		return err
	}
	if err := s.redis.LTrim(ctx, key, 0, 99).Err(); err != nil {
		return err
	}

	// 4. broadcast room
	if s.ws != nil {
		s.ws.Broadcast(msg.RoomID, data)
	}

	// 5. participants
	participants, err := s.room.GetParticipants(msg.RoomID)
	if err != nil {
		return err
	}

	// 6. unread logic
	for _, userID := range participants {

		if userID == msg.SenderID {
			continue
		}

		unreadKey := "unread:" + userID.String() + ":" + msg.RoomID.String()

		count, err := s.redis.Incr(ctx, unreadKey).Result()
		if err != nil {
			return err
		}

		event := Event{
			Type: "unread",
			Data: UnreadEvent{
				UserID: userID.String(),
				RoomID: msg.RoomID.String(),
				Count:  count,
			},
		}

		eventData, err := json.Marshal(event)
		if err != nil {
			return err
		}

		// 7. user broadcast (IMPORTANT FIX)
		if s.userWs != nil {
			s.userWs.BroadcastUser(userID, eventData)
		}
	}

	return nil
}

// =========================
// GET MESSAGES
// =========================

func (s *Service) GetMessages(ctx context.Context, roomID gocql.UUID, limit int) ([]Message, error) {

	msgs, err := s.GetRecentMessage(ctx, roomID)
	if err == nil && len(msgs) > 0 {
		if len(msgs) > limit {
			msgs = msgs[:limit]
		}
		return msgs, nil
	}

	msgs, err = s.GetMessagesFromDB(roomID, limit)
	if err != nil {
		return nil, err
	}

	if len(msgs) > 0 {
		key := "chat:" + roomID.String()

		for i := len(msgs) - 1; i >= 0; i-- {
			data, err := json.Marshal(msgs[i])
			if err != nil {
				return nil, err
			}
			s.redis.LPush(ctx, key, data)
		}

		s.redis.LTrim(ctx, key, 0, 99)
	}

	return msgs, nil
}
