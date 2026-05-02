package room

import (
	"time"

	"github.com/gocql/gocql"
)

type Room struct {
	ID        gocql.UUID `json:"id"`
	Type      string     `json:"type"` // "direct" | "group"
	CreatedAt time.Time  `json:"created_at"`
	Name      *string    `json:"name"`
}

type RoomMember struct {
	RoomID     gocql.UUID `json:"room_id"`
	UserID     gocql.UUID `json:"user_id"`
	Permission string     `json:"permission"` // "owner" | "admin" | "member"
	JoinedAt   time.Time  `json:"joined_at"`
}

type UserRoom struct {
	UserID   gocql.UUID `json:"user_id"`
	RoomID   gocql.UUID `json:"room_id"`
	LastRead time.Time  `json:"last_read"`
	Type     string     `json:"type"`
}
