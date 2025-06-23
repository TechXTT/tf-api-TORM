# Multi-stage build for tf-api-TORM
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git postgresql-client bash

# Install TORM CLI
RUN go install github.com/TechXTT/TORM/cmd/torm@latest

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Make build script executable
RUN chmod +x build.sh

# Copy TORM binary to a location in PATH
RUN cp $(go env GOPATH)/bin/torm /usr/local/bin/

# Run TORM migrate dev to generate models (DATABASE_URL will be passed as build arg)
ARG DATABASE_URL
ENV DATABASE_URL=$DATABASE_URL

# Wait for database and run TORM migrate dev
RUN if [ -n "$DATABASE_URL" ]; then \
        echo "Waiting for database..." && \
        until pg_isready -d "$DATABASE_URL" -t 1; do \
            echo "Database not ready, waiting..." && \
            sleep 2; \
        done && \
        echo "Database ready! Running TORM migrate dev..." && \
        torm migrate dev \
            --schema prisma/schema.prisma \
            --out-migrations torm/migrations \
            --out-models torm/models; \
    else \
        echo "No DATABASE_URL provided, skipping TORM generation"; \
    fi

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

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
COPY --from=builder /usr/local/bin/torm /usr/local/bin/

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