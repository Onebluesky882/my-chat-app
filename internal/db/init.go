package db

import (
	"github.com/gocql/gocql"
)

func CreateTables(session *gocql.Session) error {
	query := `
	CREATE TABLE IF NOT EXISTS room_members (
		room_id text,
		user_id text,
		PRIMARY KEY (room_id, user_id)
	);`

	return session.Query(query).Exec()
}
