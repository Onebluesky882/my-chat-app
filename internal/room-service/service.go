package room

import (
	"bytes"
	"errors"

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
func (s *Service) GetParticipants(roomID gocql.UUID) ([]gocql.UUID, error) {
	var users []gocql.UUID

	iter := s.scylla.Query(
		`SELECT user_id FROM room_members WHERE room_id = ?`,
		roomID,
	).Iter()

	var user gocql.UUID
	for iter.Scan(&user) {
		users = append(users, user)
	}

	return users, iter.Close()
}

// create room
func (s *Service) CreateRoom(roomType string, memberIDs []gocql.UUID) (gocql.UUID, error) {
	// Step 1: Validate input
	if roomType == "direct" && len(memberIDs) != 2 {
		return gocql.UUID{}, errors.New("direct room requires exactly 2 members")
	}
	if roomType == "group" && len(memberIDs) < 2 {
		return gocql.UUID{}, errors.New("group room requires at least 2 members")
	}

	// Step 2: If direct room — check if already exists (before any insert)
	if roomType == "direct" {
		existingID, err := s.FindDirectRoom(memberIDs[0], memberIDs[1])
		if err == nil {
			// room already exists, return it
			return existingID, nil
		}
		// err != nil means not found, continue to create
	}

	// Step 3: Create the room record
	roomID := gocql.TimeUUID()
	err := s.scylla.Query(
		`INSERT INTO rooms (room_id, type, created_at) VALUES (?, ?, toTimestamp(now()))`,
		roomID, roomType,
	).Exec()
	if err != nil {
		return gocql.UUID{}, err
	}

	// Step 4: Insert each member into room_members + user_rooms (reverse lookup)
	creator := memberIDs[0]
	for _, userID := range memberIDs {
		permission := "member"
		if userID == creator {
			permission = "owner"
		}

		// 4a. Insert into room_members
		err := s.scylla.Query(
			`INSERT INTO room_members (room_id, user_id, permission, joined_at)
			 VALUES (?, ?, ?, toTimestamp(now()))`,
			roomID, userID, permission,
		).Exec()
		if err != nil {
			return gocql.UUID{}, err
		}

		// 4b. Insert into user_rooms (reverse lookup: user → rooms)
		err = s.scylla.Query(
			`INSERT INTO user_rooms (user_id, room_id, type) VALUES (?, ?, ?)`,
			userID, roomID, roomType,
		).Exec()
		if err != nil {
			return gocql.UUID{}, err
		}
	}

	// Step 5: If direct room — register in direct_rooms for future dedup lookup
	if roomType == "direct" {
		a, b := normalizeUsers(memberIDs[0], memberIDs[1])
		err = s.scylla.Query(
			`INSERT INTO direct_rooms (user_a, user_b, room_id) VALUES (?, ?, ?)`,
			a, b, roomID,
		).Exec()
		if err != nil {
			return gocql.UUID{}, err
		}
	}

	// Step 6: Return the new room ID
	return roomID, nil
}
func (s *Service) FindDirectRoom(userA, userB gocql.UUID) (gocql.UUID, error) {
	a, b := normalizeUsers(userA, userB)
	var roomID gocql.UUID
	err := s.scylla.Query(`SELECT room_id FROM direct_rooms WHERE user_a = ? AND user_b = ?`, a, b).Scan(&roomID)
	if err != nil {
		return gocql.UUID{}, err
	}
	return roomID, nil
}

func normalizeUsers(a, b gocql.UUID) (gocql.UUID, gocql.UUID) {
	if bytes.Compare(a.Bytes(), b.Bytes()) < 0 {
		return a, b
	}
	return b, a
}

// invited user to room
func (s *Service) JoinRoom(roomID string, userID string) error {
	rID, err := gocql.ParseUUID(roomID)
	if err != nil {
		return errors.New("invalid room_id")
	}
	uID, err := gocql.ParseUUID(userID)
	if err != nil {
		return errors.New("invalid user_id")
	}

	// 1.room type
	var roomType string
	err = s.scylla.Query(`SELECT type FROM rooms WHERE room_id = ?`, rID).Scan(&roomType)
	if err == gocql.ErrNotFound {
		return errors.New("room not found")
	}
	if err != nil {
		return err
	}

	if roomType == "direct" {
		var count int
		err = s.scylla.Query(`SELECT COUNT(*) FROM room_members WHERE room_id = ?`, rID).Scan(&count)
		if err != nil {
			return err
		}
		if count >= 2 {
			return errors.New("direct room already has 2 members")
		}
	}

	// 3. insert member
	return s.scylla.Query(`
        INSERT INTO room_members (room_id, user_id, joined_at)
        VALUES (?, ?, toTimestamp(now()))`,
		rID, uID,
	).Exec()
}

func (s *Service) LeaveRoom(roomID, userID string) error {

	rID, err := gocql.ParseUUID(roomID)
	if err != nil {
		return errors.New("invalid room_id")
	}
	uID, err := gocql.ParseUUID(userID)
	if err != nil {
		return errors.New("invalid user_id")
	}

	// 1. remove from room_members
	if err := s.scylla.Query(
		`DELETE FROM room_members WHERE room_id = ? AND user_id = ?`,
		rID, uID,
	).Exec(); err != nil {
		return err
	}

	// 2. remove from user_rooms
	if err := s.scylla.Query(
		`DELETE FROM user_rooms WHERE user_id = ? AND room_id = ?`,
		uID, rID,
	).Exec(); err != nil {
		return err
	}

	// 3. check remaining members (ONLY ONCE, LIMIT 2)
	iter := s.scylla.Query(
		`SELECT user_id FROM room_members WHERE room_id = ? LIMIT 2`,
		rID,
	).Iter()

	count := 0
	var tmp gocql.UUID

	for iter.Scan(&tmp) {
		count++
	}

	if err := iter.Close(); err != nil {
		return err
	}

	// 4. delete room if <= 1 member left
	if count <= 1 {

		// cleanup remaining members (best-effort)
		iter := s.scylla.Query(
			`SELECT user_id FROM room_members WHERE room_id = ?`,
			rID,
		).Iter()

		var uid gocql.UUID
		for iter.Scan(&uid) {
			_ = s.scylla.Query(
				`DELETE FROM user_rooms WHERE user_id = ? AND room_id = ?`,
				uid, rID,
			).Exec()
		}
		_ = iter.Close()

		// delete room members
		if err := s.scylla.Query(
			`DELETE FROM room_members WHERE room_id = ?`,
			rID,
		).Exec(); err != nil {
			return err
		}

		// delete room
		if err := s.scylla.Query(
			`DELETE FROM rooms WHERE room_id = ?`,
			rID,
		).Exec(); err != nil {
			return err
		}
	}

	return nil
}
