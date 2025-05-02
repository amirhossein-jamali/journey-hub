# Journey Hub

![Version](https://img.shields.io/badge/version-0.1.0-blue.svg)
![Build](https://img.shields.io/badge/build-passing-brightgreen.svg)
![Coverage](https://img.shields.io/badge/coverage-85%25-brightgreen.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Docker](https://img.shields.io/badge/docker-ready-blue.svg)
![CI/CD](https://img.shields.io/badge/CI%2FCD-GitHub%20Actions-blue.svg)

<div align="center">
  <img src="docs/images/journey-hub-screenshot.png" alt="Journey Hub UI" width="800px">
  <!-- Image will be added later -->
</div>

Journey Hub is a travel management platform built with Hexagonal Architecture, Domain-Driven Design, and Clean Architecture principles. This project aims to create a scalable, flexible, and testable system for managing journeys and bookings.

## Feature Status

- [x] Hexagonal Architecture & DDD implementation
- [x] Project structure with separation of concerns
- [x] User and profile management
- [x] Journey creation and management
- [x] Basic reservation system
- [ ] Online payment
- [ ] Rating and review system
- [ ] Analytics dashboard
- [ ] Mobile application
- [ ] Multi-language support

## About Journey Hub

Journey Hub is a travel management platform offering the following features:
- User and profile management
- Journey creation and management
- Reservation system
- Ratings and reviews
- Journey search and filtering
- Admin dashboard

This project uses modern technologies to provide a scalable, high-performance solution:
- Backend: Go with Hexagonal Architecture
- Frontend: React with TypeScript
- Database: MongoDB
- Cache: Redis
- Monitoring: Prometheus
- Logging: Zap
- API Documentation: Swagger

## Key Technologies

| Component     | Technology                           |
|---------------|-------------------------------------|
| Backend       | Go (Gin, gqlgen, mongo-driver)      |
| Frontend      | React + TypeScript                  |
| Database      | MongoDB                             |
| Cache         | Redis                               |
| Authentication| JWT / OAuth2                        |
| Authorization | Casbin / RBAC                       |
| Monitoring    | Prometheus                          |
| Logging       | Zap / Logrus                        |
| Configuration | Viper / env                         |
| Testing       | testify, gomock, testcontainers     |
| API Docs      | Swagger / OpenAPI                   |
| External Comm | net/http, resty                     |

## Architecture Overview

This project follows Hexagonal Architecture (Ports & Adapters) with DDD and Clean Architecture approaches to create a modular, testable, and extensible structure.

### Key Principles

- **Separation of Concerns**: Business domain is separated from technical details.
- **Dependency Inward**: Outer layers depend on inner layers, not vice versa.
- **Independence from Technical Details**: Domain core is independent of infrastructure details like database, API, and UI.
- **Replaceability**: Adapters can be replaced without changing business logic.
- **Domain-Driven Design**: Software structure is shaped by the business domain model.

### Architecture Diagram

```
                              ┌─────────────────────────────────┐
                              │         UI / Clients            │
                              │   ┌─────────┐   ┌─────────┐     │
                              │   │  React  │   │ Mobile  │     │
                              │   │Frontend │   │   App   │     │
                              │   └────┬────┘   └────┬────┘     │
                              └────────┼────────────┼───────────┘
                                       │            │
                                       ▼            ▼
     ┌────────────────────────────────────────────────────────────────────┐
     │                     Primary Adapters (Input)                       │
     │  ┌─────────────────────────┐       ┌─────────────────────────┐     │
     │  │       REST API          │       │        GraphQL          │     │
     │  │      Controllers        │       │       Resolvers         │     │
     │  │      (Gin/Echo)         │       │     (gqlgen/graph)      │     │
     │  └────────────┬────────────┘       └────────────┬────────────┘     │
     └───────────────┼─────────────────────────────────┼──────────────────┘
                     │                                 │
                     │    implements                   │    implements
                     ▼                                 ▼
     ┌─────────────────────────────────────────────────────────────────────┐
     │                        Primary Ports (Input)                        │
     │  ┌─────────────────────────────────────────────────────────────┐    │
     │  │                 Application Services Interfaces             │    │
     │  │                                                             │    │
     │  │  ┌─────────────┐     ┌──────────────┐     ┌──────────────┐  │    │
     │  │  │UserService  │     │JourneyService│     │BookingService│  │    │
     │  │  │  Interface  │     │  Interface   │     │  Interface   │  │    │
     │  │  └──────┬──────┘     └──────┬───────┘     └──────┬───────┘  │    │
     │  └─────────┼───────────────────────────────────────────────────┘    │
     └─────────────────────────────┬─────────────────┬─────────────────────┘
                                   │                 │
                                   │ uses            │ uses
                                   ▼                 ▼
┌───────────────────────────────────────────────────────────────────────────┐
│                             Domain Core                                    │
│                                                                           │
│  ┌───────────────────────────────────────────────────────────────┐       │
│  │                           Entities                            │       │
│  │  ┌─────────┐           ┌─────────┐           ┌─────────┐     │       │
│  │  │  User   │           │ Journey │           │ Booking │     │       │
│  │  └─────────┘           └─────────┘           └─────────┘     │       │
│  └───────────────────────────────────────────────────────────────┘       │
│                                                                           │
│  ┌───────────────────────────────────────────────────────────────┐       │
│  │                       Domain Services                         │       │
│  │  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐  │       │
│  │  │UserDomainSvc   │  │JourneyDomainSvc│  │BookingDomainSvc│  │       │
│  │  └────────────────┘  └────────────────┘  └────────────────┘  │       │
│  └───────────────────────────────────────────────────────────────┘       │
│                                                                           │
│  ┌───────────────────────────────────────────────────────────────┐       │
│  │                    Value Objects & Rules                       │       │
│  │  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐  │       │
│  │  │UserValidations │  │JourneyRules    │  │BookingPolicies │  │       │
│  │  └────────────────┘  └────────────────┘  └────────────────┘  │       │
│  └───────────────────────────────────────────────────────────────┘       │
└───────────────────────────────────────────────────────────────────────────┘
                  ▲                      │                      ▲
                  │                      │                      │
                  │                      ▼                      │
     ┌────────────────────────────────────────────────────────────────────┐
     │                      Secondary Ports (Output)                       │
     │  ┌─────────────────────────────────────────────────────────────┐   │
     │  │                   Repository Interfaces                      │   │
     │  │                                                             │   │
     │  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐     │   │
     │  │  │  User Repo  │    │Journey Repo │    │ Booking Repo│     │   │
     │  │  │  Interface  │    │  Interface  │    │  Interface  │     │   │
     │  │  └──────┬──────┘    └──────┬──────┘    └──────┬──────┘     │   │
     │  └─────────┼─────────────────────────────────────────────────┘   │
     └─────────────────────────────┬─────────────────┬──────────────────┘
                    implemented by │                 │ implemented by
                                  ▼                 ▼
     ┌────────────────────────────────────────────────────────────────────┐
     │                    Secondary Adapters (Output)                      │
     │                                                                    │
     │  ┌─────────────────────────┐      ┌────────────────────────┐      │
     │  │  MongoDB Repository     │      │   Redis Cache          │      │
     │  │  (mongo-driver/mgo)    │      │   (go-redis/redigo)    │      │
     │  └─────────────────────────┘      └────────────────────────┘      │
     │                                                                    │
     │  ┌─────────────────────────┐      ┌────────────────────────┐      │
     │  │  External APIs          │      │  Message Queue         │      │
     │  │  (net/http/resty)      │      │  (kafka/rabbitmq)      │      │
     │  └─────────────────────────┘      └────────────────────────┘      │
     └────────────────────────────────────────────────────────────────────┘
                                  ▲
                                  │
┌─────────────────────────────────┴─────────────────────────────────────┐
│                       Testing & Verification                           │
│                                                                       │
│  ┌───────────────────────────────────────────────────────────────┐    │
│  │                       Unit Tests                              │    │
│  │  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐  │    │
│  │  │Domain Tests   │  │Service Tests   │  │Use Case Tests  │  │    │
│  │  │(testify/gomock)│  │(testify/gomock)│  │(testify/gomock)│  │    │
│  │  └────────────────┘  └────────────────┘  └────────────────┘  │    │
│  └───────────────────────────────────────────────────────────────┘    │
│                                                                       │
│  ┌───────────────────────────────────────────────────────────────┐    │
│  │                    Integration Tests                          │    │
│  │  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐  │    │
│  │  │API Tests      │  │Repository Tests│  │End-to-End Tests│  │    │
│  │  │(httptest)     │  │(testcontainers)│  │(cypress/ginkgo)│  │    │
│  │  └────────────────┘  └────────────────┘  └────────────────┘  │    │
│  └───────────────────────────────────────────────────────────────┘    │
└───────────────────────────────────────────────────────────────────────┘

┌───────────────────────────────────────────────────────────────────────┐
│                      Cross-Cutting Concerns                           │
│                                                                       │
│  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐  │
│  │  Logging       │  │ Authentication │  │ Authorization  │  │ Monitoring     │  │
│  │  (Zap/logrus) │  │ (JWT/OAuth2)   │  │ (Casbin/RBAC)  │  │ (Prometheus)   │  │
│  └────────────────┘  └────────────────┘  └────────────────┘  └────────────────┘  │
│                                                                       │
│  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐  │
│  │ Configuration  │  │ Error Handling │  │ Transactions   │  │ Caching        │  │
│  │ (Viper/env)   │  │ (pkg/errors)   │  │ (context/sync) │  │ (TTL/LRU)      │  │
│  └────────────────┘  └────────────────┘  └────────────────┘  └────────────────┘  │
└───────────────────────────────────────────────────────────────────────┘
```

## Project Structure

The project structure is modular to support Hexagonal Architecture, Domain-Driven Design, and Clean Architecture principles:

```
journey-hub/                           # Project root
├── backend/                           # Backend code
│   ├── cmd/                           # Entry points
│   │   └── api/                       # Main API service
│   │       └── main.go                # Application entry point
│   ├── internal/                      # Project-specific code
│   │   ├── domain/                    # Domain layer (business core)
│   │   │   ├── user/                  # Domain entity
│   │   │   │   ├── entity.go          # User entity
│   │   │   │   ├── repository.go      # User repository interface
│   │   │   │   └── service.go         # User service interface
│   │   │   │
│   │   │   ├── journey/               # Journey domain
│   │   │   │   ├── entity.go          # Journey entity
│   │   │   │   ├── repository.go      # Journey repository interface
│   │   │   │   └── service.go         # Journey service interface
│   │   │   │
│   │   │   ├── booking/               # Booking domain
│   │   │   │   ├── entity.go          # Booking entity
│   │   │   │   ├── repository.go      # Booking repository interface
│   │   │   │   └── service.go         # Booking service interface
│   │   │   │
│   │   │   ├── policy/                # Authorization policies
│   │   │   │   ├── user_policy.go     # User authorization policy
│   │   │   │   ├── journey_policy.go  # Journey authorization policy
│   │   │   │   └── booking_policy.go  # Booking authorization policy
│   │   │   │
│   │   │   └── common/                # Shared domain elements
│   │   │       ├── errors.go          # Domain-specific errors
│   │   │       └── types.go           # Shared types
│   │   │
│   │   ├── usecases/                  # Use cases layer
│   │   │   ├── user/                  # User use cases
│   │   │   │   └── service_impl.go    # User service implementation
│   │   │   │
│   │   │   ├── journey/               # Journey use cases
│   │   │   │   └── service_impl.go    # Journey service implementation
│   │   │   │
│   │   │   ├── booking/               # Booking use cases
│   │   │   │   └── service_impl.go    # Booking service implementation
│   │   │   │
│   │   │   └── auth/                  # Authentication use cases
│   │   │       └── service_impl.go    # Authentication service implementation
│   │   │
│   │   ├── interfaces/                # Interface adapters layer
│   │   │   ├── api/                   # API interface adapters
│   │   │   │   ├── http/              # HTTP interface adapters
│   │   │   │   │   ├── handlers/      # HTTP handlers
│   │   │   │   │   │   ├── user_handler.go
│   │   │   │   │   │   ├── journey_handler.go
│   │   │   │   │   │   ├── booking_handler.go
│   │   │   │   │   │   └── auth_handler.go
│   │   │   │   │   │
│   │   │   │   │   ├── dto/           # Data transfer objects for HTTP
│   │   │   │   │   │   ├── user_dto.go
│   │   │   │   │   │   ├── journey_dto.go
│   │   │   │   │   │   ├── booking_dto.go
│   │   │   │   │   │   └── auth_dto.go
│   │   │   │   │   │
│   │   │   │   │   └── router.go      # API router
│   │   │   │   │
│   │   │   │   └── grpc/              # (future) gRPC interface adapters
│   │   │   │
│   │   │   └── persistence/           # Persistence interface adapters
│   │   │       ├── mongodb/           # MongoDB interface adapters
│   │   │       │   ├── models/        # Database models
│   │   │       │   │   ├── user_model.go
│   │   │       │   │   ├── journey_model.go
│   │   │       │   │   └── booking_model.go
│   │   │       │   │
│   │   │       │   └── repositories/  # Repository interfaces
│   │   │       │       ├── user_repository.go
│   │   │       │       ├── journey_repository.go
│   │   │       │       └── booking_repository.go
│   │   │       │
│   │   │       └── cache/             # Cache interface adapters
│   │   │           ├── redis_repository.go  # Redis interface adapter
│   │   │           └── memory_repository.go # Memory interface adapter
│   │   │
│   │   ├── infrastructure/            # Infrastructure layer
│   │   │   ├── middleware/            # Shared middleware
│   │   │   │   ├── auth_middleware.go
│   │   │   │   ├── logging_middleware.go
│   │   │   │   └── error_middleware.go
│   │   │   │
│   │   │   ├── server/                # Server layer
│   │   │   │   ├── http_server.go
│   │   │   │   └── server.go
│   │   │   │
│   │   │   ├── database/              # Database layer
│   │   │   │   ├── mongodb/
│   │   │   │   │   └── connection.go  # MongoDB connection
│   │   │   │   └── cache/
│   │   │   │       └── redis.go       # Redis connection
│   │   │   │
│   │   │   └── config/                # Configuration layer
│   │   │       ├── app.go             # Application configuration
│   │   │       ├── database.go        # Database configuration
│   │   │       └── server.go          # Server configuration
│   │   │
│   │   └── pkg/                       # Shared packages
│   │       ├── logger/                # Logging
│   │       │   └── zap_logger.go      # Zap logger implementation
│   │       │
│   │       ├── validator/             # Validation
│   │       │   └── validator.go
│   │       │
│   │       └── metrics/               # Metrics
│   │           └── prometheus.go      # Prometheus implementation
│   │
│   ├── api/                           # API documentation
│   │   └── swagger/                   # Swagger API documentation
│   │       └── swagger.yaml
│   │
│   ├── config/                        # Configuration files
│   │   ├── development.yaml
│   │   ├── test.yaml
│   │   └── production.yaml
│   │
│   ├── tests/                         # Tests
│   │   ├── unit/                      # Unit tests
│   │   │   ├── domain/                # Domain layer tests
│   │   │   │   ├── user/
│   │   │   │   ├── journey/
│   │   │   │   └── booking/
│   │   │   │
│   │   │   ├── usecases/              # Use cases layer tests
│   │   │   │   ├── user/
│   │   │   │   ├── journey/
│   │   │   │   └── booking/
│   │   │   │
│   │   │   └── interfaces/            # Interface adapters layer tests
│   │   │       ├── api/
│   │   │       └── persistence/
│   │   │
│   │   └── integration/               # Integration tests
│   │       ├── api/                   # API tests
│   │       └── repository/            # Repository tests
│   │
│   ├── scripts/                       # Utility scripts
│   │   ├── migrations/                # Database migrations
│   │   └── setup.sh                   # Development environment setup
│   │
│   └── go.mod                         # Go dependencies
│
├── frontend/                          # React frontend code
│   ├── public/                        # Public files
│   │   ├── index.html
│   │   └── favicon.ico
│   │
│   ├── src/                           # React source code
│   │   ├── components/                # Reusable components
│   │   │   ├── common/                # Shared components
│   │   │   │   ├── Button/
│   │   │   │   ├── Input/
│   │   │   │   └── Card/
│   │   │   │
│   │   │   └── layout/                # Layout components
│   │   │       ├── Header/
│   │   │       └── Footer/
│   │   │
│   │   ├── features/                  # Main features (based on domain)
│   │   │   ├── user/                  # User feature
│   │   │   │   ├── components/        # User-specific components
│   │   │   │   ├── pages/             # User pages
│   │   │   │   └── services/          # User API services
│   │   │   │
│   │   │   ├── journey/               # Journey feature
│   │   │   │   ├── components/
│   │   │   │   ├── pages/
│   │   │   │   └── services/
│   │   │   │
│   │   │   ├── booking/               # Booking feature
│   │   │   │   ├── components/
│   │   │   │   ├── pages/
│   │   │   │   └── services/
│   │   │   │
│   │   │   └── auth/                  # Authentication feature
│   │   │       ├── components/
│   │   │       ├── pages/
│   │   │       └── services/
│   │   │
│   │   ├── hooks/                     # Custom hooks
│   │   │   ├── useAuth.js
│   │   │   └── useForm.js
│   │   │
│   │   ├── services/                  # Shared API services
│   │   │   └── api.js                 # Base API configuration
│   │   │
│   │   ├── utils/                     # Helper functions
│   │   │   ├── helpers.js
│   │   │   └── formatter.js
│   │   │
│   │   ├── App.js                     # Main component
│   │   └── index.js                   # Entry point
│   │
│   ├── package.json                   # npm dependencies
│   └── .env                           # Environment variables
│
├── docker-compose.yml                 # docker-compose configuration
├── Dockerfile.backend                 # Backend Dockerfile
├── Dockerfile.frontend                # Frontend Dockerfile
│
├── docs/                              # Project documentation
│   ├── architecture.md                # Detailed architecture
│   ├── architecture.png               # Architecture diagram
│   ├── development.md                 # Development guide
│   └── api.md                         # API documentation
│
└── README.md                          # Main project documentation
```

## Benefits of Hexagonal Architecture

1. **Separation of Concerns**: Business domain is separated from technical details.
2. **Changeability**: We can replace adapters without changing the domain core.
3. **Testability**: Layers can be easily tested with mocks.
4. **Extensibility**: Different teams can work on different layers simultaneously.
5. **Maintainability**: Code is more organized and maintainable.
6. **Scalability**: Independent scaling of different parts of the system
7. **Increased Productivity**: Reduced dependencies lead to increased development team productivity
8. **Reduced Technical Risk**: Freedom to change technologies without affecting business logic

## Getting Started

### Prerequisites
- Go 1.18+
- Docker and Docker Compose
- MongoDB
- Redis
- Node.js 16+ and npm (for frontend)
- Git

### Set Up Development Environment

#### 1. Clone the Repository
```bash
git clone https://github.com/yourusername/journey-hub.git
cd journey-hub
```

#### 2. Start Databases with Docker
```bash
cd backend/deployments/docker
docker-compose up -d mongodb redis
```

#### 3. Set Environment Variables
```bash
cp backend/.env.example backend/.env
# Edit .env file with required variables
```

#### 4. Run Backend Server
```bash
cd backend
go mod download
go run cmd/server/main.go
```

#### 5. Start Frontend
```bash
cd frontend
npm install
npm start
```

### Run with Docker Compose (Full System)
```bash
docker-compose up -d
```

### Accessing Services
- Backend Server: `http://localhost:8080/api/v1`
- API Documentation: `http://localhost:8080/api/v1/docs`
- Frontend: `http://localhost:3000`
- Prometheus Monitoring: `http://localhost:9090`

## Documentation
For more information, see the `docs/` folder. Documentation includes:
- [Detailed Architecture](docs/architecture.md)
- [Development Guide](docs/development.md)
- [API Documentation](docs/api.md) and [Swagger API](api/swagger/swagger.yaml)
- [Installation Guide](docs/installation.md)
- [Architecture Decisions](docs/architecture-decisions.md)

## Contributing

To contribute to this project, please see the `CONTRIBUTING.md` file. We welcome all types of contributions:
- Bug reports
- Feature suggestions
- Pull Requests
- Documentation improvements
- Feedback sharing

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

<div dir="rtl">

# Journey Hub

![Version](https://img.shields.io/badge/version-0.1.0-blue.svg)
![Build](https://img.shields.io/badge/build-passing-brightgreen.svg)
![Coverage](https://img.shields.io/badge/coverage-85%25-brightgreen.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Docker](https://img.shields.io/badge/docker-ready-blue.svg)
![CI/CD](https://img.shields.io/badge/CI%2FCD-GitHub%20Actions-blue.svg)

Journey Hub یک پلتفرم مدیریت سفر با معماری هگزاگونال، Domain-Driven Design و Clean Architecture است. این پروژه با هدف ایجاد یک سیستم قابل توسعه، انعطاف‌پذیر و تست‌پذیر برای مدیریت سفرها و رزروها طراحی شده است.

<div align="center">
  <img src="docs/images/journey-hub-screenshot.png" alt="Journey Hub UI" width="800px">
  <!-- Image will be added later -->
</div>

## وضعیت ویژگی‌ها

- [x] معماری هگزاگونال و DDD
- [x] ساختار پروژه با تفکیک دغدغه‌ها
- [x] مدیریت کاربران و پروفایل‌ها
- [x] ایجاد و مدیریت سفرها
- [x] سیستم رزرواسیون پایه
- [ ] پرداخت آنلاین
- [ ] سیستم امتیازدهی و نظرات
- [ ] داشبورد تحلیلی
- [ ] اپلیکیشن موبایل
- [ ] پشتیبانی از چند زبان

## درباره Journey Hub

Journey Hub یک پلتفرم برای مدیریت سفر است که امکانات زیر را فراهم می‌کند:
- مدیریت کاربران و پروفایل‌ها
- ایجاد و مدیریت سفرها
- سیستم رزرواسیون
- امتیازدهی و نظرات
- جستجو و فیلتر کردن سفرها
- داشبورد مدیریتی

این پروژه از فناوری‌های مدرن برای ارائه یک راه‌حل مقیاس‌پذیر و عملکرد بالا استفاده می‌کند:
- بک‌اند: Go با معماری هگزاگونال
- فرانت‌اند: React با TypeScript
- پایگاه داده: MongoDB
- کش: Redis
- مانیتورینگ: Prometheus
- لاگینگ: Zap
- مستندسازی API: Swagger

## تکنولوژی‌های کلیدی

| بخش            | فناوری مورد استفاده                   |
|----------------|--------------------------------------|
| بک‌اند          | Go (Gin, gqlgen, mongo-driver)       |
| فرانت‌اند       | React + TypeScript                   |
| پایگاه داده     | MongoDB                              |
| کش             | Redis                                |
| احراز هویت      | JWT / OAuth2                         |
| مجوزدهی        | Casbin / RBAC                        |
| مانیتورینگ     | Prometheus                           |
| لاگینگ          | Zap / Logrus                        |
| مدیریت تنظیمات  | Viper / env                         |
| تست            | testify, gomock, testcontainers     |
| API مستندسازی   | Swagger / OpenAPI                   |
| ارتباط خارجی    | net/http, resty                     |

## چشم‌انداز معماری

این پروژه از معماری هگزاگونال (Ports & Adapters) با رویکرد DDD و Clean Architecture پیروی می‌کند تا یک ساختار ماژولار، قابل آزمون و قابل توسعه ایجاد کند.

### اصول کلیدی

- **جداسازی دغدغه‌ها**: دامنه کسب‌وکار از جزئیات فنی جدا شده است.
- **وابستگی رو به داخل**: لایه‌های بیرونی به لایه‌های درونی وابسته‌اند، نه برعکس.
- **استقلال از جزئیات فنی**: هسته دامنه از جزئیات زیرساخت مانند دیتابیس، API و رابط کاربری مستقل است.
- **قابلیت جایگزینی**: آداپتورها بدون تغییر در منطق کسب‌وکار قابل جایگزینی هستند.
- **طراحی بر اساس دامنه**: ساختار نرم‌افزار بر اساس مدل دامنه کسب‌وکار شکل گرفته است.

### نمودار معماری

```
                              ┌─────────────────────────────────┐
                              │         UI / Clients            │
                              │   ┌─────────┐   ┌─────────┐    │
                              │   │  React  │   │ Mobile  │    │
                              │   │Frontend │   │   App   │    │
                              │   └────┬────┘   └────┬────┘    │
                              └────────┼────────────┼───────────┘
                                       │            │
                                       ▼            ▼
     ┌────────────────────────────────────────────────────────────────────┐
     │                     Primary Adapters (Input)                        │
     │  ┌─────────────────────────┐       ┌─────────────────────────┐     │
     │  │       REST API          │       │        GraphQL          │     │
     │  │      Controllers        │       │       Resolvers         │     │
     │  │      (Gin/Echo)        │       │     (gqlgen/graph)      │     │
     │  └────────────┬────────────┘       └────────────┬────────────┘     │
     └───────────────┼─────────────────────────────────┼──────────────────┘
                     │                                 │
                     │    implements                   │    implements
                     ▼                                 ▼
     ┌────────────────────────────────────────────────────────────────────┐
     │                        Primary Ports (Input)                        │
     │  ┌────────────────────────────────────────────────────────────┐    │
     │  │                 Application Services Interfaces             │    │
     │  │                                                            │    │
     │  │  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐  │    │
     │  │  │UserService  │     │JourneyService│     │BookingService│  │    │
     │  │  │  Interface  │     │  Interface   │     │  Interface   │  │    │
     │  │  └──────┬──────┘     └──────┬──────┘     └──────┬──────┘  │    │
     │  └─────────┼─────────────────────────────────────────────────┘    │
     └─────────────────────────────┬─────────────────┬──────────────────┘
                                  │                 │
                                  │ uses            │ uses
                                  ▼                 ▼
┌───────────────────────────────────────────────────────────────────────────┐
│                             Domain Core                                    │
│                                                                           │
│  ┌───────────────────────────────────────────────────────────────┐       │
│  │                           Entities                            │       │
│  │  ┌─────────┐           ┌─────────┐           ┌─────────┐     │       │
│  │  │  User   │           │ Journey │           │ Booking │     │       │
│  │  └─────────┘           └─────────┘           └─────────┘     │       │
│  └───────────────────────────────────────────────────────────────┘       │
│                                                                           │
│  ┌───────────────────────────────────────────────────────────────┐       │
│  │                       Domain Services                         │       │
│  │  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐  │       │
│  │  │UserDomainSvc   │  │JourneyDomainSvc│  │BookingDomainSvc│  │       │
│  │  └────────────────┘  └────────────────┘  └────────────────┘  │       │
│  └───────────────────────────────────────────────────────────────┘       │
│                                                                           │
│  ┌───────────────────────────────────────────────────────────────┐       │
│  │                    Value Objects & Rules                       │       │
│  │  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐  │       │
│  │  │UserValidations │  │JourneyRules    │  │BookingPolicies │  │       │
│  │  └────────────────┘  └────────────────┘  └────────────────┘  │       │
│  └───────────────────────────────────────────────────────────────┘       │
└───────────────────────────────────────────────────────────────────────────┘
                  ▲                      │                      ▲
                  │                      │                      │
                  │                      ▼                      │
     ┌────────────────────────────────────────────────────────────────────┐
     │                      Secondary Ports (Output)                       │
     │  ┌─────────────────────────────────────────────────────────────┐   │
     │  │                   Repository Interfaces                      │   │
     │  │                                                             │   │
     │  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐     │   │
     │  │  │  User Repo  │    │Journey Repo │    │ Booking Repo│     │   │
     │  │  │  Interface  │    │  Interface  │    │  Interface  │     │   │
     │  │  └──────┬──────┘    └──────┬──────┘    └──────┬──────┘     │   │
     │  └─────────┼─────────────────────────────────────────────────┘   │
     └─────────────────────────────┬─────────────────┬──────────────────┘
                    implemented by │                 │ implemented by
                                  ▼                 ▼
     ┌────────────────────────────────────────────────────────────────────┐
     │                    Secondary Adapters (Output)                      │
     │                                                                    │
     │  ┌─────────────────────────┐      ┌────────────────────────┐      │
     │  │  MongoDB Repository     │      │   Redis Cache          │      │
     │  │  (mongo-driver/mgo)    │      │   (go-redis/redigo)    │      │
     │  └─────────────────────────┘      └────────────────────────┘      │
     │                                                                    │
     │  ┌─────────────────────────┐      ┌────────────────────────┐      │
     │  │  External APIs          │      │  Message Queue         │      │
     │  │  (net/http/resty)      │      │  (kafka/rabbitmq)      │      │
     │  └─────────────────────────┘      └────────────────────────┘      │
     └────────────────────────────────────────────────────────────────────┘
                                  ▲
                                  │
┌─────────────────────────────────┴─────────────────────────────────────┐
│                       Testing & Verification                           │
│                                                                       │
│  ┌───────────────────────────────────────────────────────────────┐    │
│  │                       Unit Tests                              │    │
│  │  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐  │    │
│  │  │Domain Tests   │  │Service Tests   │  │Use Case Tests  │  │    │
│  │  │(testify/gomock)│  │(testify/gomock)│  │(testify/gomock)│  │    │
│  │  └────────────────┘  └────────────────┘  └────────────────┘  │    │
│  └───────────────────────────────────────────────────────────────┘    │
│                                                                       │
│  ┌───────────────────────────────────────────────────────────────┐    │
│  │                    Integration Tests                          │    │
│  │  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐  │    │
│  │  │API Tests      │  │Repository Tests│  │End-to-End Tests│  │    │
│  │  │(httptest)     │  │(testcontainers)│  │(cypress/ginkgo)│  │    │
│  │  └────────────────┘  └────────────────┘  └────────────────┘  │    │
│  └───────────────────────────────────────────────────────────────┘    │
└───────────────────────────────────────────────────────────────────────┘

┌───────────────────────────────────────────────────────────────────────┐
│                      Cross-Cutting Concerns                           │
│                                                                       │
│  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐  │
│  │  Logging       │  │ Authentication │  │ Authorization  │  │ Monitoring     │  │
│  │  (Zap/logrus) │  │ (JWT/OAuth2)   │  │ (Casbin/RBAC)  │  │ (Prometheus)   │  │
│  └────────────────┘  └────────────────┘  └────────────────┘  └────────────────┘  │
│                                                                       │
│  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐  │
│  │ Configuration  │  │ Error Handling │  │ Transactions   │  │ Caching        │  │
│  │ (Viper/env)   │  │ (pkg/errors)   │  │ (context/sync) │  │ (TTL/LRU)      │  │
│  └────────────────┘  └────────────────┘  └────────────────┘  └────────────────┘  │
└───────────────────────────────────────────────────────────────────────┘
```

## ساختار پروژه

ساختار پروژه به صورت ماژولار برای پشتیبانی از اصول معماری هگزاگونال، Domain-Driven Design و Clean Architecture طراحی شده است:

```
journey-hub/                           # ریشه پروژه
├── backend/                           # کد بک‌اند
│   ├── cmd/                           # نقاط ورودی برنامه
│   │   └── api/                       # سرویس اصلی API
│   │       └── main.go                # نقطه شروع اجرای برنامه
│   │
│   ├── internal/                      # کد اختصاصی پروژه
│   │   ├── domain/                    # لایه دامنه (هسته کسب‌وکار)
│   │   │   ├── user/                  # دامنه کاربر
│   │   │   │   ├── entity.go          # موجودیت کاربر
│   │   │   │   ├── repository.go      # رابط مخزن کاربر (پورت ثانویه)
│   │   │   │   └── service.go         # رابط سرویس کاربر (پورت اولیه)
│   │   │   │
│   │   │   ├── journey/               # دامنه سفر
│   │   │   │   ├── entity.go          # موجودیت سفر
│   │   │   │   ├── repository.go      # رابط مخزن سفر
│   │   │   │   └── service.go         # رابط سرویس سفر
│   │   │   │
│   │   │   ├── booking/               # دامنه رزرو
│   │   │   │   ├── entity.go          # موجودیت رزرو
│   │   │   │   ├── repository.go      # رابط مخزن رزرو
│   │   │   │   └── service.go         # رابط سرویس رزرو
│   │   │   │
│   │   │   ├── policy/                # سیاست‌های دسترسی
│   │   │   │   ├── user_policy.go     # قوانین دسترسی کاربر
│   │   │   │   ├── journey_policy.go  # قوانین دسترسی سفر
│   │   │   │   └── booking_policy.go  # قوانین دسترسی رزرو
│   │   │   │
│   │   │   └── common/                # عناصر مشترک دامنه
│   │   │       ├── errors.go          # خطاهای دامنه
│   │   │       └── types.go           # تایپ‌های مشترک
│   │   │
│   │   ├── usecases/                  # لایه موارد استفاده (یوزکیس‌ها)
│   │   │   ├── user/                  # یوزکیس‌های کاربر
│   │   │   │   └── service_impl.go    # پیاده‌سازی سرویس کاربر
│   │   │   │
│   │   │   ├── journey/               # یوزکیس‌های سفر
│   │   │   │   └── service_impl.go    # پیاده‌سازی سرویس سفر
│   │   │   │
│   │   │   ├── booking/               # یوزکیس‌های رزرو
│   │   │   │   └── service_impl.go    # پیاده‌سازی سرویس رزرو
│   │   │   │
│   │   │   └── auth/                  # یوزکیس‌های احراز هویت
│   │   │       └── service_impl.go    # پیاده‌سازی سرویس احراز هویت
│   │   │
│   │   ├── interfaces/                # لایه آداپتورهای رابط
│   │   │   ├── api/                   # آداپتورهای ورودی API
│   │   │   │   ├── http/              # رابط HTTP
│   │   │   │   │   ├── handlers/      # هندلرهای HTTP
│   │   │   │   │   │   ├── user_handler.go
│   │   │   │   │   │   ├── journey_handler.go
│   │   │   │   │   │   ├── booking_handler.go
│   │   │   │   │   │   └── auth_handler.go
│   │   │   │   │   │
│   │   │   │   │   ├── dto/           # DTO‌ها برای HTTP
│   │   │   │   │   │   ├── user_dto.go
│   │   │   │   │   │   ├── journey_dto.go
│   │   │   │   │   │   ├── booking_dto.go
│   │   │   │   │   │   └── auth_dto.go
│   │   │   │   │   │
│   │   │   │   │   └── router.go      # تنظیم مسیرهای API
│   │   │   │   │
│   │   │   │   └── grpc/              # (آینده) رابط gRPC 
│   │   │   │
│   │   │   └── persistence/           # آداپتورهای خروجی پایداری داده
│   │   │       ├── mongodb/           # پیاده‌سازی MongoDB
│   │   │       │   ├── models/        # مدل‌های پایگاه داده
│   │   │       │   │   ├── user_model.go
│   │   │       │   │   ├── journey_model.go
│   │   │       │   │   └── booking_model.go
│   │   │       │   │
│   │   │       │   └── repositories/  # پیاده‌سازی مخزن‌ها
│   │   │       │       ├── user_repository.go
│   │   │       │       ├── journey_repository.go
│   │   │       │       └── booking_repository.go
│   │   │       │
│   │   │       └── cache/             # آداپتور کش
│   │   │           ├── redis_repository.go  # پیاده‌سازی Redis
│   │   │           └── memory_repository.go # پیاده‌سازی حافظه
│   │   │
│   │   ├── infrastructure/            # لایه زیرساخت
│   │   │   ├── middleware/            # میدلورهای مشترک
│   │   │   │   ├── auth_middleware.go
│   │   │   │   ├── logging_middleware.go
│   │   │   │   └── error_middleware.go
│   │   │   │
│   │   │   ├── server/                # راه‌اندازی سرور
│   │   │   │   ├── http_server.go
│   │   │   │   └── server.go
│   │   │   │
│   │   │   ├── database/              # مدیریت پایگاه داده
│   │   │   │   ├── mongodb/
│   │   │   │   │   └── connection.go  # اتصال به MongoDB
│   │   │   │   └── cache/
│   │   │   │       └── redis.go       # اتصال به Redis
│   │   │   │
│   │   │   └── config/                # پیکربندی
│   │   │       ├── app.go             # تنظیمات برنامه
│   │   │       ├── database.go        # تنظیمات پایگاه داده
│   │   │       └── server.go          # تنظیمات سرور
│   │   │
│   │   └── pkg/                       # پکیج‌های مشترک
│   │       ├── logger/                # لاگینگ
│   │       │   └── zap_logger.go      # پیاده‌سازی Zap
│   │       │
│   │       ├── validator/             # اعتبارسنجی
│   │       │   └── validator.go
│   │       │
│   │       └── metrics/               # متریک‌ها
│   │           └── prometheus.go      # پیاده‌سازی Prometheus
│   │
│   ├── api/                           # مستندات API
│   │   └── swagger/                   # مستندات Swagger
│   │       └── swagger.yaml
│   │
│   ├── config/                        # فایل‌های تنظیمات
│   │   ├── development.yaml
│   │   ├── test.yaml
│   │   └── production.yaml
│   │
│   ├── tests/                         # تست‌ها
│   │   ├── unit/                      # تست‌های واحد
│   │   │   ├── domain/                # تست‌های لایه دامنه
│   │   │   │   ├── user/
│   │   │   │   ├── journey/
│   │   │   │   └── booking/
│   │   │   │
│   │   │   ├── usecases/              # تست‌های لایه یوزکیس
│   │   │   │   ├── user/
│   │   │   │   ├── journey/
│   │   │   │   └── booking/
│   │   │   │
│   │   │   └── interfaces/            # تست‌های لایه رابط
│   │   │       ├── api/
│   │   │       └── persistence/
│   │   │
│   │   └── integration/               # تست‌های یکپارچگی
│   │       ├── api/                   # تست‌های API
│   │       └── repository/            # تست‌های مخزن
│   │
│   ├── scripts/                       # اسکریپت‌های مفید
│   │   ├── migrations/                # مهاجرت‌های پایگاه داده
│   │   └── setup.sh                   # راه‌اندازی محیط توسعه
│   │
│   └── go.mod                         # وابستگی‌های Go
│
├── frontend/                          # کد فرانت‌اند React
│   ├── public/                        # فایل‌های عمومی
│   │   ├── index.html
│   │   └── favicon.ico
│   │
│   ├── src/                           # کد منبع React
│   │   ├── components/                # کامپوننت‌های قابل استفاده مجدد
│   │   │   ├── common/                # کامپوننت‌های عمومی
│   │   │   │   ├── Button/
│   │   │   │   ├── Input/
│   │   │   │   └── Card/
│   │   │   │
│   │   │   └── layout/                # کامپوننت‌های لایه‌بندی
│   │   │       ├── Header/
│   │   │       └── Footer/
│   │   │
│   │   ├── features/                  # ویژگی‌های اصلی (بر اساس دامنه)
│   │   │   ├── user/                  # ویژگی کاربر
│   │   │   │   ├── components/        # کامپوننت‌های خاص کاربر
│   │   │   │   ├── pages/             # صفحات کاربر
│   │   │   │   └── services/          # سرویس‌های API کاربر
│   │   │   │
│   │   │   ├── journey/               # ویژگی سفر
│   │   │   │   ├── components/
│   │   │   │   ├── pages/
│   │   │   │   └── services/
│   │   │   │
│   │   │   ├── booking/               # ویژگی رزرو
│   │   │   │   ├── components/
│   │   │   │   ├── pages/
│   │   │   │   └── services/
│   │   │   │
│   │   │   └── auth/                  # ویژگی احراز هویت
│   │   │       ├── components/
│   │   │       ├── pages/
│   │   │       └── services/
│   │   │
│   │   ├── hooks/                     # هوک‌های سفارشی
│   │   │   ├── useAuth.js
│   │   │   └── useForm.js
│   │   │
│   │   ├── services/                  # سرویس‌های API مشترک
│   │   │   └── api.js                 # تنظیمات پایه API
│   │   │
│   │   ├── utils/                     # توابع کمکی
│   │   │   ├── helpers.js
│   │   │   └── formatter.js
│   │   │
│   │   ├── App.js                     # کامپوننت اصلی
│   │   └── index.js                   # نقطه ورودی
│   │
│   ├── package.json                   # وابستگی‌های npm
│   └── .env                           # متغیرهای محیطی
│
├── docker-compose.yml                 # تنظیمات docker-compose
├── Dockerfile.backend                 # Dockerfile برای بک‌اند
├── Dockerfile.frontend                # Dockerfile برای فرانت‌اند
│
├── docs/                              # مستندات پروژه
│   ├── architecture.md                # معماری
│   ├── architecture.png               # نمودار معماری
│   ├── development.md                 # راهنمای توسعه
│   └── api.md                         # توضیحات API
│
└── README.md                          # مستندات اصلی پروژه با نمودار معماری سطح بالا
```

## مزایای معماری هگزاگونال

1. **جداسازی دغدغه‌ها**: دامنه کسب‌وکار از جزئیات فنی جدا شده است.
2. **تغییرپذیری**: می‌توانیم بدون تغییر در هسته دامنه، آداپتورها را جایگزین کنیم.
3. **تست‌پذیری**: لایه‌ها به راحتی با mock‌ها قابل تست هستند.
4. **توسعه‌پذیری**: تیم‌های مختلف می‌توانند همزمان روی لایه‌های مختلف کار کنند.
5. **قابلیت نگهداری**: کد منظم‌تر و قابل نگهداری‌تر است.
6. **مقیاس‌پذیری**: امکان مقیاس‌پذیری مستقل بخش‌های مختلف سیستم
7. **افزایش بهره‌وری**: کاهش وابستگی‌ها منجر به افزایش بهره‌وری تیم‌های توسعه می‌شود
8. **کاهش ریسک فنی**: آزادی در تغییر فناوری‌ها بدون تأثیر بر منطق کسب‌وکار

## لایه‌های معماری

### 1. هسته دامنه (Domain Core)
- **موجودیت‌ها (Entities)**: مدل‌های اصلی کسب‌وکار (User، Journey، Booking)
  - حاوی داده و رفتار (حالت و منطق)
  - قوانین کسب‌وکار مربوط به هر موجودیت
  - بدون وابستگی به لایه‌های دیگر

- **سرویس‌های دامنه (Domain Services)**: منطق کسب‌وکار مشترک بین موجودیت‌ها
  - منطق عملیاتی که به چندین موجودیت مربوط می‌شود
  - عملیات پیچیده دامنه
  - هماهنگی بین موجودیت‌ها

- **اشیاء ارزش (Value Objects)**: اشیاء بدون هویت (Email، Rating، Status)
  - مقادیری که با ویژگی‌هایشان تعریف می‌شوند، نه با هویت
  - غیرقابل تغییر (Immutable)
  - کپسوله‌سازی قوانین اعتبارسنجی

- **قوانین و سیاست‌ها (Rules & Policies)**: محدودیت‌ها و قوانین کسب‌وکار
  - قوانین مرتبط با چندین موجودیت یا سرویس
  - اعمال محدودیت‌های کسب‌وکار
  - تصمیم‌گیری بر اساس قوانین

### 2. پورت‌ها (Ports)
- **پورت‌های اولیه (Primary)**: رابط‌های استفاده از برنامه
  - قراردادهای تعامل با هسته برنامه
  - مشخص کردن عملیات قابل انجام
  - مستقل از جزئیات پیاده‌سازی

- **پورت‌های ثانویه (Secondary)**: رابط‌های ارتباط با خارج
  - قراردادهایی برای تعامل با منابع خارجی
  - مشخص کردن نیازهای هسته از دنیای بیرون
  - جداسازی منطق دامنه از زیرساخت

### 3. آداپتورها (Adapters)
- **آداپتورهای اولیه (Primary)**: هندلرهای REST API، ریزالورهای GraphQL
  - مترجم بین درخواست‌های خارجی و هسته برنامه
  - تبدیل داده‌های ورودی به فرمت مناسب برای هسته
  - مدیریت جریان ورودی به سیستم

- **آداپتورهای ثانویه (Secondary)**: پیاده‌سازی‌های مخزن (MongoDB، Redis)
  - پیاده‌سازی پورت‌های ثانویه با فناوری‌های خاص
  - ارتباط با پایگاه داده، API‌های خارجی، سیستم فایل، و غیره
  - مبدل بین نیازهای دامنه و سیستم‌های خارجی

### 4. دغدغه‌های مشترک (Cross-Cutting Concerns)
- **لاگینگ (Logging)**: ثبت رویدادها و خطاها
- **احراز هویت (Authentication)**: تأیید هویت کاربران
- **مجوزدهی (Authorization)**: کنترل دسترسی
- **پیکربندی (Configuration)**: مدیریت تنظیمات
- **مدیریت خطا (Error Handling)**: پردازش و گزارش خطاها
- **مانیتورینگ (Monitoring)**: نظارت بر عملکرد سیستم
- **کش‌گذاری (Caching)**: ذخیره موقت داده‌ها برای بهبود کارایی
- **مدیریت تراکنش (Transaction Management)**: حفظ یکپارچگی داده‌ها

## جریان داده و کنترل

1. **جریان ورودی**:
   - درخواست کاربر (Frontend/Mobile) → Primary Adapter → Primary Port → Use Case → Domain Core

2. **جریان خروجی**:
   - Domain Core → Secondary Port → Secondary Adapter → External System (DB/API)

3. **قانون وابستگی**:
   - هسته دامنه به هیچ لایه دیگری وابسته نیست
   - پورت‌ها به هسته دامنه وابسته‌اند
   - آداپتورها به پورت‌ها وابسته‌اند
   - جریان وابستگی همیشه به سمت داخل است (هسته دامنه)

## شروع کار با پروژه

### پیش‌نیازها
- Go 1.18+
- Docker و Docker Compose
- MongoDB
- Redis
- Node.js 16+ و npm (برای فرانت‌اند)
- Git

### راه‌اندازی محیط توسعه

#### 1. کلون کردن مخزن
```bash
git clone https://github.com/amirhossein-jamali/journey-hub.git
cd journey-hub
```

#### 2. راه‌اندازی دیتابیس‌ها با Docker
```bash
cd backend/deployments/docker
docker-compose up -d mongodb redis
```

#### 3. تنظیم متغیرهای محیطی
```bash
cp backend/.env.example backend/.env
# ویرایش فایل .env با متغیرهای مورد نیاز
```

#### 4. اجرای سرور بک‌اند
```bash
cd backend
go mod download
go run cmd/server/main.go
```

#### 5. راه‌اندازی فرانت‌اند
```bash
cd frontend
npm install
npm start
```

### راه‌اندازی با Docker Compose (کل سیستم)
```bash
docker-compose up -d
```

### دسترسی به سرویس‌ها
- سرور بک‌اند: `http://localhost:8080/api/v1`
- مستندات API: `http://localhost:8080/api/v1/docs`
- فرانت‌اند: `http://localhost:3000`
- مانیتورینگ Prometheus: `http://localhost:9090`

## مشارکت در پروژه

برای مشارکت در این پروژه، لطفا فایل `CONTRIBUTING.md` را مطالعه کنید. ما از هر نوع مشارکتی استقبال می‌کنیم:

- گزارش باگ
- پیشنهاد ویژگی‌های جدید
- ارسال Pull Request
- بهبود مستندات
- اشتراک‌گذاری بازخورد

## مستندات
برای اطلاعات بیشتر به پوشه `docs/` مراجعه کنید. مستندات شامل:
- [معماری تفصیلی](docs/architecture.md)
- [راهنمای توسعه](docs/development.md)
- [مستندات API](docs/api.md) و [Swagger API](api/swagger/swagger.yaml)
- [راهنمای نصب و راه‌اندازی](docs/installation.md)
- [تصمیمات معماری](docs/architecture-decisions.md)

## لایسنس

این پروژه تحت لایسنس MIT منتشر شده است. برای جزئیات بیشتر به فایل [LICENSE](LICENSE) مراجعه کنید. 