package websocket

import "github.com/fasthttp/websocket"

type WSMessage struct {
	Type    string `json:"type"`
	RoomID  string `json:"room_id"`
	Content string `json:"content"`
}

type client struct {
	conn *websocket.Conn
	send chan []byte
}
