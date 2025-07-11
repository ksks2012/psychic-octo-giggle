# Use the official Go image as the base image
FROM golang:1.23 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY billing-service/ billing-service/
COPY octoevents/ octoevents/
COPY eventstore/ eventstore/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /billing-service ./billing-service/main.go

# Use a minimal base image for the final stage
FROM alpine:latest

# Install ca-certificates for HTTPS calls (if needed)
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /billing-service /billing-service

# Set the binary as the entrypoint
ENTRYPOINT ["/billing-service"]