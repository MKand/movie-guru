# Movie Search MCP Server

Go-based MCP server for movie search using SQLite with vector embeddings.

## Prerequisites

- Go 1.22 or higher
- SQLite3
- Google Cloud credentials (for Vertex AI embeddings)

## Setup

1. **Install dependencies:**
   ```bash
   make deps
   ```

2. **Set up Google Cloud credentials:**
   ```bash
   export GOOGLE_APPLICATION_CREDENTIALS=/path/to/your/key.json
   export GOOGLE_CLOUD_PROJECT=your-project-id
   ```

3. **Build the application:**
   ```bash
   make build
   ```

## Testing

### Run all tests
```bash
make test
```

### Run tests with coverage
```bash
make test-coverage
# Opens coverage.html in your browser
```

### Run specific package tests
```bash
make test-db      # Database tests only
make test-types   # Types tests only
```

### Run quick check (lint + test)
```bash
make check
```

## Running

### Build and run
```bash
make run
```

### Run directly with go
```bash
go run ./cmd/main.go
```

## Development

### Format code
```bash
make fmt
```

### Run linters
```bash
make lint
```

### Clean build artifacts
```bash
make clean
```

## Testing without Genkit

To test the database and search functionality without full Genkit setup:

```bash
cd pkg/db
go test -v
```

This will run unit tests using an in-memory SQLite database.

## Project Structure

```
.
├── cmd/
│   └── main.go              # Main application entry
├── pkg/
│   ├── db/
│   │   ├── db.go           # Database operations
│   │   └── db_test.go      # Database tests
│   ├── embedder/
│   │   └── embedder.go     # Vertex AI embedding generation
│   ├── types/
│   │   ├── types.go        # Type definitions
│   │   └── types_test.go   # Types tests
│   └── search/
│       └── search.go       # Search logic
├── Makefile                # Build and test commands
└── go.mod                  # Dependencies
```

## Available Make Targets

Run `make help` to see all available commands.
