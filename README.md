# Sona Web Server

A simple Go web server that responds with "Hello, World!".

## Requirements

- Go 1.22 or later
- Docker (optional)
- Make (optional)

## Running the Server

### Using Go directly

1. Start the server:
```bash
go run main.go
```

2. Open your browser and visit:
```
http://localhost:8080
```

You should see "Hello, World!" displayed in your browser.

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

3. Open your browser and visit:
```
http://localhost:8080
```

You should see "Hello, World!" displayed in your browser. 