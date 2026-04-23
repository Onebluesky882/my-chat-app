package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Onebluesky882/my-chat-app/internal/cache"
	"github.com/Onebluesky882/my-chat-app/internal/chat-service"
	"github.com/Onebluesky882/my-chat-app/internal/config"
	"github.com/Onebluesky882/my-chat-app/internal/db"
	"github.com/Onebluesky882/my-chat-app/internal/room-service"
	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()
	cfg := config.LoadConfig()

	// connect db scyllaDB
	scyllaSession, err := db.ConnectScylla()
	if err != nil {
		fmt.Println("DB connection error:", err)
		return
	}
	defer scyllaSession.Close()

	err = db.CreateTables(scyllaSession)
	if err != nil {
		log.Fatal("create table error:", err)
	}

	// redis
	rdb, err := cache.NewRedis(cfg)
	if err != nil {
		log.Fatal(err)
	}
	roomSvc := room.New(scyllaSession)
	chatSvc := chat.New(scyllaSession, rdb, roomSvc)

	ctx := context.Background()

	chat.ChatRouter(app, chatSvc)

	// test redis
	msgs, err := chatSvc.GetMessages(ctx, "room-1")
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range msgs {
		fmt.Println(m)
	}
	// run server
	log.Println("Server running on :3000 🚀")
	log.Fatal(app.Listen(":3000"))
}
