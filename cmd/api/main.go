package main

import (
	"apollo-image-processor/internal/server"
	"log"
)

func main() {
	_, err := server.NewServer()
	if err != nil {
		log.Fatalf("error building application: %v", err)
	}
}
