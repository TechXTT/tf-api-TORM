#!/bin/bash
set -e

echo "🚀 Starting tf-api-TORM build process..."

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    echo "❌ DATABASE_URL environment variable is required"
    exit 1
fi

# Install TORM CLI if not already installed
if ! command -v torm &> /dev/null; then
    echo "📦 Installing TORM CLI..."
    go install github.com/TechXTT/TORM/cmd/torm@latest
fi

# Wait for database to be ready
echo "⏳ Waiting for database to be ready..."
until pg_isready -d "$DATABASE_URL" -t 1; do
    echo "Database not ready, waiting..."
    sleep 2
done

echo "✅ Database is ready!"

# Create torm directory if it doesn't exist
mkdir -p torm

# Run TORM migrate dev to generate models and run migrations
echo "🔄 Running TORM migrate dev..."
torm migrate dev \
    --schema prisma/schema.prisma \
    --out-migrations torm/migrations \
    --out-models torm/models

echo "✅ TORM models generated successfully!"

# Build the application
echo "🔨 Building application..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

echo "✅ Application built successfully!"

echo "🎉 Build process completed!" 