package websocket

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Onebluesky882/my-chat-app/internal/chat-service"
)

func (s *Service) join(roomID string, c *client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rooms[roomID]; !ok {
		s.rooms[roomID] = make(map[*client]bool)
	}
	s.rooms[roomID][c] = true
	log.Println("JOIN ROOM:", roomID, "clients:", len(s.rooms[roomID]))
}

func (s *Service) leave(roomID string, c *client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if room, ok := s.rooms[roomID]; ok {
		delete(room, c)
		if len(room) == 0 {
			delete(s.rooms, roomID)
		}
	}
}

func (s *Service) Broadcast(roomID string, payload []byte) {
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

func (s *Service) handleMessage(roomID, userID string, m WSMessage) {
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

func (s *Service) handleTyping(roomID, userID string) {
	payload, _ := json.Marshal(map[string]any{
		"type": "typing",
		"data": map[string]string{
			"user_id": userID,
			"room_id": roomID,
		},
	})
	s.Broadcast(roomID, payload)
}
