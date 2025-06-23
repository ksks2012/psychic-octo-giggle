
# Structure
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