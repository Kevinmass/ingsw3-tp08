# Build stage
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates (needed for go modules), create appuser
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates && adduser -D -g '' appuser

# Set working directory
WORKDIR /build

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Build settings for PostgreSQL
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Copy source code - use .dockerignore to filter out unnecessary files
COPY . .

# Build the application
RUN go build -ldflags="-w -s" -o main ./cmd/api/main.go

# Final stage
FROM alpine:3.20

# Install ca-certificates for SSL/TLS and create appuser in one layer
RUN apk --no-cache add ca-certificates && adduser -D -g '' appuser

# Copy the binary from builder
COPY --from=builder /build/main /app/main

# Change ownership to appuser
RUN chown appuser:appuser /app/main

# Switch to appuser
USER appuser

# Set working directory
WORKDIR /app

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./main"]
