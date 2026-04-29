package websocket

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/Onebluesky882/my-chat-app/internal/chat-service"
	"github.com/Onebluesky882/my-chat-app/internal/room-service"
	"github.com/gorilla/websocket"
)

type Service struct {
	Upgrader websocket.Upgrader
	chatSvc  *chat.Service
	roomSvc  *room.Service
	mu       sync.RWMutex
	rooms    map[string]map[*websocket.Conn]bool
}

func New(
	chatSvc *chat.Service,
	roomSvc *room.Service,
) *Service {
	return &Service{
		chatSvc: chatSvc,
		roomSvc: roomSvc,
		rooms:   make(map[string]map[*websocket.Conn]bool),
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (s *Service) ConnectUser(w http.ResponseWriter, r *http.Request, roomID string, userID string) error {
	ok, err := s.roomSvc.IsMember(

		roomID,
		userID,
	)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("forbidden")
	}
	conn, err := s.Upgrader.Upgrade(
		w, r, nil,
	)
	if err != nil {
		return err
	}
	s.join(roomID, conn)
	defer s.leave(roomID, conn)
	log.Println(
		"WS connected:",
		userID,
		roomID,
	)
	return s.readLoop(
		roomID,
		userID,
		conn,
	)

}
func (s *Service) readLoop(
	roomID string,
	userID string,
	conn *websocket.Conn,
) error {

	for {

		_, msg, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		log.Println(
			"incoming:",
			string(msg),
		)
	}
}
