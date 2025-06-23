# tf-api-TORM

A Go API application using TORM (TypeORM-inspired ORM) with PostgreSQL database.

## Prerequisites

- Docker
- Kubernetes cluster (local or remote)
- kubectl configured to access your cluster
- Go 1.23+ (for local development)

## Quick Start

### 1. Clone and Navigate
```bash
cd tf-api-TORM
```

### 2. Deploy to Kubernetes
```bash
# Make deployment script executable
chmod +x deploy.sh

# Run deployment
./deploy.sh
```

### 3. Access the Application
The service is exposed via a LoadBalancer. After running `./deploy.sh`, it may take a minute for the external IP to be available.

You can check the status and find the `EXTERNAL-IP` by running:
```bash
kubectl get svc tf-api-torm -n tf-api-torm
```

Once the `EXTERNAL-IP` is assigned (for Docker Desktop, this will typically be `localhost`), you can access the API on port 80. For example:
```bash
# Visit the API (replace localhost if your EXTERNAL-IP is different)
curl http://localhost/v1/
```

## Architecture

The deployment includes:

- **PostgreSQL Database**: Persistent storage with ConfigMap and Secret management
- **tf-api-TORM Application**: Go API with TORM ORM integration
- **Kubernetes Resources**: Namespace, Deployments, Services, and Ingress

### Components

All Kubernetes resources are defined in a single file (`k8s/deployment.yaml`):

- **Namespace**: `tf-api-torm` namespace for isolation
- **ConfigMaps**: Database and application configuration
- **Secrets**: Database password and sensitive data
- **PostgreSQL**: Database deployment with persistent storage
- **tf-api-TORM API**: Application deployment with init containers
- **Services**: Internal communication between components
- **Ingress**: External access (commented out by default)

## Configuration

### Environment Variables

The application requires the following environment variables:

- `DATABASE_URL`: PostgreSQL connection string
- `PRIVATE_KEY`: JWT private key
- `PUBLIC_KEY`: JWT public key
- `GMAIL_CLIENT_ID`: Gmail OAuth client ID
- `GMAIL_CLIENT_SECRET`: Gmail OAuth client secret
- `GMAIL_REFRESH_TOKEN`: Gmail OAuth refresh token

### Database Configuration

Default database settings:
- Database: `tfapi`
- User: `tfapi_user`
- Password: `tfapi_password`

## Development

### Local Development

1. **Start PostgreSQL**:
```bash
docker run -d \
  --name postgres-tfapi \
  -e POSTGRES_DB=tfapi \
  -e POSTGRES_USER=tfapi_user \
  -e POSTGRES_PASSWORD=tfapi_password \
  -p 5432:5432 \
  postgres:15-alpine
```

2. **Set environment variables**:
```bash
export DATABASE_URL="postgresql://tfapi_user:tfapi_password@localhost:5432/tfapi?sslmode=disable"
```

3. **Generate TORM models and run migrations**:
```bash
chmod +x build.sh
./build.sh
```

4. **Run the application**:
```bash
go run main.go
```

### Building Docker Image

```bash
docker build -t tf-api-torm:latest .
```

### Manual Kubernetes Deployment

If you prefer to deploy manually:

```bash
# Build the Docker image
docker build -t tf-api-torm:latest .

# Apply the Kubernetes configuration
kubectl apply -f k8s/deployment.yaml

# Check the deployment status
kubectl get pods -n tf-api-torm
```

## API Endpoints

The API provides the following endpoints:

- `GET /v1/get/projects` - Get all projects
- `GET /v1/get/project/{id}` - Get specific project
- `GET /v1/get/projects/{category}` - Get projects by category
- `POST /v1/post/vote` - Submit a vote
- `PUT /v1/update/verify_vote` - Verify a vote

## Monitoring and Logs

### View Logs
```bash
# Application logs
kubectl logs -f deployment/tf-api-torm -n tf-api-torm

# Database logs
kubectl logs -f deployment/postgres -n tf-api-torm
```

### Check Pod Status
```bash
kubectl get pods -n tf-api-torm
```

### Check Services
```bash
kubectl get svc -n tf-api-torm
```

## Troubleshooting

### Common Issues

1. **Database Connection Issues**
   - Check if PostgreSQL pod is running: `kubectl get pods -n tf-api-torm`
   - Check database logs: `kubectl logs deployment/postgres -n tf-api-torm`

2. **TORM Generation Issues**
   - Ensure DATABASE_URL is correctly set
   - Check if TORM CLI is installed in the container
   - Verify database is accessible from the application pod

3. **Application Startup Issues**
   - Check application logs: `kubectl logs deployment/tf-api-torm -n tf-api-torm`
   - Verify all environment variables are set
   - Check if the application can connect to the database

### Cleanup

To remove the entire deployment:
```bash
kubectl delete namespace tf-api-torm
```

## Security Notes

- Update the default database password in production
- Store sensitive environment variables in Kubernetes Secrets
- Configure proper ingress rules for production use
- Enable SSL/TLS for database connections in production

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test the deployment
5. Submit a pull request