package room

import "github.com/gocql/gocql"

func (s *Service) IsMember(roomID, userID string) (bool, error) {
	var id string

	err := s.scylla.Query(
		`SELECT user_id FROM room_members WHERE room_id = ? AND user_id = ? LIMIT 1`,
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
