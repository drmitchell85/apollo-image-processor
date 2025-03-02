package main

import (
	apiserver "apollo-image-processor/internal/api/server"
	"log"
)

func main() {
	apiServer, err := apiserver.NewServer()
	if err != nil {
		log.Fatalf("API service: error building application: %v", err)
	}

	err = apiServer.Start()
	if err != nil {
		log.Fatalf("API service: failed to start: %v", err)
	}
}
