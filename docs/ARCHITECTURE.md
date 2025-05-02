 # Journey Hub Project Architecture

This document describes the overall architecture of the Journey Hub project. The purpose of this document is to provide an overview of how the project is structured and how different components interact with each other.

## Hexagonal Architecture

The Journey Hub project uses Hexagonal Architecture (or Ports and Adapters Architecture). This architecture is based on the following principles:

1. **Separation of Concerns**: Each part of the system has a specific responsibility
2. **Dependency on Abstractions, not Implementations**: Inner layers depend on interfaces, not specific implementations
3. **Testability**: Ability to test each layer separately and independently
4. **Replaceability**: Ability to replace external implementations without affecting the core business logic

## Architecture Layers

![Hexagonal Architecture](https://miro.medium.com/max/1400/1*LF3qzk0dgk9kHVVzxOcXEA.png)

### 1. Domain Core

The heart of the application is located in the `internal/domain` folder:

- **Entities**: Core business models such as `User`, `Journey`, `Booking`
- **Business Rules**: Pure business logic
- **Value Objects**: Immutable objects such as `Email`, `Password`, `Money`
- **Errors**: Domain-specific errors
- **Repository Interfaces**: Interfaces for storing and retrieving entities

This layer has no dependencies on outer layers and is independent of any specific technology.

### 2. Usecases

The usecases layer is located in the `internal/usecases` folder:

- **Services**: Implementation of system use cases
- **Output Ports**: Interfaces for external systems
- **Orchestration**: Coordination between entities and different operations

This layer implements application logic without dependency on technical details.

### 3. Primary Adapters

Primary adapters are located in the `internal/interfaces` folder:

- **HTTP Controllers**: RESTful API implementations
- **DTO Models**: Data transfer classes for API input and output
- **Middleware**: Such as authentication, logging, and access control

These adapters translate external requests into usecase calls.

### 4. Secondary Adapters

Secondary adapters are located in the `internal/infrastructure` folder:

- **Repository Implementations**: Data storage and retrieval (e.g., with MongoDB)
- **External Clients**: Communication with external APIs
- **Cache Systems**: Such as Redis
- **File Storage Systems**: Such as MinIO or S3

These adapters implement the output ports defined in the usecases layer.

### 5. Common Packages

Shared libraries are located in the `pkg` folder:

- **Configuration**: Settings management
- **Logging**: Logging system
- **Error Management**: Common error structures
- **Tools and Helpers**: General utility functions

## Data Flow Path

1. HTTP request enters the primary adapter
2. Input data is converted to domain models
3. A usecase is invoked
4. The usecase works with domain entities
5. I/O operations are performed through output ports and corresponding adapters
6. The result is returned to the primary adapter
7. The primary adapter converts the result to an appropriate format (e.g., JSON) and sends it

## Folder Structure

```
journey-hub/
├── .github/                     # CI/CD and GitHub Actions
├── backend/                     # Backend
│   ├── cmd/                     # Application entry points
│   │   └── api/                 # API server entry point
│   ├── config/                  # Configuration files
│   │   ├── development/         # Development environment configuration
│   │   ├── production/          # Production environment configuration
│   │   └── test/                # Test environment configuration
│   ├── internal/                # Project-specific code
│   │   ├── domain/              # Domain core
│   │   │   ├── user/            # User entity
│   │   │   ├── journey/         # Journey entity
│   │   │   └── booking/         # Booking entity
│   │   ├── usecases/            # Use cases
│   │   │   ├── user/            # User services
│   │   │   ├── journey/         # Journey services
│   │   │   └── booking/         # Booking services
│   │   ├── interfaces/          # Primary adapters
│   │   │   └── api/             # API interfaces
│   │   │       └── http/        # HTTP APIs
│   │   │           ├── dto/     # Data transfer models
│   │   │           ├── handlers/# Controllers
│   │   │           ├── middleware/ # Middleware
│   │   │           └── routes/  # API routes
│   │   └── infrastructure/      # Secondary adapters
│   │       ├── mongodb/         # Repository implementation with MongoDB
│   │       ├── redis/           # Cache implementation with Redis
│   │       └── minio/           # File storage implementation
│   ├── pkg/                     # Shared libraries
│   │   ├── logger/              # Logging module
│   │   ├── errors/              # Error management
│   │   └── auth/                # Authentication tools
│   ├── api/                     # API documentation
│   │   └── swagger/             # Swagger files
│   ├── scripts/                 # Useful scripts
│   └── tests/                   # Integration tests
├── frontend/                    # React frontend
├── docker-compose.yml           # Docker configuration
└── docs/                        # Project documentation
```

## Database Interaction

We use MongoDB as the main database:

1. **Repository Ports**: Defined in the domain layer
2. **Repository Implementation**: Implemented in the infrastructure layer
3. **Transaction Operations**: Performed in a separate session
4. **Indexing**: Based on common query patterns

## Caching

We use Redis for caching:

1. **Frequently Used Data Caching**: Such as journey information or search results
2. **Authentication Token Caching**: Storage of JWT tokens
3. **Session Management**: Storage of user session information
4. **Invalidation Strategy**: Using TTL and manual invalidation

## Error Management

The error management system includes:

1. **Domain Errors**: Defined in the domain layer
2. **Application Errors**: Defined in the usecases layer
3. **Infrastructure Errors**: Defined in the infrastructure layer
4. **Error Translation**: Converting internal errors to appropriate HTTP responses

## Logging Management

Our logging system supports:

1. **Various Log Levels**: Debug, Info, Warn, Error, Fatal
2. **Different Output Formats**: JSON for production, text for development
3. **Metadata Logging**: Request information, user, time, etc.
4. **Log Rotation**: Management of large log files
5. **Alert System**: Sending alerts for critical errors

## Authentication and Authorization

1. **JWT**: Using JWT tokens for authentication
2. **RBAC**: Role-based access control
3. **Email Verification**: User email verification process
4. **Password Recovery**: Forgotten password recovery process
5. **Rate Limiting**: Limiting the number of failed login attempts

## Architecture FAQs

### 1. Why are we using Hexagonal Architecture?

Hexagonal Architecture allows us to:

- Separate business logic from technical details
- Test different parts independently
- Easily change infrastructure technologies
- Enable parallel development of different system parts

### 2. How do we add a new feature?

To add a new feature, follow these steps:

1. Define entities and business rules in the domain layer
2. Define necessary ports
3. Implement use cases in the usecases layer
4. Create input adapters (APIs) and output adapters (e.g., repository implementation)
5. Cover all parts with tests

### 3. How do we prevent crossing architecture boundaries?

- Use one-way dependencies (inner layers have no knowledge of outer layers)
- Use interfaces to define contracts between layers
- Use DTOs for data transfer between layers
- Use code reviews to ensure boundaries are respected

## Further Resources

- [Hexagonal Architecture (Alistair Cockburn)](https://alistair.cockburn.us/hexagonal-architecture/)
- [Domain-Driven Design (Eric Evans)](https://domainlanguage.com/ddd/)
- [Clean Architecture (Robert C. Martin)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
