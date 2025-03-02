package main

import (
	processorserver "apollo-image-processor/internal/processor/server"
	"log"
)

func main() {
	processorServer, err := processorserver.NewServer()
	if err != nil {
		log.Fatalf("Processor Service: error building application: %v", err)
	}

	err = processorServer.Start()
	if err != nil {
		log.Fatalf("API service: failed to start: %v", err)
	}
}
