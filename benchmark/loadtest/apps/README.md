# Load Test Applications

This directory contains sample HTTP server applications for load testing.

## Applications

Each application is in its own subdirectory and implements the same API endpoints on different ports:

| Framework  | Port | Directory | File       |
|------------|------|-----------|------------|
| Tinh Tinh  | 3000 | tinhtinh/ | main.go    |
| Gin        | 3001 | gin/      | main.go    |
| Echo       | 3002 | echo/     | main.go    |
| Fiber      | 3003 | fiber/    | main.go    |
| Chi        | 3004 | chi/      | main.go    |

## API Endpoints

All applications expose the following endpoints:

- `GET /api/` - Simple text response
- `GET /api/json` - JSON response
- `POST /api/json` - JSON request/response
- `GET /api/user/:id` - Path parameter handling

## Running Applications

### Run Tinh Tinh App
```bash
cd apps/tinhtinh
go run main.go
```

### Run Gin App
```bash
cd apps/gin
go run main.go
```

### Run All Apps (in separate terminals)
```bash
# Terminal 1
cd apps/tinhtinh && go run main.go

# Terminal 2
cd apps/gin && go run main.go

# Terminal 3
cd apps/echo && go run main.go

# Terminal 4
cd apps/fiber && go run main.go

# Terminal 5
cd apps/chi && go run main.go
```

## Testing

Test if an app is running:
```bash
curl http://localhost:3000/api/
curl http://localhost:3000/api/json
```

## Dependencies

Make sure to install dependencies first:
```bash
cd ../../frameworks
go mod download
```
