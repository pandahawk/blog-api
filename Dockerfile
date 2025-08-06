# 1. Build stage: use the official Go image to build the binary
FROM golang:alpine AS builder

WORKDIR /app

# Only download go modules first (caching layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary
RUN go build -o blog-api .

# 2. Minimal final image
FROM alpine:3.19

# Set up an unprivileged user (optional but good for security)
RUN adduser -D appuser

# Copy the binary from the builder stage
COPY --from=builder /app/blog-api /usr/local/bin/blog-api

# Copy migrations folder from builder
COPY --from=builder /app/internal/database/migrations /home/appuser/internal/database/migrations


USER appuser

# Set working directory
WORKDIR /home/appuser

# Expose the port your app listens on (default 8080 for Gin or net/http)
EXPOSE 8080

# Start the app (read PORT env for cloud platforms; defaults to 8080)
ENV PORT=8080

CMD ["blog-api"]
