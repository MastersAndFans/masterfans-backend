# Build
FROM golang:1.22.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/masterfansapp ./cmd/masterfansapp

# Production
FROM alpine:latest AS production

# Install bash and CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates bash

WORKDIR /app

# Copy binary from builder to production stage
COPY --from=builder /bin/masterfansapp /app/masterfansapp

# Add wait-for-it script
COPY wait-for-it.sh /app/wait-for-it.sh
RUN chmod +x /app/wait-for-it.sh

EXPOSE 5000

# Use wait-for-it.sh to wait for the db, then start the application
CMD ["/app/wait-for-it.sh", "db:5432", "--timeout=30", "--", "/app/masterfansapp"]

