# Go Authentication Service

A lightweight Go authentication service that validates admin credentials and issues JWT tokens.

This service is framework-minimal (`net/http`) and designed to be easy to understand, test, and integrate with multiple backend apps (for example Spring Boot APIs).

## Features

- `POST /login` endpoint for credential validation
- JWT token generation (`HS256`)
- Environment-based configuration (`.env`)
- Health endpoint (`GET /health`)
- Small and clean project structure for learning and extension

## Requirements
- Go 1.25.0 (see go.mod)
- Git (optional)

## Run Locally
```
go mod tidy
go run main.go
```

Expected startup log: `Starting auth service on : 8081`

## API Endpoints
### Health Check
- Method: GET
- URL: /health
- Response: `Auth service is running`

### Login
- Method: POST
- URL: /login
```
{
  "username": "admin",
  "password": "change-this-password"
}
```

## Roadmap ideas

- Add persistent user store (PostgreSQL/SQLite)
- Add roles/permissions claims
- Add refresh token flow
- Add key rotation and JWKS
- Add middleware-protected routes for internal verification
- Add unit/integration tests

