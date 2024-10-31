# Debt Solver - Authentication Microservice

This repository contains the authentication microservice for the **Debt Solver** project, an expense management mobile application that allows users to track and manage their finances.

## Key Features

- **User Registration**: Securely create and store user accounts.
- **User Login**: Validate credentials and issue JWT tokens for session management.
- **Password Management**: Hash and securely manage user passwords.
- **Authorization**: Protect routes with JWT authentication.

## Technologies Used

- **Golang & Gin**: For building the service.
- **PostgreSQL**: For storing user data.
- **GORM**: For ORM database interactions.
- **JWT**: For user authentication and authorization.
- **Bcrypt**: For password hashing.
- **Viper**: For configuration management.

## Directory Structure

```plaintext
auth-service/
│
├── cmd/
│   └── auth-service/
│       └── main.go                  # Entry point for the application
│
├── configs/
│   └── config.yaml                  # Configuration file for the service
│
├── db/
│   └── migrate.go                   # Database migrations
│
├── internal/
│   ├── controller/
│   │   └── auth_controller.go       # Controller for authentication handlers
│   ├── middleware/
│   │   └── auth_middleware.go       # Middleware for JWT authentication
│   ├── model/
│   │   └── user.go                  # User model and database interactions
│   └── routes/
│       └── routes.go                # Define routes for authentication endpoints
│
├── utils/
│   └── response.go                  # Utility functions for handling responses
│
├── Dockerfile                       # Dockerfile for building the container
├── go.mod                           # Go module file
└── README.md                        # Project documentation
```

## Setup and Installation

<code>git clone https://github.com/debt-solver/DB-auth.git <br>
cd debt-solver-auth
</code>

## Setup PostgreSQL

<code>
docker run --name debt-solver-postgres -e POSTGRES_PASSWORD=yourpassword -d -p 5432:5432 postgres
</code>

## Install Dependencies

<code>go mod tidy</code>

## SMTP Testing

<p>Mailtrap.io Server Will be used for testing purpose</p>

<code>{{d.backend.developer.personal.email.will.receive.all.code}}</code>

## Run Database Migrations

Setup the database schema using the migration file

## Run the application

go run cmd/auth-service/main.go

## Build and Run with Docker

<code>
  docker build -t auth-service . <br>
  docker run -p 8080:8080 auth-service
</code>

## API Endpoints

POST /signup: Register a new user
POST /login: Authenticate and receive a JWT token

<!-- GET /profile: Retrieve the authenticated user's profile(protected by JWT) -->

POST /reset-password
POST /logout
POST /verify-email

## Environment Varibles

DB_HOST=localhost
DB_PORT=5432
DB_USER=auth_user
DB_PASSWORD=yourpassword
JWT_SECRET={{SuperSecret}}

## License

This project is open-source and licensed under the MIT License.

## Contributions

Contributions are welcome! Feel free to open an issue or submit a pull request.
