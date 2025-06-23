#!/bin/bash
set -e

echo "🚀 Starting tf-api-TORM deployment..."

# Configuration
IMAGE_NAME="tf-api-torm"
IMAGE_TAG="latest"
NAMESPACE="tf-api-torm"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    print_error "kubectl is not installed. Please install kubectl and try again."
    exit 1
fi

# Check if we're connected to a Kubernetes cluster
if ! kubectl cluster-info &> /dev/null; then
    print_error "Not connected to a Kubernetes cluster. Please configure kubectl and try again."
    exit 1
fi

print_status "Applying PostgreSQL manifests..."
kubectl apply -f k8s/postgres.yaml

print_status "Waiting for PostgreSQL to be ready..."
kubectl wait --for=condition=ready pod -l app=tf-api-torm,component=database -n ${NAMESPACE} --timeout=300s

print_status "Setting up port forward for database access during build..."
# Start port forward in background
kubectl port-forward svc/postgres 5432:5432 -n ${NAMESPACE} &
PF_PID=$!

# Ensure port forward is killed on exit
trap "kill $PF_PID 2>/dev/null || true" EXIT

# Wait a moment for port forward to establish
sleep 5

print_status "Building Docker image with database access..."
# Build the Docker image with DATABASE_URL as build argument
docker build \
    --build-arg DATABASE_URL="postgresql://tfapi_user:tfapi_password@host.docker.internal:5432/tfapi?sslmode=disable" \
    -t ${IMAGE_NAME}:${IMAGE_TAG} .

# Stop port forward
kill $PF_PID 2>/dev/null || true
trap - EXIT # Clear the trap

print_status "Deploying tf-api-TORM application..."
kubectl apply -f k8s/tf-api-torm.yaml

print_status "Waiting for tf-api-TORM to be ready..."
kubectl wait --for=condition=ready pod -l app=tf-api-torm,component=api -n ${NAMESPACE} --timeout=300s

print_status "Getting service information..."
kubectl get svc -n ${NAMESPACE}

print_status "Getting pod status..."
kubectl get pods -n ${NAMESPACE}

print_status "🎉 Deployment completed successfully!"

echo ""
print_status "The service is being exposed via a LoadBalancer."
print_status "It might take a minute for the external IP to be available."
print_status "Check the status with: kubectl get svc tf-api-torm -n ${NAMESPACE}"
echo ""
print_status "Once the EXTERNAL-IP is assigned (it will likely be 'localhost' on Docker Desktop),"
print_status "you can access the application at port 80:"
echo "  http://<EXTERNAL-IP>/v1/   (e.g., http://localhost/v1/)"
echo ""
print_status "To view logs:"
echo "  - Application: kubectl logs -f deployment/tf-api-torm -n ${NAMESPACE}"
echo "  - Database: kubectl logs -f deployment/postgres -n ${NAMESPACE}"
echo ""
print_status "To delete the deployment:"
echo "  - kubectl delete namespace ${NAMESPACE}" 