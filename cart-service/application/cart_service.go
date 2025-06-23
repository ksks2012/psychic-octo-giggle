package application

import (
	"context"

	"psychic-octo-giggle/cart-service/domain"
	"psychic-octo-giggle/cart-service/infrastructure"
	"psychic-octo-giggle/eventstore"
)

type CartService struct {
	eventStore eventstore.EventStore
	publisher  *infrastructure.EventPublisher
}

func NewCartService(eventStore eventstore.EventStore, publisher *infrastructure.EventPublisher) *CartService {
	return &CartService{
		eventStore: eventStore,
		publisher:  publisher,
	}
}

func (s *CartService) AddItem(ctx context.Context, cartID, itemID string, quantity int) error {
	// Load the shopping cart (by replaying events)
	cart := domain.NewShoppingCart(cartID)
	events, err := s.eventStore.LoadEvents(ctx, cartID)
	if err != nil {
		return err
	}
	cart.RebuildFromEvents(events)

	// Execute command
	if err := cart.AddItem(itemID, quantity); err != nil {
		return err
	}

	// Save events to EventStore
	if err := s.eventStore.SaveEvents(ctx, cartID, cart.Events); err != nil {
		return err
	}

	// Publish events to NATS
	for _, event := range cart.Events {
		if err := s.publisher.PublishEvent(ctx, event); err != nil {
			return err
		}
	}

	// Clear processed events
	cart.ClearEvents()
	return nil
}
