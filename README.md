# rdev-auth-go

A standardized, production-ready Authentication and Authorization boilerplate for Go services. 

## 🔐 Features
*   **JWT Integration:** Secure token generation and validation.
*   **Hierarchical RBAC:** Role-based access control with support for inheritance.
*   **Middleware:** Plug-and-play Gin middleware for protecting routes.
*   **SQL-First:** Optimized for Bun ORM and PostgreSQL.

## 🛠 Usage
This repository is designed to be the security foundation for all `rdev` projects.

## Folder Structure
```
.
├── cmd/
│   ├── api/                # Main entry point; initializes the server, services, and database
├── internal/
│   ├── config/             # Configuration management; loads environment variables
│   ├── data/               # The "Source of Truth"; contains database models and repositories
│   │   └── migrations/     # SQL files defining schema changes over time
│   ├── dto/                # Data Transfer Objects; handles request/response shapes and query sanitization
│   ├── server/             # HTTP layer; contains route definitions and controller handlers
│   ├── service/            # Business logic layer; bridges DTOs and Data models
├── go.mod                  # Project dependencies
└── go.sum                  # Dependency checksums
```