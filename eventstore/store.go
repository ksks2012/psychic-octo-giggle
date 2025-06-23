package eventstore

import (
	"context"
	"time"
)

type Event struct {
	ID          uint      `gorm:"primaryKey"`
	AggregateID string    `gorm:"index"`              // Aggregate root ID (CartID)
	EventType   string    `gorm:"not null"`           // Event type (ItemAdded, ItemRemoved)
	EventData   []byte    `gorm:"type:blob;not null"` // Serialized event data (JSON)
	Timestamp   time.Time `gorm:"index"`              // Event timestamp
}

type EventStore interface {
	SaveEvents(ctx context.Context, aggregateID string, events []interface{}) error
	LoadEvents(ctx context.Context, aggregateID string) ([]interface{}, error)
}
