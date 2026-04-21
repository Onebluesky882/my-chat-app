package db

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
)

func ConnectScylla() (*gocql.Session, error) {

	_ = godotenv.Load()

	hosts := os.Getenv("SCYLLA_HOSTS")
	keyspace := os.Getenv("SCYLLA_KEYSPACE")

	fmt.Println("hosts:", hosts)
	fmt.Println("keyspace:", keyspace)

	if hosts == "" {
		return nil, fmt.Errorf("SCYLLA_HOSTS not set")
	}

	// create cluster
	hostList := strings.Split(hosts, ",")
	cluster := gocql.NewCluster(hostList...)

	cluster.Port = 9042
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum

	cluster.Timeout = 20 * time.Second
	cluster.ConnectTimeout = 20 * time.Second
	cluster.NumConns = 1
	cluster.ProtoVersion = 4
	cluster.DisableInitialHostLookup = true
	fmt.Println("creating session...")

	session, err := cluster.CreateSession()

	fmt.Println("session created")
	if err != nil {
		return nil, fmt.Errorf("scylla connect error: %w", err)
	}
	// insert query
	roomID := "room2"
	userID := "user2"
	content := "hello wansing"

	err = session.Query(`
		INSERT INTO messages (room_id, message_id, user_id, content, created_at)
		VALUES (?, now(), ?, ?, toTimestamp(now()))
	`, roomID, userID, content).Exec()

	if err != nil {
		log.Fatal("insert error:", err)
	}

	fmt.Println("✅ Connected to ScyllaDB")

	return session, nil
}
