package websocket

import (
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/Onebluesky882/my-chat-app/internal/chat-service"
	"github.com/Onebluesky882/my-chat-app/internal/room-service"
	"github.com/fasthttp/websocket"
	"github.com/gocql/gocql"
)

type Service struct {
	chatSvc   *chat.Service
	roomSvc   *room.Service
	mu        sync.RWMutex
	userConns map[gocql.UUID]*client
	rooms     map[gocql.UUID]map[*client]bool
}

func New(chatSvc *chat.Service, roomSvc *room.Service) *Service {
	return &Service{
		chatSvc:   chatSvc,
		roomSvc:   roomSvc,
		userConns: make(map[gocql.UUID]*client),
		rooms:     make(map[gocql.UUID]map[*client]bool),
	}
}

func newClient(conn *websocket.Conn) *client {
	return &client{
		conn: conn,
		send: make(chan []byte, 256),
	}
}

func (c *client) writePump() {
	defer c.conn.Close()
	for payload := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			log.Println("writePump error:", err)
			return
		}
	}
}

func (s *Service) safeWrite(c *client, payload []byte) error {
	select {
	case c.send <- payload:
		return nil
	default:
		return errors.New("client send buffer full")
	}
}

func (s *Service) Connect(conn *websocket.Conn, roomID, userID gocql.UUID) error {

	ok, err := s.roomSvc.IsMember(roomID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("forbidden")
	}

	c := newClient(conn)
	go c.writePump()

	s.mu.Lock()
	s.userConns[userID] = c
	s.mu.Unlock()

	s.join(roomID, c)

	defer func() {
		s.leave(roomID, c)
		s.mu.Lock()
		delete(s.userConns, userID)
		s.mu.Unlock()
		close(c.send) // → writePump จบ → conn.Close() อัตโนมัติ
	}()

	return s.readLoop(roomID, userID, c)
}

func (s *Service) SendToUser(userID gocql.UUID, payload []byte) {
	s.mu.RLock()
	c, ok := s.userConns[userID]
	s.mu.RUnlock()
	if !ok {
		return
	}
	if err := s.safeWrite(c, payload); err != nil {
		log.Println("send to user error:", err)
	}
}

func (s *Service) readLoop(roomID, userID gocql.UUID, c *client) error {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("WS READ ERROR:", err)
			return err
		}
		log.Println("INCOMING:", userID, string(msg))

		var m WSMessage
		if err := json.Unmarshal(msg, &m); err != nil {
			log.Println("invalid json:", err)
			continue
		}
		if m.Type == "" {
			continue
		}

		switch m.Type {
		case "message":
			s.handleMessage(roomID, userID, m)
		case "typing":
			s.handleTyping(roomID, userID)
		case "leave":
			return nil
		}
	}
}
