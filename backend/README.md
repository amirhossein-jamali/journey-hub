# Journey Hub Backend

Backend service for the Journey Hub platform, built with a hexagonal architecture using Go.

## Getting Started

### Prerequisites
- Go 1.20 or newer
- Docker and Docker Compose
- golangci-lint
- swag (for Swagger generation)

### Development Environment Setup

#### Option 1: Using VS Code Dev Containers
1. Install the [Remote - Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension in VS Code
2. Open the project in VS Code
3. Press F1, select "Remote-Containers: Reopen in Container"
4. VS Code will build the dev container and open the project inside it

#### Option 2: Using Docker Compose
1. Build and start the development environment:
   ```bash
   docker-compose -f docker-compose.dev.yml up -d
   ```
2. Access the backend service at http://localhost:8080
3. Access Swagger UI at http://localhost:8081

#### Option 3: Local Development
1. Install dependencies:
   ```bash
   go mod download
   go mod verify
   ```
2. Install Git hooks:
   ```bash
   chmod +x backend/scripts/install-hooks.sh
   ./backend/scripts/install-hooks.sh
   ```
3. Run the application:
   ```bash
   cd backend
   go run cmd/app/main.go
   ```

### Project Structure
```
backend/
├── api/              # API documentation (Swagger)
├── cmd/              # Application entry points
├── config/           # Configuration files
├── internal/         # Private application code
│   ├── domain/       # Domain layer (entities, repositories interfaces)
│   ├── usecases/     # Application layer (use cases, business logic)
│   ├── interfaces/   # Interface adapters (controllers, presenters)
│   └── infrastructure/ # Infrastructure layer (DB, external services)
├── scripts/          # Build and development scripts
└── tests/            # Integration and end-to-end tests
```

### Hexagonal Architecture
This project follows the Hexagonal Architecture (Ports and Adapters) pattern:

- **Domain Layer**: Core business logic and entities
- **Application Layer**: Use cases that orchestrate the domain
- **Interface Adapters**: Controllers, presenters, and gateways
- **Infrastructure**: External systems and frameworks

## Development Workflows

### Generating Swagger Documentation
```bash
cd backend
swag init -g cmd/app/main.go -o api/swagger
```

### Running Linters
```bash
cd backend
golangci-lint run
```

### Running Tests
```bash
cd backend
go test -v ./...
```

### Config Encryption System
Journey Hub uses an encryption system to protect sensitive configuration values:

- For development: Run `./scripts/setup-encrypted-configs.sh` to set up encrypted values
- For CI/CD: Set `JH_SKIP_CONFIG_ENCRYPTION=true` to disable encryption
- For production: Set `JH_CONFIG_ENCRYPTION_KEY` to a secure key
- For more details, see [Config Encryption Guide](config/ENCRYPTION.md)

## Docker Support

### Building the Docker Image
```bash
docker build -t journey-hub-backend .
```

### Running the Docker Container
```bash
docker run -p 8080:8080 journey-hub-backend
```

## CI/CD Pipeline

The project is set up with the following pipeline stages:
1. Code linting
2. Unit tests
3. Integration tests
4. Build and package
5. Deploy to staging/production

## Contributing
1. Create a feature branch from the main branch
2. Make your changes
3. Run tests and linters
4. Submit a pull request

All contributions must pass the pre-commit hooks for code quality and formatting. 