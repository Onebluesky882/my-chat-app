package chat

import "github.com/gocql/gocql"

type Message struct {
	RoomID    gocql.UUID `json:"room_id"  `
	MessageID gocql.UUID `json:"message_id"  `
	SenderID  gocql.UUID `json:"sender_id"  `
	Content   string     `json:"content" `
}

type SendRequest struct {
	RoomID       gocql.UUID `json:"room_id"`
	SenderID     gocql.UUID `json:"sender_id"`
	Content      string     `json:"content"`
	Participants []string   `json:"participants"`
}
