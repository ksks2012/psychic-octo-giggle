package main

import (
	"context"
	"log"
	"os"

	"psychic-octo-giggle/billing-service/application"
	"psychic-octo-giggle/billing-service/infrastructure"
)

func main() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	billingService := application.NewBillingService()

	subscriber, err := infrastructure.NewEventSubscriber(natsURL, billingService)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if err := subscriber.Subscribe(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Billing service started, subscribing to events...")

	select {}
}
