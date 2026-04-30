package room

import (
	"github.com/gocql/gocql"
)

func (s *Service) IsMember(roomID, userID gocql.UUID) (bool, error) {

	var id gocql.UUID

	err := s.scylla.Query(
		`SELECT user_id FROM room_members WHERE room_id = ? AND user_id = ?`,
		roomID, userID,
	).Scan(&id)

	if err == gocql.ErrNotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}
