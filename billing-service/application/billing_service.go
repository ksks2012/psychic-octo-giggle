// billing-service/application/billing_service.go
package application

import (
	"context"
	"log"

	"psychic-octo-giggle/octoevents"
)

type BillingService struct{}

func NewBillingService() *BillingService {
	return &BillingService{}
}

func (s *BillingService) HandleItemAdded(ctx context.Context, event octoevents.ItemAdded) error {
	log.Printf("Processing ItemAdded: CartID=%s, ItemID=%s, Quantity=%d",
		event.CartID, event.ItemID, event.Quantity)
	// Simulate billing logic, for example, calculate total amount
	return nil
}

func (s *BillingService) HandleItemRemoved(ctx context.Context, event octoevents.ItemRemoved) error {
	log.Printf("Processing ItemRemoved: CartID=%s, ItemID=%s, Quantity=%d",
		event.CartID, event.ItemID, event.Quantity)
	// Simulate updating billing logic
	return nil
}
