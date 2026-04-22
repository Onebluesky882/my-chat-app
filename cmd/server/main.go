package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Onebluesky882/my-chat-app/internal/cache"
	"github.com/Onebluesky882/my-chat-app/internal/chat-service"
	"github.com/Onebluesky882/my-chat-app/internal/config"
	"github.com/Onebluesky882/my-chat-app/internal/db"
)

func main() {

	cfg := config.LoadConfig()

	// connect db scyllaDB
	scyllaSession, err := db.ConnectScylla()
	if err != nil {
		fmt.Println("DB connection error:", err)
		return
	}
	defer scyllaSession.Close()

	// redis
	rdb, err := cache.NewRedis(cfg)
	if err != nil {
		log.Fatal(err)
	}

	chatSvc := chat.New(scyllaSession, rdb)

	ctx := context.Background()

	msg := chat.Message{
		RoomID:   "room-3",
		SenderID: "user-3",
		Content:  "hello scylla + redis",
	}

	err = chatSvc.Send(ctx, msg)
	if err != nil {
		log.Fatalf("send message error %v", err)
	}
	fmt.Println("✅ message sent")

	// test redis
	err = rdb.Set(ctx, "test:key", "hello-redis", 0).Err()
	if err != nil {
		log.Fatal(err)
	}
	val, err := rdb.Get(ctx, "test:key").Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Redis value:", val)
	fmt.Println("Application started 🚀")
}
