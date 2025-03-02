package main

import (
	"apollo-image-processor/internal/servers/apiserver"
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
