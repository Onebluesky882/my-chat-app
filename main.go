package main

import (
	"fmt"

	"github.com/Onebluesky882/my-chat-app/internal/db"
)

func main() {

	session, err := db.ConnectScylla()
	if err != nil {
		fmt.Println("DB connection error:", err)
		return
	}

	defer session.Close()

	fmt.Println("Application started 🚀")
}
