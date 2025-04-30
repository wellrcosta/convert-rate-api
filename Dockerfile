# Stage 1: Build
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the full source
COPY . .

# Build the application
RUN go build -o converter ./cmd/main.go

# Stage 2: Final
FROM alpine:latest
WORKDIR /root/

# Add CA certificates (required for HTTPS calls)
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/converter .

# Expose port and run
EXPOSE 3000
CMD ["./converter"]
