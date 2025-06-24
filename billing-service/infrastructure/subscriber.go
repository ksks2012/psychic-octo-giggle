package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"psychic-octo-giggle/octoevents"

	"github.com/nats-io/nats.go"
)

type EventHandler interface {
	HandleItemAdded(ctx context.Context, event octoevents.ItemAdded) error
	HandleItemRemoved(ctx context.Context, event octoevents.ItemRemoved) error
}

type EventSubscriber struct {
	nc      *nats.Conn
	js      nats.JetStreamContext
	handler EventHandler
	sub     *nats.Subscription
}

func NewEventSubscriber(natsURL string, handler EventHandler) (*EventSubscriber, error) {
	// Connect to NATS server
	nc, err := nats.Connect(natsURL, nats.Name("billing-service-subscriber"))
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
			MaxMsgs:  -1,
			MaxAge:   0,
		})
		if err != nil {
			nc.Close()
			return nil, fmt.Errorf("failed to create stream %s: %w", streamName, err)
		}
		log.Printf("Created stream %s", streamName)
	} else {
		log.Printf("Stream %s already exists", streamName)
	}

	return &EventSubscriber{
		nc:      nc,
		js:      js,
		handler: handler,
	}, nil
}

func (s *EventSubscriber) Close() error {
	if s.sub != nil {
		if err := s.sub.Unsubscribe(); err != nil {
			log.Printf("Failed to unsubscribe: %v", err)
		}
	}
	if s.nc != nil {
		s.nc.Close()
	}
	return nil
}

func (s *EventSubscriber) Subscribe(ctx context.Context) error {
	topic := "cart.events"
	var subErr error
	s.sub, subErr = s.js.QueueSubscribe(topic, "billing-service-queue", func(msg *nats.Msg) {
		if err := s.processMessage(ctx, msg); err != nil {
			log.Printf("Failed to process message: %v", err)
			if err := msg.Nak(); err != nil {
				log.Printf("Failed to NAK message: %v", err)
			}
			return
		}
		if err := msg.Ack(); err != nil {
			log.Printf("Failed to ACK message: %v", err)
		}
	}, nats.Durable("billing-service"), nats.ManualAck(), nats.Context(ctx))
	if subErr != nil {
		return fmt.Errorf("failed to subscribe to %s: %w", topic, subErr)
	}
	log.Printf("Subscribed to %s with queue group billing-service-queue", topic)
	return nil
}

func (s *EventSubscriber) processMessage(ctx context.Context, msg *nats.Msg) error {
	log.Printf("Received message: Subject=%s, Payload=%s, Headers=%v", msg.Subject, string(msg.Data), msg.Header)

	eventType := msg.Header.Get("event_type")
	switch eventType {
	case "ItemAdded":
		var event octoevents.ItemAdded
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			return fmt.Errorf("failed to unmarshal ItemAdded: %w", err)
		}
		return s.handler.HandleItemAdded(ctx, event)
	case "ItemRemoved":
		var event octoevents.ItemRemoved
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			return fmt.Errorf("failed to unmarshal ItemRemoved: %w", err)
		}
		return s.handler.HandleItemRemoved(ctx, event)
	default:
		log.Printf("Unknown event type: %s", eventType)
		return nil
	}
}
