package main

import (
	"log"

	"app-server/internal/server"
)

func main() {
	srv := server.New()
	
	log.Println("Starting server on :8080")
	if err := srv.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}