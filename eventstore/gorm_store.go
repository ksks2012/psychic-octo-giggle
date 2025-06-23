// eventstore/gorm_eventstore.go
package eventstore

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"psychic-octo-giggle/octoevents"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GormEventStore struct {
	db *gorm.DB
}

func NewGormEventStore(dbPath string) (*GormEventStore, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Auto migrate, create event table
	if err := db.AutoMigrate(&Event{}); err != nil {
		return nil, err
	}
	return &GormEventStore{db: db}, nil
}

func (s *GormEventStore) SaveEvents(ctx context.Context, aggregateID string, events []interface{}) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, event := range events {
			// Serialize event to JSON
			eventData, err := json.Marshal(event)
			if err != nil {
				return err
			}
			// Determine event type
			eventType := getEventType(event)
			if eventType == "" {
				return errors.New("unknown event type")
			}
			// Create event record
			dbEvent := Event{
				AggregateID: aggregateID,
				EventType:   eventType,
				EventData:   eventData,
				Timestamp:   getEventTimestamp(event),
			}
			if err := tx.Create(&dbEvent).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *GormEventStore) LoadEvents(ctx context.Context, aggregateID string) ([]interface{}, error) {
	var dbEvents []Event
	// Load events ordered by time
	if err := s.db.WithContext(ctx).
		Where("aggregate_id = ?", aggregateID).
		Order("timestamp ASC").
		Find(&dbEvents).Error; err != nil {
		return nil, err
	}

	// Deserialize events
	var result []interface{}
	for _, dbEvent := range dbEvents {
		event, err := deserializeEvent(dbEvent.EventType, dbEvent.EventData)
		if err != nil {
			return nil, err
		}
		result = append(result, event)
	}
	return result, nil
}

// Helper function: get event type
func getEventType(event interface{}) string {
	switch event.(type) {
	case octoevents.ItemAdded:
		return "ItemAdded"
	case octoevents.ItemRemoved:
		return "ItemRemoved"
	default:
		return ""
	}
}

// Helper function: get event timestamp
func getEventTimestamp(event interface{}) time.Time {
	switch e := event.(type) {
	case octoevents.ItemAdded:
		return e.Timestamp
	case octoevents.ItemRemoved:
		return e.Timestamp
	default:
		return time.Now()
	}
}

// Helper function: deserialize event
func deserializeEvent(eventType string, data []byte) (interface{}, error) {
	switch eventType {
	case "ItemAdded":
		var event octoevents.ItemAdded
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		return event, nil
	case "ItemRemoved":
		var event octoevents.ItemRemoved
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		return event, nil
	default:
		return nil, errors.New("unknown event type")
	}
}
