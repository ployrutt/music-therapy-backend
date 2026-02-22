# ---- Build Stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (required for some Go modules)
RUN apk add --no-cache git

# Copy go.mod and go.sum first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .

# ---- Runtime Stage ----
FROM alpine:3.21

WORKDIR /app

# Install ca-certificates for HTTPS and tzdata for timezone
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/server .

# Create uploads directory
RUN mkdir -p /app/uploads

# Expose port
EXPOSE 8080

# Environment variables (override via Dokploy)
# ENV DB_HOST=localhost \
#     DB_PORT=5432 \
#     DB_USER=postgres \
#     DB_PASSWORD= \
#     DB_NAME=project-backend \
#     GIN_MODE=release

ENTRYPOINT ["./server"]
