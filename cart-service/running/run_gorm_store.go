// cart-service/main.go
package main

import (
	"context"
	"log"
	"time"

	"psychic-octo-giggle/events"
	"psychic-octo-giggle/eventstore"
)

func main() {
	// Initialize event store
	store, err := eventstore.NewGormEventStore("./events.db")
	if err != nil {
		log.Fatal(err)
	}

	// Simulate adding an item to the cart
	cartID := "cart-123"
	event := events.ItemAdded{
		CartID:    cartID,
		ItemID:    "item-001",
		Quantity:  2,
		Timestamp: time.Now(),
	}

	// Save event
	ctx := context.Background()
	if err := store.SaveEvents(ctx, cartID, []interface{}{event}); err != nil {
		log.Fatal(err)
	}

	// Load events
	loadedEvents, err := store.LoadEvents(ctx, cartID)
	if err != nil {
		log.Fatal(err)
	}

	// Print events
	for _, e := range loadedEvents {
		switch event := e.(type) {
		case events.ItemAdded:
			log.Printf("ItemAdded: %+v", event)
		case events.ItemRemoved:
			log.Printf("ItemRemoved: %+v", event)
		}
	}
}
