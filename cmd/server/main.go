package main

import (
	"fmt"
	"log"

	"github.com/Onebluesky882/my-chat-app/internal/cache"
	"github.com/Onebluesky882/my-chat-app/internal/chat-service"
	"github.com/Onebluesky882/my-chat-app/internal/config"
	"github.com/Onebluesky882/my-chat-app/internal/db"
	"github.com/Onebluesky882/my-chat-app/internal/room-service"
	"github.com/Onebluesky882/my-chat-app/internal/websocket-service"
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

	// redis
	rdb, err := cache.NewRedis(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// instant
	roomSvc := room.New(scyllaSession)
	chatSvc := chat.New(scyllaSession, rdb, roomSvc)
	// routers
	chat.ChatRouter(app, chatSvc, roomSvc)
	room.RoomRouter(app, roomSvc)

	// websocket
	wsSvc := websocket.New(chatSvc, roomSvc)
	websocket.WSRouter(app, wsSvc, chatSvc)
	chatSvc.SetBroadcaster(wsSvc)
	chatSvc.SetUserBroadcaster(wsSvc)

	// run server
	log.Println("Server running on :3000 🚀")
	log.Fatal(app.Listen(":3000"))
}
