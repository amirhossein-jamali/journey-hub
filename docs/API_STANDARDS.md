# API Standards in Journey Hub

This document outlines the design and documentation standards for APIs in the Journey Hub project.

## REST API Design Principles

### 1. Using Resources and HTTP Verbs

- **Resources**: Resources should be identified with appropriate plural names (e.g., `/users`, `/journeys`)
- **HTTP Verbs**: Use standard HTTP verbs for CRUD operations:
  - `GET`: Retrieve one or more resources
  - `POST`: Create a new resource
  - `PUT`: Complete update of an existing resource
  - `PATCH`: Partial update of an existing resource
  - `DELETE`: Delete a resource

### 2. Endpoint Naming

- Use meaningful REST names
- Use kebab-case for multi-part names
- Use API versioning (e.g., `/api/v1/users`)

#### Suggested Examples:

```
GET    /api/v1/users              # Get list of users
POST   /api/v1/users              # Create a new user
GET    /api/v1/users/{id}         # Get a specific user by ID
PUT    /api/v1/users/{id}         # Full update of a user
PATCH  /api/v1/users/{id}         # Partial update of a user
DELETE /api/v1/users/{id}         # Delete a user

GET    /api/v1/users/{id}/journeys # Get a user's journeys
```

### 3. Query Parameters

- Use query parameters for filtering, sorting, and pagination
- Parameter names should be lowercase and camelCase

```
GET /api/v1/journeys?destination=tehran     # Filter by destination
GET /api/v1/journeys?minPrice=100&maxPrice=500   # Filter by price range
GET /api/v1/journeys?sortBy=price&sortDir=asc   # Sorting
GET /api/v1/journeys?page=2&pageSize=10    # Pagination
```

### 4. HTTP Status Codes

Use standard HTTP status codes:

- **2xx**: Success
  - `200 OK`: Successful request (for GET, PUT, PATCH)
  - `201 Created`: New resource created (for POST)
  - `204 No Content`: Successful request with no content (for DELETE)
- **4xx**: Client Error
  - `400 Bad Request`: Invalid request
  - `401 Unauthorized`: Authentication required
  - `403 Forbidden`: No access
  - `404 Not Found`: Resource not found
  - `409 Conflict`: Conflict with current state of resource (e.g., duplicate email)
  - `422 Unprocessable Entity`: Validation failed
- **5xx**: Server Error
  - `500 Internal Server Error`: Server error
  - `503 Service Unavailable`: Service unavailable

### 5. Response Format

- Use JSON as the main response format
- Use consistent structure for responses
- Always return an appropriate HTTP status code

#### Success Response:

```json
{
  "data": {
    "id": "5f8d0c1b8f67b1a4f2a3d9e7",
    "name": "Mahdi Ahmadi",
    "email": "mahdi@example.com"
  },
  "meta": {
    "timestamp": "2023-10-18T10:15:30Z"
  }
}
```

#### List Response:

```json
{
  "data": [
    {
      "id": "5f8d0c1b8f67b1a4f2a3d9e7",
      "name": "Mahdi Ahmadi",
      "email": "mahdi@example.com"
    },
    // ...
  ],
  "meta": {
    "page": 1,
    "pageSize": 10,
    "totalItems": 42,
    "totalPages": 5
  }
}
```

#### Error Response:

```json
{
  "error": {
    "code": "USER_NOT_FOUND",
    "message": "The requested user was not found",
    "details": "User with ID 5f8d0c1b8f67b1a4f2a3d9e7 does not exist in the system"
  },
  "meta": {
    "timestamp": "2023-10-18T10:15:30Z"
  }
}
```

## API Documentation with Swagger

### 1. Swagger Setup

