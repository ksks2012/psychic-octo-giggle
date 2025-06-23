package main

import (
	"context"
	"log"
	"psychic-octo-giggle/cart-service/application"
	"psychic-octo-giggle/cart-service/infrastructure"
	"psychic-octo-giggle/eventstore"
)

func main() {
	// Initialize event store
	eventStore, err := eventstore.NewGormEventStore("events.db")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize event publisher
	publisher, err := infrastructure.NewEventPublisher("nats://localhost:4222")
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
}
