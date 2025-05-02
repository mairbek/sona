# Sona Web Server

A Go web server with ConnectRPC support.

## Requirements

- Go 1.23 or later
- Docker (optional)
- Make (optional)

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

### Using Docker

#### Using Docker commands directly

1. Build the Docker image:
```bash
docker build -t sona .
```

2. Run the container:
```bash
docker run -p 8080:8080 sona
```

#### Using Make commands

1. Build the Docker image:
```bash
make build
```

2. Run the container:
```bash
make run
```

3. Clean up (optional):
```bash
make clean
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
  -d '{"name": "World"}' \
  http://localhost:8080/sona.v1.HelloService/Hello
```

Expected response:
```json
{"message": "Hello, World!"}
``` 