# Multi-stage build for tf-api-TORM
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the CLI from the same release as the generated client's runtime.
RUN go install "github.com/TechXTT/TORM/v2/cmd/torm@$(go list -m -f '{{.Version}}' github.com/TechXTT/TORM/v2)" && \
    torm generate --out-models torm/models && \
    CGO_ENABLED=0 go build -o main .

# Stage 2: Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates postgresql-client bash

# Create app user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy TORM binary from builder stage
COPY --from=builder /go/bin/torm /usr/local/bin/

# Copy the built application
COPY --from=builder /app/main /app/main

# Copy prisma schema
COPY --from=builder /app/prisma /app/prisma

# Copy the pkg directory which contains email templates and other assets
COPY --from=builder /app/pkg /app/pkg

# Copy generated TORM models and migrations
COPY --from=builder /app/torm /app/torm

# Create torm directory if it doesn't exist
RUN mkdir -p torm && chown -R appuser:appgroup /app

# Switch to app user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/v1/ || exit 1

# Run the application
CMD ["./main"] 