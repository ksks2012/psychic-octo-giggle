# Use the official Go image as the base image
FROM golang:1.23-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY cart-service/ cart-service/
COPY octoevents/ octoevents/
COPY eventstore/ eventstore/

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -v -o /cart-service ./cart-service/main.go && ls -l /cart-service || exit 1

# Use a minimal base image for the final stage
FROM alpine:latest

# Install ca-certificates for HTTPS calls (if needed)
RUN apk --no-cache add ca-certificates sqlite-libs

# Copy the binary from the builder stage
COPY --from=builder /cart-service /cart-service

# Set the binary as the entrypoint
ENTRYPOINT ["/cart-service"]