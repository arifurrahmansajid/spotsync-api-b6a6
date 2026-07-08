# Build stage
FROM golang:latest AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary (static, no CGO)
RUN CGO_ENABLED=0 GOOS=linux go build -a -o spotsync .

# Run stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/spotsync .

# Expose port
EXPOSE 8080

# Run
CMD ["./spotsync"]
