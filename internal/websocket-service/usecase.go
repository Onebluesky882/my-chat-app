package websocket

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Onebluesky882/my-chat-app/internal/chat-service"
	"github.com/gocql/gocql"
)

func (s *Service) join(roomID gocql.UUID, c *client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rooms[roomID]; !ok {
		s.rooms[roomID] = make(map[*client]bool)
	}
	s.rooms[roomID][c] = true
	log.Println("JOIN ROOM:", roomID, "clients:", len(s.rooms[roomID]))
}

func (s *Service) leave(roomID gocql.UUID, c *client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if room, ok := s.rooms[roomID]; ok {
		delete(room, c)
		if len(room) == 0 {
			delete(s.rooms, roomID)
		}
	}
}

func (s *Service) Broadcast(roomID gocql.UUID, payload []byte) {
	s.mu.RLock()
	clients := make([]*client, 0, len(s.rooms[roomID]))
	for c := range s.rooms[roomID] {
		clients = append(clients, c)
	}
	s.mu.RUnlock()

	var dead []*client
	for _, c := range clients {
		if err := s.safeWrite(c, payload); err != nil {
			log.Println("broadcast error:", err)
			dead = append(dead, c)
		}
	}

	if len(dead) > 0 {
		s.mu.Lock()
		for _, c := range dead {
			delete(s.rooms[roomID], c)
		}
		s.mu.Unlock()
	}
}
func (s *Service) BroadcastUser(userID gocql.UUID, payload []byte) {
	s.mu.RLock()
	c, ok := s.userConns[userID]
	s.mu.RUnlock()
	if !ok {
		return // user offline — ไม่เป็นไร
	}
	if err := s.safeWrite(c, payload); err != nil {
		log.Println("BroadcastUser error:", err)
	}
}

func (s *Service) handleMessage(roomID, userID gocql.UUID, m WSMessage) {
	if m.Content == "" {
		return
	}
	msg := chat.Message{
		RoomID:   roomID,
		SenderID: userID,
		Content:  m.Content,
	}
	if err := s.chatSvc.Send(context.Background(), msg); err != nil {
		log.Println("send message error:", err)
		return
	}
	payload, _ := json.Marshal(map[string]any{
		"type": "message",
		"data": msg,
	})
	s.Broadcast(roomID, payload)
}

func (s *Service) handleTyping(roomID, userID gocql.UUID) {
	payload, _ := json.Marshal(map[string]any{
		"type": "typing",
		"data": map[string]string{
			"user_id": userID.String(),
			"room_id": roomID.String(),
		},
	})
	s.Broadcast(roomID, payload)
}
