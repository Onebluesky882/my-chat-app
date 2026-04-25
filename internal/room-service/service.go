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

func (s *Service) JoinRoom(roomID string, userID string) error {
	return s.scylla.Query(
		`INSERT INTO room_members (room_id, user_id) VALUES(? ,? )`,
		roomID, userID,
	).Exec()
}

func (s *Service) LeaveRoom(roomID, userID string) error {
	// remove user first
	if err := s.scylla.Query(
		`DELETE FROM room_members WHERE room_id = ? AND user_id = ?`,
		roomID, userID,
	).Exec(); err != nil {
		return err
	}

	var memberCount int
	if err := s.scylla.Query(
		`SELECT COUNT(*) FROM room_members WHERE room_id = ?`,
		roomID,
	).Scan(&memberCount); err != nil {
		return err
	}

	if memberCount <= 1 {
		if err := s.scylla.Query(
			`DELETE FROM room_members WHERE room_id = ?`,
			roomID,
		).Exec(); err != nil {
			return err
		}

		// ignore if rooms table doesn't exist
		if err := s.scylla.Query(
			`DELETE FROM rooms WHERE room_id = ?`,
			roomID,
		).Exec(); err != nil {
			// log only, don't fail
			return nil
		}
	}

	return nil
}