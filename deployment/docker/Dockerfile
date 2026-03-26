# Build stage
FROM golang:1.26.1-alpine AS builder

# Install build dependencies (needed for CGO if using certain DB drivers)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy dependency files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the binary
RUN go build -o sense cmd/main.go

# Final stage
FROM alpine:latest

# Install CA certificates for secure connections (HTTPS/TLS)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/sense .

# Copy necessary runtime assets
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/config ./config

EXPOSE 8080

# Run the application
CMD ["./sense"]