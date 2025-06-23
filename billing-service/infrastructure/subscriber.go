package infrastructure

import (
	"context"
	"encoding/json"
	"log"

	"psychic-octo-giggle/octoevents"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
)

type EventHandler interface {
	HandleItemAdded(ctx context.Context, event octoevents.ItemAdded) error
	HandleItemRemoved(ctx context.Context, event octoevents.ItemRemoved) error
}

type EventSubscriber struct {
	subscriber message.Subscriber
	handler    EventHandler
}

func NewEventSubscriber(natsURL string, handler EventHandler) (*EventSubscriber, error) {
	// Initialize NATS subscriber
	subscriber, err := nats.NewSubscriber(
		nats.SubscriberConfig{
			URL:         natsURL,
			Unmarshaler: &nats.GobMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		return nil, err
	}
	return &EventSubscriber{
		subscriber: subscriber,
		handler:    handler,
	}, nil
}

func (s *EventSubscriber) Subscribe(ctx context.Context) error {
	topic := "cart.events"
	messages, err := s.subscriber.Subscribe(ctx, topic)
	if err != nil {
		return err
	}

	// Handle messages
	go func() {
		for msg := range messages {
			if err := s.processMessage(ctx, msg); err != nil {
				log.Printf("Failed to process message: %v", err)
				msg.Nack() // Processing failed, notify Watermill to retry
				continue
			}
			msg.Ack() // Processing succeeded, acknowledge the message
		}
	}()
	return nil
}

func (s *EventSubscriber) processMessage(ctx context.Context, msg *message.Message) error {
	eventType := msg.Metadata.Get("event_type")
	switch eventType {
	case "ItemAdded":
		var event octoevents.ItemAdded
		if err := json.Unmarshal(msg.Payload, &event); err != nil {
			return err
		}
		return s.handler.HandleItemAdded(ctx, event)
	case "ItemRemoved":
		var event octoevents.ItemRemoved
		if err := json.Unmarshal(msg.Payload, &event); err != nil {
			return err
		}
		return s.handler.HandleItemRemoved(ctx, event)
	default:
		log.Printf("Unknown event type: %s", eventType)
		return nil
	}
}
