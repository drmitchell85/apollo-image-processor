package main

import (
	apiserver "apollo-image-processor/internal/api/server"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	apiServer, err := apiserver.NewServer()
	if err != nil {
		log.Fatalf("API service: error building application: %v", err)
	}

	errc := make(chan error)
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		err = apiServer.Start()
		if err != nil {
			log.Fatalf("API service: failed to start: %v", err)
			errc <- err
		}
	}()

	select {
	case <-sigc:
		log.Println("received signal to shut down...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		apiServer.Shutdown(ctx)
		cancel()

	case <-errc:
		log.Println("error starting up server, shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		apiServer.Shutdown(ctx)
		cancel()
	}

	os.Exit(0)
}
