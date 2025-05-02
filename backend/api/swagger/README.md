# Swagger Documentation for Journey Hub

## Overview
This directory contains Swagger/OpenAPI documentation for the Journey Hub API. The documentation is generated automatically from code comments.

## Setup
To generate Swagger documentation, we use [swag](https://github.com/swaggo/swag) tool.

### Installation
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Generating Documentation
Run the following command from the root of the backend directory:
```bash
swag init -g cmd/app/main.go -o api/swagger
```

## Documentation Format
When writing API handlers, use the following comment format to document your endpoints:

```go
// @Summary Create a new user
// @Description Register a new user in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration details"
// @Success 201 {object} dto.TokenResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
    // Implementation
}
```

## Component Documentation

### Data Transfer Objects (DTOs)
Document your DTOs using the following format:

```go
// RegisterRequest represents the request body for user registration.
// @Description User registration request
type RegisterRequest struct {
	// User's full name
	Name string `json:"name" binding:"required" example:"John Doe"`
	
	// User's email address
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
	
	// User's password (min 8 characters)
	Password string `json:"password" binding:"required,min=8" example:"secureP@ss123"`
}
```

## Viewing Documentation
After generating the documentation, you can view it using Swagger UI:

1. Start the application
2. Navigate to `/swagger/index.html` in your browser

## Integration with Frontend
The frontend application can use this Swagger documentation to generate API clients automatically. 