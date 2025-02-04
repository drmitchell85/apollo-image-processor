package main

import (
	"apollo-image-processor/internal/servers/apiserver"
	"log"
)

func main() {
	_, err := apiserver.NewServer()
	if err != nil {
		log.Fatalf("API service: error building application: %v", err)
	}
}