Using [swaggo/swag](https://github.com/swaggo/swag) to generate Swagger documentation from Go comments:

```go
// @title Journey Hub API
// @version 1.0
// @description Journey Hub API
// @termsOfService http://journeyhub.ir/terms/
// @contact.name API Support
// @contact.email support@journeyhub.ir
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host api.journeyhub.ir
// @BasePath /api/v1
```

### 2. Handler Documentation

Each handler should have complete Swagger documentation:

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
// @Failure 409 {object} dto.ErrorResponse
// @Failure 422 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    // ...
}
```

### 3. Swagger Model Definitions

DTO models should be defined with appropriate Swagger tags:

```go
// CreateUserRequest model for creating a user
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required" example:"Ali Mohammadi" minLength:"3" maxLength:"50"`
    Email    string `json:"email" binding:"required,email" example:"ali@example.com" format:"email"`
    Password string `json:"password" binding:"required" example:"StrongP@ss123" minLength:"8" format:"password"`
}

// UserResponse model for user information
type UserResponse struct {
    ID        string    `json:"id" example:"5f8d0c1b8f67b1a4f2a3d9e7"`
    Name      string    `json:"name" example:"Ali Mohammadi"`
    Email     string    `json:"email" example:"ali@example.com"`
    CreatedAt time.Time `json:"createdAt" example:"2023-10-18T10:15:30Z"`
    UpdatedAt time.Time `json:"updatedAt" example:"2023-10-18T10:15:30Z"`
}
```

### 4. Tag Organization

Use Swagger tags to logically group endpoints:

```go
// @title Journey Hub API
// @version 1.0
// ...
// @tag.name users
// @tag.description User operations
// @tag.name journeys
// @tag.description Journey operations
// @tag.name bookings
// @tag.description Booking operations
// @tag.name auth
// @tag.description Authentication operations
```

### 5. Path and Query Parameter Documentation

URL and query parameters should be well documented:

```go
// GetJourneys godoc
// @Summary Get list of journeys
// @Description Get list of journeys with filtering, sorting and pagination
// @Tags journeys
// @Accept json
// @Produce json
// @Param destination query string false "Filter by destination"
// @Param minPrice query number false "Minimum price"
// @Param maxPrice query number false "Maximum price"
// @Param startDate query string false "Start date (YYYY-MM-DD)" format(date)
// @Param endDate query string false "End date (YYYY-MM-DD)" format(date)
// @Param page query int false "Page number (starts from 1)" minimum(1) default(1)
// @Param pageSize query int false "Items per page" minimum(1) maximum(100) default(10)
// @Success 200 {object} dto.PaginatedResponse{data=[]dto.JourneySummaryResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /journeys [get]
```

### 6. Authentication Documentation

For endpoints requiring authentication, use `@Security`:

```go
// @SecurityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token with Bearer {token} format
```

And then in each handler that requires authentication:

```go
// CreateJourney godoc
// @Summary Create new journey
// @Description Create a new journey by a user
// @Tags journeys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param journey body dto.CreateJourneyRequest true "Journey information"
// @Success 201 {object} dto.JourneyResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 422 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /journeys [post]
```

## API Development Best Practices

### 1. Input Validation

- Validate all input data on the server
- Use `validator` for structural validation
- Use security policies such as length restrictions, allowed patterns, etc.

### 2. Error Handling

- Provide friendly and useful error messages
- Don't show internal error details to users
- Define standard error codes

### 3. Access Control

- Use middleware for access control
- Define access levels based on user roles
- Document access requirements for each endpoint

### 4. Pagination and Rate Limiting

- Implement pagination for endpoints that return large lists
- Implement rate limiting to prevent DoS attacks
- Apply limits based on IP or user

### 5. CORS Standard

- Use CORS to restrict access to resources
- Only define allowed domains in CORS
- Use different CORS settings for development environment

## API Versioning

We use URL-based versioning:

```
/api/v1/users
/api/v2/users
```

### Versioning Rules:

1. Never make breaking changes to a published API
2. For major (breaking) changes, create a new version
3. You can add new features to older versions, as long as backward compatibility is maintained
4. Document and announce deprecation schedule for older versions

## API FAQs

### 1. Can we use GraphQL?

Currently, we use REST API. Using GraphQL in the future will be evaluated based on project needs.

### 2. Should all responses follow a standard format?

Yes, using the standard response pattern (`data`/`error` and `meta`) is required for all APIs.

### 3. How do we ensure API security?

- Always use HTTPS
- Use JWT for authentication
- Use CSRF tokens for forms
- Implement rate limiting
- Validate all inputs
- Use prepared statements for database queries 