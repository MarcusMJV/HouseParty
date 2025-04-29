# 1. Build the Go app
FROM golang:1.21 AS builder

WORKDIR /app

# Copy go.mod and go.sum first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Now copy the full source code
COPY . .

# Build the Go binary
RUN go build -o app

# 2. Create a minimal image to run the app
FROM alpine:latest

WORKDIR /root/

# Copy the built Go binary from builder stage
COPY --from=builder /app/app .

# Expose port 8080 (you can change if needed)
EXPOSE 8080

# Start the Go app
CMD ["./app"]
