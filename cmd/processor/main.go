package main

import (
	"apollo-image-processor/internal/servers/processorserver"
	"log"
)

func main() {
	_, err := processorserver.NewServer()
	if err != nil {
		log.Fatalf("Processor Service: error building application: %v", err)
	}
}
