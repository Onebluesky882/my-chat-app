package room

import (
	"errors"

	"github.com/gocql/gocql"
)

func (s *Service) IsMember(roomID, userID string) (bool, error) {
	rID, err := gocql.ParseUUID(roomID)
	if err != nil {
		return false, errors.New("invalid room_id")
	}
	uID, err := gocql.ParseUUID(userID)
	if err != nil {
		return false, errors.New("invalid user_id")
	}

	var id gocql.UUID

	err = s.scylla.Query(
		`SELECT user_id FROM room_members WHERE room_id = ? AND user_id = ?`,
		rID, uID,
	).Scan(&id)

	if err == gocql.ErrNotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}
