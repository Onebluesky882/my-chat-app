package chat

import (
	"context"
	"testing"

	"github.com/Onebluesky882/my-chat-app/internal/cache"
	"github.com/Onebluesky882/my-chat-app/internal/config"
	"github.com/Onebluesky882/my-chat-app/internal/db"
	"github.com/gocql/gocql"
)

func UUIDZero() gocql.UUID {
	return gocql.UUID{}
}

func TestGetMessagesWithCache_MissThenHit(t *testing.T) {
	ctx := context.Background()
	// load config
	cfg := config.LoadConfig()

	// connect scylla
	scylla, err := db.ConnectScylla()
	if err != nil {
		t.Fatalf("scylla error: %v", err)
	}
	defer scylla.Close()

	// connect redis
	rdb, err := cache.NewRedis(cfg)
	if err != nil {
		t.Fatalf("redis error: %v", err)
	}

	// todo nil
	svc := New(scylla, rdb, nil)

	roomID := "room-test"

	// 🔥 clear redis
	rdb.Del(ctx, "chat:"+roomID)

	// 🔥 insert test messages (DB + Redis)
	for i := 0; i < 5; i++ {
		err := svc.Send(ctx, Message{
			RoomID:   roomID,
			SenderID: "user",
			Content:  "msg",
		})
		if err != nil {
			t.Fatalf("send error: %v", err)
		}
	}

	// ❗ clear redis อีกครั้ง เพื่อ force MISS
	rdb.Del(ctx, "chat:"+roomID)

	// ===== FIRST CALL (MISS) =====
	msgs, err := svc.GetMessagesWithCache(ctx, roomID, UUIDZero(), 3)
	if err != nil {
		t.Fatalf("get messages error: %v", err)
	}

	if len(msgs) == 0 {
		t.Fatal("expected messages, got empty")
	}

	// ===== SECOND CALL (HIT) =====
	msgs2, err := svc.GetMessagesWithCache(ctx, roomID, UUIDZero(), 3)
	if err != nil {
		t.Fatalf("get messages error: %v", err)
	}

	if len(msgs2) == 0 {
		t.Fatal("expected cached messages")
	}
}
