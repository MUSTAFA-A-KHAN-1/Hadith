# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o hadith-bot ./cmd/bot

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/hadith-bot .

# Copy data directory (optional - for offline data)
COPY --from=builder /app/data ./data

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

# Expose port (if using webhooks)
EXPOSE 8080

# Set environment variables
ENV TELEGRAM_BOT_TOKEN=""
ENV API_URL="https://api.sunnah.com/v1"
ENV API_KEY=""
ENV LOG_LEVEL="info"
ENV RATE_LIMIT_REQUESTS=10
ENV RATE_LIMIT_WINDOW="1m"

# Run the bot
CMD ["./hadith-bot"]

