package main

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func testNATSConnection(natsURL string) error {
	nc, err := nats.Connect(natsURL, nats.Name("test-connection"))
	if err != nil {
		return err
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		return err
	}

	streamName := "cart"
	_, err = js.StreamInfo(streamName)
	if err != nil {
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{"cart.events"},
			Storage:  nats.FileStorage,
		})
		if err != nil {
			return err
		}
		log.Printf("Created stream %s", streamName)
	} else {
		log.Printf("Stream %s already exists", streamName)
	}

	_, err = js.Publish("cart.events", []byte("Test message"))
	if err != nil {
		return err
	}
	log.Println("Successfully published test message to cart.events")

	return nil
}

func main() {
	natsURL := "nats://localhost:4222"
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := testNATSConnection(natsURL); err != nil {
		log.Fatalf("NATS test failed: %v", err)
	}
	log.Println("NATS test successful")
}
