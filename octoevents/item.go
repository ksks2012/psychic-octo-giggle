package octoevents

import "time"

type ItemAdded struct {
	CartID    string    `json:"cart_id"`
	ItemID    string    `json:"item_id"`
	Quantity  int       `json:"quantity"`
	Timestamp time.Time `json:"timestamp"`
}

type ItemRemoved struct {
	CartID    string    `json:"cart_id"`
	ItemID    string    `json:"item_id"`
	Quantity  int       `json:"quantity"`
	Timestamp time.Time `json:"timestamp"`
}
