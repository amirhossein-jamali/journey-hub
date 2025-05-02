# Contributing to Journey Hub

This document provides guidance for contributors to the Journey Hub project. All contributors are expected to read and follow this guide.

## Git Workflow

We use the [GitHub Flow](https://guides.github.com/introduction/flow/) model for development:

1. Create a new branch from `develop` for each feature or fix
2. Commit your code
3. Submit a Pull Request to the `develop` branch
4. After code review and approval, your branch will be merged into `develop`
5. The `main` branch is used for stable, releasable versions

### Branch Naming Conventions

- For new features: `feature/name-of-feature`
- For bug fixes: `bugfix/issue-number-or-description`
- For technical improvements: `tech/short-description`
- For documentation changes: `docs/what-is-changed`

### Commit Message Standards

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Main `type` values:
- `feat`: a new feature
- `fix`: a bug fix
- `docs`: documentation changes
- `style`: changes that do not affect code meaning (whitespace, formatting)
- `refactor`: code changes that neither fix a bug nor add a feature
- `test`: adding or modifying tests
- `chore`: updating build tasks, package manager configs, etc.

Example:
```
feat(auth): add JWT authentication

Implement JWT token generation and validation for authentication.
```

## Code Standards

### Go Language

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Go code must pass `gofmt`
- All code must pass `golangci-lint` before submission
- All public functions, methods, and structs must be documented
- Use [idiomatic error handling](https://github.com/golang/go/wiki/Errors)

### Naming

- Use meaningful and descriptive names
- Use camelCase for variables and fields
- Use PascalCase for struct names, interfaces, and exported functions
- Use common abbreviations like `ID` instead of `Id`
- Don't use `I` prefix for interfaces

### Architecture

We use Hexagonal Architecture (or Ports and Adapters):

- **Domain**: Business entities and core business rules
- **Usecases**: Application logic and services
- **Interfaces**: Input adapters like APIs and CLI
- **Infrastructure**: Output adapters like databases and external services
- **Pkg**: Shared libraries and reusable components

## API Documentation

We use Swagger/OpenAPI for API documentation:

- All API endpoints must have Swagger comments
- Follow [swaggo/swag](https://github.com/swaggo/swag) guidelines
- All parameters, responses, and error codes must be documented

Example:
```go
// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User information"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users [post]
```

## Testing

- Write tests for all operational code
- Create a `_test.go` file for each `.go` file
- Use [testify](https://github.com/stretchr/testify) for better assertions
- Prefix integration tests with `Test_Integration_`
- Benchmark tests should start with `Benchmark`

## Linting and Formatting

- Go code must be formatted with `gofmt` before each commit
- Use `golangci-lint` to check code quality
- The `golangci-lint` configuration is in the `.golangci.yml` file

## Development Environment

- Use VS Code or GoLand as the recommended IDEs
- Use Dev Container for development environment consistency
- Place appropriate configuration files for your development environment in the `config/development` folder before running the application

## Questions and Help

If you have questions or need help, please use GitHub Issues or team communication channels. 