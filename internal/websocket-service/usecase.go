package websocket

import "github.com/gorilla/websocket"

func (s *Service) join(roomID string, conn *websocket.Conn) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.rooms[roomID]; !ok {
		s.rooms[roomID] = make(map[*websocket.Conn]bool)
	}
	s.rooms[roomID][conn] = true

	return nil
}

func (s *Service) leave(roomID string, conn *websocket.Conn) {

	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.rooms[roomID], conn)
	conn.Close()

}

func (s *Service) Broadcast(roomID string, payload []byte) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	for conn := range s.rooms[roomID] {
		conn.WriteMessage(
			websocket.TextMessage,
			payload,
		)
	}
}
