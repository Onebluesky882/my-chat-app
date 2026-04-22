package chat

import "github.com/gocql/gocql"

type Message struct {
	RoomID    string     `json:"room_id"  `
	MessageID gocql.UUID `json:"message_id"  `
	SenderID  string     `json:"sender_id"  `
	Content   string     `json:"content" `
}
