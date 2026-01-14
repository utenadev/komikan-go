package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kench/komikan-go/internal/nostr"
)

func main() {
	log.Println("Starting Komikan Bot...")

	// TODO: Load config
	// TODO: Initialize database
	// TODO: Initialize Nostr client

	// Start Nostr client
	client, err := nostr.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Nostr client: %v", err)
	}

	// Connect to relays
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect to relays: %v", err)
	}

	log.Println("Connected to Nostr relays")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	client.Disconnect()
	fmt.Println("Bye!")
}
