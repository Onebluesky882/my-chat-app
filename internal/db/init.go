package db

import (
	"github.com/uptrace/bun"
)

var (
	DBCon *bun.DB
)

func InitDb() {
	// 	newDB := `USE chat;

	// CREATE TABLE messages (
	//   room_id text,
	//   message_id timeuuid,
	//   user_id text,
	//   content text,
	//   PRIMARY KEY (room_id, message_id)
	// ) WITH CLUSTERING ORDER BY (message_id DESC);`
	// fmt.Println(newDB)
}



