# Build stage
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates (needed for go modules)
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

# Create appuser
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /build

# Copy go.mod first
COPY go.mod ./

# Download dependencies (creates go.sum)
RUN go mod download

# Enable CGO for sqlite3
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

# Install gcc for CGO
RUN apk add --no-cache gcc musl-dev

# Copy source code
COPY . .

# Build the application
RUN go build -ldflags="-w -s" -o main cmd/api/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for SSL/TLS
RUN apk --no-cache add ca-certificates

# Create appuser
RUN adduser -D -g '' appuser

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
