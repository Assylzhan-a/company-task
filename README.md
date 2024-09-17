# Company Management Service

## Description

This project is a Company Management Service built with Go, providing a RESTful API for managing company information, user authentication, and featuring event sourcing with Kafka integration.

### Key Features

- CRUD operations for company management
- User registration and authentication with JWT
- Event sourcing using outbox pattern
- Kafka integration for event publishing
- PostgreSQL database for data persistence
- Docker support for easy deployment and development

## Prerequisites

- **Docker** is required to run the application and its dependencies.
- **Make** is recommended for running commands easily, but not strictly required.

## Quick Start

1. Clone the repository:
   ```sh
   git clone https://github.com/Assylzhan-a/company-task.git
   cd company-task
   ```

2. Run the application:
   ```sh
   make run
   ```
   or, if you don't have Make:
   ```sh
   docker-compose up
   ```

The application and all its dependencies (PostgreSQL, Kafka, etc.) will start up using Docker Compose.

## API Endpoints

Here are the main API endpoints with curl examples:

### User Registration

```sh
curl -X POST http://localhost:8080/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{"username": "newuser", "password": "password123"}'
```

### User Login

```sh
curl -X POST http://localhost:8080/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"username": "newuser", "password": "password123"}'
```

This will return a JWT token to use for authenticated requests.

### Create Company

```sh
curl -X POST http://localhost:8080/v1/companies \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "Tech Corp",
    "description": "Innovative technology solutions",
    "amount_of_employees": 100,
    "registered": true,
    "type": "Corporations"
  }'
```

> **Note**: The `id` field should be a valid UUID generated by the client.

### Get Company

```sh
curl -X GET http://localhost:8080/v1/companies/123e4567-e89b-12d3-a456-426614174000 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Update Company

```sh
curl -X PATCH http://localhost:8080/v1/companies/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Updated Tech Corp",
    "amount_of_employees": 150
  }'
```

### Delete Company

```sh
curl -X DELETE http://localhost:8080/v1/companies/123e4567-e89b-12d3-a456-426614174000 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Additional Features and Commands

- **Kafka UI**: View Kafka messages at http://localhost:8090
- Run tests: `make test`
- Run linter: `make lint`
- Run database migrations: `make migrate-up`
- Rollback database migrations: `make migrate-down`
- Postman collection is provided - CompanyService.postman_collection.json

## Configuration

To modify the configuration, edit the `app.env` file in the root directory.

## Project Structure

```
.
├── cmd
│   └── api
│       └── main.go
├── config
│   └── config.go
├── internal
│   ├── auth
│   ├── db
│   ├── delivery
│   ├── domain
│   ├── kafka
│   ├── ports
│   └── worker
├── migrations
├── pkg
│   ├── errors
│   └── logger
├── tests
├── Dockerfile
├── docker-compose.yml
├── app.env
├── go.mod
├── go.sum
└── Makefile
```

This structure follows clean architecture principles, separating concerns into different layers for improved maintainability and testability.
