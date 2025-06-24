package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"psychic-octo-giggle/octoevents"

	"github.com/nats-io/nats.go"
)

type EventPublisher struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

func NewEventPublisher(natsURL string) (*EventPublisher, error) {
	// Connect to NATS server
	nc, err := nats.Connect(natsURL, nats.Name("cart-service-publisher"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Initialize JetStream context
	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to initialize JetStream: %w", err)
	}

	// Check or create the stream
	streamName := "cart"
	_, err = js.StreamInfo(streamName)
	if err != nil {
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{"cart.events"},
			Storage:  nats.FileStorage,
			MaxMsgs:  -1, // No limit on number of messages
			MaxAge:   0,  // No expiration
		})
		if err != nil {
			nc.Close()
			return nil, fmt.Errorf("failed to create stream %s: %w", streamName, err)
		}
		log.Printf("Created stream %s", streamName)
	} else {
		log.Printf("Stream %s already exists", streamName)
	}

	return &EventPublisher{
		nc: nc,
		js: js,
	}, nil
}

func (p *EventPublisher) Close() error {
	if p.nc != nil {
		p.nc.Close()
	}
	return nil
}

func (p *EventPublisher) PublishEvent(ctx context.Context, event interface{}) error {
	if event == nil {
		return fmt.Errorf("event is nil")
	}

	// Serialize event to JSON
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	if len(payload) == 0 || string(payload) == "{}" {
		return fmt.Errorf("empty or invalid payload for event: %T", event)
	}

	// Set event type as NATS message header
	eventType := getEventType(event)
	headers := nats.Header{}
	headers.Set("event_type", eventType)

	// Log the event being published
	log.Printf("Publishing event: type=%s, payload=%s", eventType, string(payload))

	// Publish to NATS JetStream topic
	topic := "cart.events"
	_, err = p.js.PublishMsg(&nats.Msg{
		Subject: topic,
		Data:    payload,
		Header:  headers,
	}, nats.Context(ctx))
	if err != nil {
		return fmt.Errorf("failed to publish event to %s: %w", topic, err)
	}

	return nil
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
