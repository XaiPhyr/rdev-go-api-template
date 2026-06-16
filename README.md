# rdev-go-api-template

A lightweight, high-performance production-ready REST API boilerplate built with Go, utilizing the Gin Web Framework and Bun ORM. Optimized for speed, developer ergonomics, and clear separation of concerns in a monolithic architecture.

[![Go Report Card](https://goreportcard.com/badge/github.com/XaiPhyr/rdev-go-api-template)](https://goreportcard.com/report/github.com/XaiPhyr/rdev-go-api-template)
[![Build Status](https://img.shields.io/github/actions/workflow/status/XaiPhyr/rdev-go-api-template/go-test.yml?label=test&logo=github&logoColor=white&style=flat)](https://github.com/XaiPhyr/rdev-go-api-template/actions)
[![Build Status](https://img.shields.io/github/actions/workflow/status/XaiPhyr/rdev-go-api-template/go-sec.yml?label=security&logo=github&logoColor=white&style=flat)](https://github.com/XaiPhyr/rdev-go-api-template/actions)
[![Build Status](https://img.shields.io/github/actions/workflow/status/XaiPhyr/rdev-go-api-template/go-lint.yml?label=lint&logo=github&logoColor=white&style=flat)](https://github.com/XaiPhyr/rdev-go-api-template/actions)
[![Build Status](https://img.shields.io/github/actions/workflow/status/XaiPhyr/rdev-go-api-template/go-cyclo.yml?label=cyclo&logo=github&logoColor=white&style=flat)](https://github.com/XaiPhyr/rdev-go-api-template/actions)
[![Go Version](https://img.shields.io/github/go-mod/go-version/XaiPhyr/rdev-go-api-template)](https://golang.org)
[![License](https://img.shields.io/github/license/XaiPhyr/rdev-go-api-template)](LICENSE)
[![LinkedIn](https://img.shields.io/badge/LinkedIn-0A66C2?style=flat&logo=linkedin&logoColor=white)](https://linkedin.com/in/r-lozada)

---

## 🚀 Tech Stack

- **Framework:** [Gin Gonic](https://github.com/gin-gonic/gin) - Fast, lightweight HTTP web framework.
- **ORM:** [Bun](https://github.com/uptrace/bun) - SQL-first ORM for Go supporting PostgreSQL, MySQL, and SQLite.
- **Database:** PostgreSQL (Default)
- **Caching/Session:** Redis (Optional integration setup ready)
- **Config Management:** Env-based or Viper/Godotenv

---

## 📂 Project Structure

The project follows a scalable, layered monolithic layout to keep domain logic decoupled from external delivery mechanisms.

```text
rdev-go-api-template/
├── cmd/                          # Application entry points
│   ├── api/
│   │   └── main.go               # Web API server entry point
│   └── migration/
│       └── main.go               # Database migration CLI entry point
├── internal/                     # Private application code (non-importable externally)
│   ├── audit_logs/               # Audit logging domain components
│   ├── auth/                     # Authentication & authorization domain components
│   ├── config/                   # Application configuration parsing
│   │   ├── config.go             # Core config structure & parsing logic
│   │   ├── db.go                 # Database connection configurations
│   │   └── redis.go              # Redis client configurations
│   ├── db/                       # Database migrations engine
│   │   └── migrations/
│   │       └── migrations.go     # Schema setup and migration scripts
│   ├── middleware/               # Gin custom HTTP middlewares
│   │   └── middleware.go         # CORS, Recovery, JWT verification, etc.
│   ├── server/                   # HTTP server wrapper
│   │   └── routes.go             # Route grouping and engine setup
│   ├── shared/                   # Cross-cutting concerns and shared helpers
│   │   ├── aws/                  # AWS service integrations (S3, SES, etc.)
│   │   ├── dto/                  # Shared Data Transfer Objects (Request/Response shapes)
│   │   ├── email/                # Email dispatch utilities
│   │   ├── helpers/              # Cryptography, string manipulation utilities
│   │   ├── models/               # Shared Bun ORM database schemas
│   │   └── testers/              # Testing suites and mocking utilities
│   ├── templates/                # Static files and layout resources
│   │   └── landing.html          # Server-rendered HTML landing views
│   └── users/                    # User management domain (Vertical Slice)
│       ├── handler.go            # HTTP Controllers/Gin Context parsing
│       ├── mock.go               # Mock structures for unit testing
│       ├── repository.go         # Direct database access executing Bun queries
│       ├── service.go            # Core business rules processing
│       ├── types.go              # Domain-specific structures
│       ├── users_handler_test.go # HTTP Entry Point Test: Asserts status codes, JSON binding, headers, and routing using httptest.ResponseRecorder
│       └── users_service_test.go # Business Logic Test: Asserts core domain validation rules, errors, calculations, and data transformations
├── scripts/
│   └── entrypoint.sh             # Docker container initialization script
├── compose.yaml                  # Local multi-container Docker assembly (DB, Redis)
├── config.sample.yaml            # Shared application configuration blueprint
└── go.mod                        # Go module dependency manifest
```

---

## 🛠️ Getting Started

### Prerequisites

- Go `1.26` or higher
- PostgreSQL instance running locally or via Docker

### Local Installation

1. **Clone the repository:**
```bash
   git clone [https://github.com/XaiPhyr/rdev-go-api-template.git](https://github.com/XaiPhyr/rdev-go-api-template.git)
   cd rdev-go-api-template
```

2. **Setup application configuration:**
```bash
   cp config.sample.yaml config.yaml
```
Open `config.yaml` and fill in your local PostgreSQL credentials and server port configurations.

3. **Download dependencies:**

```bash
   go mod download
```

4. **Run the application:**

```bash
   go run cmd/api/main.go
```

The server should spin up by default on `http://localhost:8200`.

***

### 💡 Quick Tip for your `.gitignore`
Make sure you add `config.yaml` to your `.gitignore` file so you don't accidentally push your actual passwords, database strings, or API tokens to public GitHub!

```text
# Configuration files
config.yaml
```

---

## 🔄 Architecture Flow

Requests follow a strict, unidirectional path down the stack to ensure predictability and ease of testing:

```text
[ Client Request ]
       │
       ▼
 [ Middleware ] ──► (Auth, CORS, Rate Limiting)
       │
       ▼
  [ Handlers ]  ──► (Validates HTTP inputs, parses JSON binding)
       │
       ▼
  [ Services ]  ──► (Executes core business rules and validations)
       │
       ▼
[ Repositories ] ──► (Performs direct Bun ORM SQL statements)
       │
       ▼
  [ Database ]

```
---

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](https://www.google.com/search?q=LICENSE) file for details.
