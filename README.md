# Simple Task Management API

A RESTful API for task and project management built with Go, Gin, GORM, and PostgreSQL. This application provides user authentication, project management, task tracking, and audit logging functionality.

## Features

- 🔐 **Authentication & Authorization**: JWT-based authentication with role-based access control (Admin, Manager, Employee)
- 📋 **Project Management**: Create, read, update, and delete projects
- ✅ **Task Management**: Comprehensive task CRUD operations with status tracking
- 👥 **User Management**: User registration, profile management, and role assignment
- 📊 **Audit Logging**: Track all system activities for compliance and monitoring
- 🔄 **Background Jobs**: Async task processing with Redis and Asynq
- 📧 **Email Notifications**: Automated email notifications for important events
- 🐳 **Docker Support**: Containerized deployment with Docker and Docker Compose
- 📖 **API Documentation**: Auto-generated Swagger documentation

## Tech Stack

- **Backend**: Go 1.24, Gin Web Framework
- **Database**: PostgreSQL with GORM ORM
- **Cache/Queue**: Redis with Asynq for background jobs
- **Authentication**: JWT tokens
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker & Docker Compose
- **Migration**: golang-migrate

## Quick Start

### Prerequisites

- Go 1.24+
- PostgreSQL
- Redis
- Docker (optional)

### Environment Setup

1. Clone the repository:

```bash
git clone https://github.com/mnizarzr/dot-test.git
cd dot-test
```

2. Copy environment file and configure:

```bash
cp .env.example .env
# Edit .env with your database and Redis configurations
```

### Running with Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f
```

### Running Locally

1. Install dependencies:

```bash
go mod download
```

2. Run database migrations:

```bash
make migrate
```

3. Create an admin user:

```bash
go run main.go create-admin
```

4. Seed the database (optional):

```bash
make seed
```

5. Start the server:

```bash
make serve
# or
go run main.go serve
```

6. Start the queue worker (in another terminal):

```bash
go run main.go start-queue
```

The API will be available at `http://localhost:8080`

## API Documentation

Once the server is running, you can access the Swagger documentation at:

- **Swagger UI**: `http://localhost:8080/swagger/index.html`

## Available Commands

The application provides several CLI commands:

```bash
# Start the HTTP server
go run main.go serve

# Start the background queue worker
go run main.go start-queue

# Create an admin user
go run main.go create-admin

# Seed the database with sample data
go run main.go seed
```

## Makefile Commands

```bash
# Development
make serve          # Start the server
make seed           # Seed the database

# Database migrations
make migration name=create_table_name  # Create new migration
make migrate                          # Run pending migrations
make rollback                        # Rollback last migration
make rollback-all                    # Rollback all migrations
make force-migrate version=N         # Force migration to specific version
```

## Development

### Project Structure

```
.
├── app/                 # Application setup and routing
├── cmd/                 # CLI commands
├── common/              # Shared utilities and responses
├── config/              # Configuration management
├── db/                  # Database connection and migrations
├── docs/                # Swagger documentation
├── entity/              # Database entities/models
├── jobs/                # Background job definitions
├── middleware/          # HTTP middleware
├── modules/             # Feature modules (auth, user, project, task)
├── template/            # Email templates
└── utils/               # Utility functions
```

## License

This project is licensed under the Unlicense - see the [LICENSE](LICENSE) file for details.
