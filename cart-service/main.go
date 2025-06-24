package main

import (
	"context"
	"log"
	"os"
	"psychic-octo-giggle/cart-service/application"
	"psychic-octo-giggle/cart-service/infrastructure"
	"psychic-octo-giggle/eventstore"
)

func main() {
	// Initialize event store
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	eventStore, err := eventstore.NewGormEventStore("./events.db")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize event publisher
	publisher, err := infrastructure.NewEventPublisher(natsURL)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize service
	cartService := application.NewCartService(eventStore, publisher)

	// Simulate adding an item
	ctx := context.Background()
	if err := cartService.AddItem(ctx, "cart-123", "item-001", 2); err != nil {
		log.Fatal(err)
	}
	log.Println("Item added successfully")

	select {}
}
