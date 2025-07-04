# Multi-stage build for Go backend
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates for fetching dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o qafka ./cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN adduser -D -g '' qafka

# Set working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/qafka .

# Change ownership to non-root user
RUN chown qafka:qafka qafka

# Switch to non-root user
USER qafka

# Expose port
EXPOSE 8080

# Command to run
CMD ["./qafka"]
