# Sona Web Server

Go+connectrpc+pgx+sqlc+testcontainers+pgtestdb

## Requirements

- Go 1.23 or later
- Make (optional)
- Nix (optional)

## Setup

1. Install dependencies:
```bash
make deps
```

2. Generate protobuf files:
```bash
make proto
```

## Running the Server

### Using Go directly

1. Start the server:
```bash
go run main.go
```

2. The server will be available at:
```
http://localhost:8080
```

#### Using Make commands

1. Build the binary
```bash
make build
```

2. Run the server:
```bash
make run
```

3. Run tests
```bash
make test
```

4. View available commands:
```bash
make help
```

## API Documentation

The server exposes a ConnectRPC endpoint at `/sona.v1.HelloService/`.

### Hello Service

- **Endpoint**: `/sona.v1.HelloService/Hello`
- **Method**: POST
- **Request Body**:
  ```json
  {
    "name": "string"
  }
  ```
- **Response Body**:
  ```json
  {
    "message": "string"
  }
  ```

## Example Usage

You can test the API using curl:

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"name": "Pepe"}' \
  http://localhost:8080/sona.v1.HelloService/Hello
```

Expected response:
```json
{"message": "Hello, Pepe!"}
``` 