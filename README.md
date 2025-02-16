# do-database-metrics-adapter

A Kubernetes metrics adapter that collects and exposes metrics from DigitalOcean Managed Databases. This adapter retrieves metrics from the DigitalOcean API and exposes them in a format compatible with Kubernetes metrics API, making it ideal for monitoring DigitalOcean databases in a Kubernetes environment.

## Features

- Retrieves metrics from DigitalOcean Managed Databases
- Exposes metrics via HTTP endpoints
- Supports Kubernetes metrics API format
- Includes health and readiness checks
- Provides version information
- Supports secure metrics authentication

## Prerequisites

- Go 1.21.6 or later
- DigitalOcean API Token
- Database name to monitor
- Docker (for container deployment)
- Kubernetes (for Helm deployment)

## Configuration

The application is configured through environment variables:

- `DEBUG`: Enable debug logging (boolean)
- `PORT`: Server port number
- `DO_TOKEN`: DigitalOcean API token
- `DATABASE_NAME`: Name of the DigitalOcean database to monitor

## Deployment Options

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/supporttools/do-database-metrics-adapter.git
   cd do-database-metrics-adapter
   ```

2. Build and run:
   ```bash
   go mod download
   go build
   ./do-database-metrics-adapter
   ```

### Docker Deployment

1. Build the container:
   ```bash
   docker build -t do-database-metrics-adapter .
   ```

2. Run the container:
   ```bash
   docker run -e DO_TOKEN=<your-token> \
             -e DATABASE_NAME=<your-database> \
             -e PORT=8080 \
             -p 8080:8080 \
             do-database-metrics-adapter
   ```

### Helm Deployment

The project includes a Helm chart for Kubernetes deployment:

1. Add the repository:
   ```bash
   helm repo add supporttools https://charts.support.tools
   helm repo update
   ```

2. Install the chart:
   ```bash
   helm install do-database-metrics-adapter supporttools/do-database-metrics-adapter \
     --set doToken=<your-token> \
     --set databaseName=<your-database>
   ```

## API Endpoints

- `/metrics`: Exposes database metrics in Prometheus format
- `/healthz`: Health check endpoint
- `/readyz`: Readiness check endpoint
- `/version`: Returns version information

## Project Structure

```
.
├── main.go                 # Application entry point
├── Dockerfile             # Container build instructions
├── go.mod                # Go module definition
├── pkg/
│   ├── config/          # Configuration handling
│   ├── digitalocean/    # DigitalOcean API integration
│   ├── health/         # Health check endpoints
│   ├── logging/        # Logging configuration
│   └── version/        # Version information
└── charts/              # Helm chart for Kubernetes deployment
```

## Dependencies

- Go 1.21.6
- `github.com/sirupsen/logrus`: Structured logging
- `github.com/stretchr/testify`: Testing framework
- Additional indirect dependencies managed through Go modules

## Development

### Running Tests

```bash
go test -v ./...
```

### Building from Source

```bash
make build
```

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

## Maintainers

- Matt Mattox (mmattox@support.tools)

## Links

- [GitHub Repository](https://github.com/supporttools/do-database-metrics-adapter)
- [Support Tools](https://support.tools)
