package domain

import (
	"errors"
	"time"

	"psychic-octo-giggle/octoevents"
)

type ShoppingCart struct {
	ID     string
	Items  map[string]int // itemID -> quantity
	Events []interface{}  // Stores events to be persisted
}

func NewShoppingCart(id string) *ShoppingCart {
	return &ShoppingCart{
		ID:    id,
		Items: make(map[string]int),
	}
}

func (c *ShoppingCart) AddItem(itemID string, quantity int) error {
	// Validation logic
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	// Generate event
	c.Events = append(c.Events, octoevents.ItemAdded{
		CartID:    c.ID,
		ItemID:    itemID,
		Quantity:  quantity,
		Timestamp: time.Now(),
	})
	return nil
}

func (c *ShoppingCart) RebuildFromEvents(events []interface{}) {
	for _, event := range events {
		switch e := event.(type) {
		case octoevents.ItemAdded:
			c.Items[e.ItemID] += e.Quantity
		case octoevents.ItemRemoved:
			c.Items[e.ItemID] -= e.Quantity
		}
	}
}

func (c *ShoppingCart) ClearEvents() {
	c.Events = nil
}
