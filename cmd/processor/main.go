package main

import (
	processorserver "apollo-image-processor/internal/processor/server"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	processorServer, err := processorserver.NewServer()
	if err != nil {
		log.Fatalf("Processor Service: error building application: %v", err)
	}

	errc := make(chan error)
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		err = processorServer.Start()
		if err != nil {
			log.Fatalf("API service: failed to start: %v", err)
			errc <- err
		}
	}()

	select {
	case <-sigc:
		log.Println("received signal to shut down...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		processorServer.Shutdown(ctx)
		cancel()

	case <-errc:
		log.Println("error starting up server, shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		processorServer.Shutdown(ctx)
		cancel()
	}

	os.Exit(0)

}
