# Build
FROM golang:1.22.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/myapp ./cmd/myapp

# Run
FROM alpine:latest

# CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder to run stage
COPY --from=builder /bin/myapp /app/

EXPOSE 5000

CMD ["./myapp"]

