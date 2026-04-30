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
	// 1. validate
	if roomType == "direct" && len(memberIDs) != 2 {
		return gocql.UUID{}, errors.New("direct room requires exactly 2 members")
	}
	if roomType == "group" && len(memberIDs) < 2 {
		return gocql.UUID{}, errors.New("group room requires at least 2 members")
	}

	// 2. ถ้า direct — เช็คว่ามี room อยู่แล้วหรือยัง
	if roomType == "direct" {
		existingID, err := s.FindDirectRoom(memberIDs[0], memberIDs[1])
		if err == nil {
			return existingID, nil // ✅ คืน room เดิม
		}
	}

	roomID := gocql.TimeUUID()
	err := s.scylla.Query(`INSERT INTO rooms (room_id, type , created_at) VALUES (? ,? , toTimestamp(now()))`,
		roomID, roomType).Exec()
	if err != nil {
		return gocql.UUID{}, err
	}
	creator := memberIDs[0]
	for _, userID := range memberIDs {
		permission := "member"
		if userID == creator {
			permission = "owner"
		}
		err := s.scylla.Query(`

        INSERT INTO room_members (room_id, user_id, permission, joined_at)

        VALUES (?, ?, ?, toTimestamp(now()))`,

			roomID, userID, permission,
		).Exec()

		if err != nil {

			return gocql.UUID{}, err

		}

		// reverse lookup
		err = s.scylla.Query(`
			INSERT INTO user_rooms (user_id, room_id , type) 
			VALUES (?, ?, ?)`,
			userID, roomID, roomType).Exec()
		if err != nil {
			return gocql.UUID{}, err
		}
	}

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
	// remove user first
	if err := s.scylla.Query(
		`DELETE FROM room_members WHERE room_id = ? AND user_id = ?`,
		rID, uID,
	).Exec(); err != nil {
		return err
	}
	// 1. ลบออกจาก room_members

	if err := s.scylla.Query(
		`DELETE FROM room_members WHERE room_id = ? AND user_id = ?`,
		rID, uID,
	).Exec(); err != nil {
		return err
	}
	// 2. ลบจาก user_rooms ด้วย (สำคัญ!)
	if err := s.scylla.Query(
		`DELETE FROM user_rooms WHERE user_id = ? AND room_id = ?`,
		uID, rID,
	).Exec(); err != nil {
		return err
	}
	// 3. เช็ค member ที่เหลือ (แบบเบา ๆ ไม่ใช้ COUNT)
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

	// 4. delete room if <= 1
	if count <= 1 {
		// get remaining users first (for cleanup user_rooms)
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
		// delete all members
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
