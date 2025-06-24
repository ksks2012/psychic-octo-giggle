# Microservice Connector

Microservice Connector is a simple side project to learn **Event Sourcing** and event-driven architecture using Go. It demonstrates a shopping cart service that generates events (e.g., `ItemAdded`, `ItemRemoved`) and stores them in a SQLite-based event store using GORM. These events are published to other services (e.g., billing service) via NATS using the Watermill library.

## Project Overview

The project implements a basic Event Sourcing scenario where:
- The **Cart Service** handles commands (e.g., adding/removing items from a cart) and generates events.
- Events are stored in a SQLite event store using GORM.
- Events are published to a NATS topic (`cart.events`) using Watermill.
- The **Billing Service** subscribes to these events and processes them (e.g., simulating billing logic).

This project is designed for learning purposes, focusing on:
- Event Sourcing concepts (event storage, state rebuilding).
- Event-driven communication between microservices.
- Integration of Go libraries like GORM and Watermill.

## Tech Stack

- **Go**: Programming language for implementing services.
- **GORM**: ORM for SQLite-based event storage.
- **Watermill**: Event-driven messaging library for publishing and subscribing to events.
- **NATS**: Lightweight messaging system for event transmission.
- **SQLite**: Lightweight database for event storage.
- **encoding/json**: Used for event serialization.

## Project Structure

```
/microservice-connector
├── /cart-service
│   ├── /domain         # Defines aggregates (e.g., ShoppingCart) and events (e.g., ItemAdded, ItemRemoved)
│   ├── /application    # Command handling logic (e.g., AddItemCommand)
│   ├── /infrastructure # Event storage (SQLite) and message publishing (Watermill/NATS)
│   └── main.go         # Service entry point
├── /billing-service
│   ├── /domain         # Logic for consuming events
│   ├── /application    # Event handling logic
│   ├── /infrastructure # Message subscription (Watermill/NATS)
│   └── main.go         # Service entry point
├── /eventstore         # Shared SQLite event store implementation
├── /events             # Shared Protobuf/JSON event definitions
└── docker-compose.yml  # Deploys NATS and services
```

## Setup Instructions

### Prerequisites

- **Go**: Version 1.23 or higher.
- **Docker**: For running NATS and services.
- **SQLite**: No separate installation needed (embedded in Go via GORM).

### Installation

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/ksks2012/psychic-octo-giggle.git
   cd microservice-connector
   ```

2. **Install Dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run NATS and Services with Docker Compose**:
   ```bash
   docker-compose up --build
   ```

   This starts:
   - NATS server (`nats://localhost:4222`).
   - Cart Service (publishes events).
   - Billing Service (subscribes to events).

4. **Database Setup**:
   - The cart service automatically creates a SQLite database (`events.db`) for event storage.
   - No manual database setup is required.

### Running Locally Without Docker

1. **Start NATS Server**:
   ```bash
   docker run -d -p 4222:4222 nats
   ```

2. **Run Cart Service**:
   ```bash
   cd cart-service
   go run main.go
   ```

3. **Run Billing Service**:
   ```bash
   cd billing-service
   go run main.go
   ```

## Usage

1. **Add an Item to the Cart**:
   - The `cart-service` handles commands like `AddItem`, generating an `ItemAdded` event.
   - The event is stored in the SQLite event store (`events.db`) and published to the `cart.events` topic in NATS.
   - Example command (triggered in code):
     ```go
     cartService.AddItem(ctx, "cart-123", "item-001", 2)
     ```

2. **Event Consumption**:
   - The `billing-service` subscribes to the `cart.events` topic and processes events (e.g., logs `ItemAdded` events).
   - Check the `billing-service` logs to verify event consumption.

3. **Verify Events in SQLite**:
   - Open the SQLite database to inspect stored events:
     ```bash
     sqlite3 events.db
     SELECT * FROM events;
     ```

## Event Store Schema

The event store uses a single `events` table in SQLite with the following schema:

```sql
CREATE TABLE events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    aggregate_id TEXT,
    event_type TEXT NOT NULL,
    event_data BLOB NOT NULL,
    timestamp DATETIME,
    INDEX idx_events_aggregate_id (aggregate_id),
    INDEX idx_events_timestamp (timestamp)
);
```

- `aggregate_id`: The cart ID (e.g., `cart-123`).
- `event_type`: The type of event (e.g., `ItemAdded`, `ItemRemoved`).
- `event_data`: JSON-serialized event data.
- `timestamp`: Event creation time.

## Event Types

Events are defined in the `events` package:

- **ItemAdded**:
  ```go
  type ItemAdded struct {
      CartID    string    `json:"cart_id"`
      ItemID    string    `json:"item_id"`
      Quantity  int       `json:"quantity"`
      Timestamp time.Time `json:"timestamp"`
  }
  ```

- **ItemRemoved**:
  ```go
  type ItemRemoved struct {
      CartID    string    `json:"cart_id"`
      ItemID    string    `json:"item_id"`
      Quantity  int       `json:"quantity"`
      Timestamp time.Time `json:"timestamp"`
  }
  ```

## TODO: Future Improvements

- **Snapshots**: Implement snapshotting to optimize state rebuilding for large event streams.
- **Protobuf**: Replace JSON with Protobuf for more efficient event serialization.
- **Kafka Support**: Add support for Kafka as an alternative to NATS using Watermill.
- **Outbox Pattern**: Ensure consistency between event storage and publishing using the transactional outbox pattern.
- **Monitoring**: Integrate Prometheus for tracking event publishing and consumption metrics.

## Learning Resources

- [Watermill Documentation](https://watermill.io/)
- [NATS Documentation](https://docs.nats.io/)
- [GORM Documentation](https://gorm.io/)
- [Event Sourcing by Martin Fowler](https://martinfowler.com/eaaDev/EventSourcing.html)