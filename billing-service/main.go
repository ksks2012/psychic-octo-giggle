package main

import (
	"context"
	"log"

	"psychic-octo-giggle/billing-service/application"
	"psychic-octo-giggle/billing-service/infrastructure"
)

func main() {
	billingService := application.NewBillingService()

	subscriber, err := infrastructure.NewEventSubscriber("nats://localhost:4222", billingService)
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
