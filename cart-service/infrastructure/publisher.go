package infrastructure

import (
	"context"
	"encoding/json"

	"psychic-octo-giggle/octoevents"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
)

type EventPublisher struct {
	publisher message.Publisher
}

func NewEventPublisher(natsURL string) (*EventPublisher, error) {
	// Initialize NATS publisher
	publisher, err := nats.NewPublisher(
		nats.PublisherConfig{
			URL:       natsURL,
			Marshaler: &nats.GobMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		return nil, err
	}
	return &EventPublisher{publisher: publisher}, nil
}

func (p *EventPublisher) PublishEvent(ctx context.Context, event interface{}) error {
	// Serialize event to JSON
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Create Watermill message
	msg := message.NewMessage(watermill.NewUUID(), payload)

	// Set event type as message metadata
	eventType := getEventType(event)
	msg.Metadata.Set("event_type", eventType)

	// Publish to NATS topic
	topic := "cart.events"
	return p.publisher.Publish(topic, msg)
}

// Helper function: get event type
func getEventType(event interface{}) string {
	switch event.(type) {
	case octoevents.ItemAdded:
		return "ItemAdded"
	case octoevents.ItemRemoved:
		return "ItemRemoved"
	default:
		return "Unknown"
	}
}
