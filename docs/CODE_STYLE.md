# Journey Hub Code Style Guide

This document provides more details about coding standards in the Journey Hub project.

## Go Language

### Package Structure

- Each package should have a single responsibility
- Package names should be short, meaningful, and lowercase
- Don't use under_score in package names
- Each package should have `package` documentation in the main file

### Imports

- Use import grouping:
  1. Standard packages
  2. External packages
  3. Internal project packages
- Avoid using `.` for imports
- Avoid unused imports

```go
import (
    "context"
    "time"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    "github.com/amirhossein-jamali/journey-hub/pkg/logger"
)
```

### Error Handling

- Never ignore errors
- Use predefined error variables for common errors
- When returning an error, use `fmt.Errorf` with the `%w` verb to preserve the error chain
- Use `errors.Is` and `errors.As` to check and convert errors
- Use custom error structures with meta fields

```go
// Define predefined errors
var (
    ErrNotFound = errors.New("resource not found")
    ErrInvalidInput = errors.New("invalid input")
)

// Return error with chain preservation
if err != nil {
    return fmt.Errorf("failed to fetch user: %w", err)
}

// Check error type
if errors.Is(err, ErrNotFound) {
    // Take appropriate action
}
```

### Code Documentation

- Every exported package, type, function, and method must be documented
- Documentation should be accurate, concise, and helpful
- Avoid just repeating names in documentation; explain what is being done

```go
// UserService provides operations for managing users
type UserService interface {
    // FindByID retrieves a user by their ID
    // Returns ErrNotFound if the user doesn't exist
    FindByID(ctx context.Context, id string) (*User, error)
    
    // Create adds a new user to the system
    // Returns ErrInvalidInput if the user data is invalid
    Create(ctx context.Context, user *User) error
}
```

### Testing

#### Test Structure

- Use the AAA (Arrange-Act-Assert) pattern:
  1. Arrange data and conditions
  2. Execute the code under test
  3. Assert results
- Use descriptive and meaningful names for tests

```go
func TestUserService_FindByID_ExistingUser(t *testing.T) {
    // Arrange
    userRepo := mockUserRepository{...}
    service := NewUserService(userRepo)
    expected := &User{ID: "123", Name: "Test User"}
    
    // Act
    user, err := service.FindByID(context.Background(), "123")
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expected, user)
}
```

#### Test Types

- **Unit Tests**: Test small units of code independently of other parts
- **Integration Tests**: Test collaboration between multiple parts
- **Performance Tests**: Test execution time and resource consumption with benchmarks

#### Mock Testing

- For unit tests, mock dependencies
- Use `gomock` or manual mock interfaces
- When using mocks, only define behavior needed for the test

## Hexagonal Architecture

### Architecture Layers

#### 1. Domain (Core)

- Contains entities and core business rules
- Independent of any technology or external framework
- All entities should be encapsulated and modified through methods

```go
// domain/user/entity.go
package user

import "time"

type User struct {
    id        string
    email     string
    name      string
    createdAt time.Time
    updatedAt time.Time
}

func NewUser(email, name string) *User {
    return &User{
        id:        generateID(),
        email:     email,
        name:      name,
        createdAt: time.Now(),
        updatedAt: time.Now(),
    }
}

func (u *User) ID() string {
    return u.id
}

func (u *User) Email() string {
    return u.email
}

func (u *User) UpdateName(name string) {
    u.name = name
    u.updatedAt = time.Now()
}
```

#### 2. Usecases (Application)

- Application logic
- Uses domain entities
- Defines ports (interfaces) for communicating with the outside world

```go
// usecases/user/service.go
package user

import (
    "context"
    
    "github.com/amirhossein-jamali/journey-hub/internal/domain/user"
)

type Repository interface {
    FindByID(ctx context.Context, id string) (*user.User, error)
    Save(ctx context.Context, user *user.User) error
}

type Service struct {
    repo Repository
}

func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) GetUser(ctx context.Context, id string) (*user.User, error) {
    return s.repo.FindByID(ctx, id)
}

func (s *Service) UpdateUserName(ctx context.Context, id, name string) error {
    u, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return err
    }
    
    u.UpdateName(name)
    return s.repo.Save(ctx, u)
}
```

#### 3. Infrastructure

- Implementation of secondary ports (output adapters)
- Communication with external systems like databases, web services, etc.

```go
// infrastructure/mongodb/user_repository.go
package mongodb

import (
    "context"
    
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    
    "github.com/amirhossein-jamali/journey-hub/internal/domain/user"
)

type UserRepository struct {
    collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
    return &UserRepository{
        collection: db.Collection("users"),
    }
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
    // Implementation for MongoDB search
}

func (r *UserRepository) Save(ctx context.Context, user *user.User) error {
    // Implementation for MongoDB save
}
```

#### 4. Interfaces

- Implementation of primary ports (input adapters)
- Such as HTTP controllers, gRPC clients, etc.

```go
// interfaces/api/http/handlers/user_handler.go
package handlers

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    
    "github.com/amirhossein-jamali/journey-hub/internal/usecases/user"
)

type UserHandler struct {
    userService *user.Service
}

func NewUserHandler(userService *user.Service) *UserHandler {
    return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUser(c *gin.Context) {
    id := c.Param("id")
    
    user, err := h.userService.GetUser(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, user)
}
```

## DTO (Data Transfer Objects)

- Use DTOs for data transfer between layers
- DTO structures should use appropriate tags for validation and binding
- Never use domain entities directly in the API

```go
// interfaces/api/http/dto/user_dto.go
package dto

type CreateUserRequest struct {
    Email string `json:"email" binding:"required,email"`
    Name  string `json:"name" binding:"required"`
}

type UserResponse struct {
    ID    string `json:"id"`
    Email string `json:"email"`
    Name  string `json:"name"`
}

func MapToUserResponse(user *user.User) UserResponse {
    return UserResponse{
        ID:    user.ID(),
        Email: user.Email(),
        Name:  user.Name(),
    }
}
```

## REST API Error Handling

- Use appropriate HTTP status codes
- Standardize error response structure
- Don't display technical error details in production

```go
// Standard error response
type ErrorResponse struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

// Example usage
if err != nil {
    if errors.Is(err, ErrNotFound) {
        c.JSON(http.StatusNotFound, ErrorResponse{
            Code:    "NOT_FOUND",
            Message: "The requested resource was not found",
        })
        return
    }
    
    c.JSON(http.StatusInternalServerError, ErrorResponse{
        Code:    "INTERNAL_ERROR",
        Message: "An unexpected error occurred",
    })
    return
}
```

## Other Considerations

For points not covered in this document, refer to [Effective Go](https://golang.org/doc/effective_go.html) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments). 