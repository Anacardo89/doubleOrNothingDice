package main

import (
	"log"

	"github.com/Anacardo89/doubleOrNothingDice/internal/server"
)

func main() {
	server := server.NewServer()
	log.Println("Starting the server...")
	server.Run(":8080")
}
