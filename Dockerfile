# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Download dependencies (if any)
RUN go mod download || true

# Copy source code
COPY main.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o load-tester .

# Final stage - use scratch for minimal image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/load-tester .

# Set the entrypoint
ENTRYPOINT ["./load-tester"]
