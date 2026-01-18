# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrations ./cmd/migrations/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binaries from builder
COPY --from=builder /app/main .
COPY --from=builder /app/migrations .

# Copy migration files
COPY --from=builder /app/internal/repository/db/migrations ./internal/repository/db/migrations

# Copy misc files (for fonts used in the application)
COPY --from=builder /app/misc ./misc

# Copy .env file
COPY --from=builder /app/.env .

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./main"]
