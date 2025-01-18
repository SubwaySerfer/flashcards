# Flashcards API

Flashcards API is a backend service for managing flashcards. It provides endpoints for creating, retrieving, updating, and deleting flashcards, as well as managing tags and learning progress.

## Features

- Create, read, update, and delete flashcards
- List flashcards with filters
- Find flashcards by tags
- Manage learning progress

## Technologies Used

- Go
- PostgreSQL
- GORM (Go ORM)
- Gin (HTTP web framework)
- Docker (for containerization)

## Getting Started

### Prerequisites

- Go 1.16+
- PostgreSQL
- Docker (optional, for containerization)

### Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/flashcards.git
   cd flashcards
   ```

2. Create a `.env` file in the project root with the following content:
   ```properties
   PORT=8081
   DATABASE_URL=
   DB_USER=postgres
   DB_PASSWORD=
   DB_HOST=localhost
   DB_PORT=5432
   DB_SSLMODE=disable
   ```

3. Run the application:
   ```sh
   go run cmd/api/main.go
   ```

### Usage

The API will be available at `http://localhost:8081/api/v1`. You can use tools like Postman or curl to interact with the endpoints.

### License

This project is licensed under the MIT License.
