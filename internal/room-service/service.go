package room

import (
	"github.com/gocql/gocql"
)

type Service struct {
	scylla *gocql.Session
}

func New(session *gocql.Session) *Service {
	return &Service{
		scylla: session,
	}
}

func (s *Service) GetParticipants(roomID string) ([]string, error) {
	var users []string

	iter := s.scylla.Query(
		`SELECT user_id FROM room_members WHERE room_id = ?`,
		roomID,
	).Iter()

	var user string
	for iter.Scan(&user) {
		users = append(users, user)
	}

	return users, iter.Close()
}
