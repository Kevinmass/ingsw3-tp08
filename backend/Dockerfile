FROM alpine:3.20

# Install ca-certificates for SSL/TLS and create appuser in one layer
RUN apk --no-cache add ca-certificates && adduser -D -g '' appuser

# Copy the pre-built binary
COPY main /app/main

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
