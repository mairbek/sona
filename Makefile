.PHONY: build run clean

# Variables
IMAGE_NAME = sona
PORT = 8080

# Build the Docker image
build:
	docker build -t $(IMAGE_NAME) .

# Run the Docker container
run:
	docker run -p $(PORT):$(PORT) $(IMAGE_NAME)

# Clean up Docker resources
clean:
	docker rmi $(IMAGE_NAME) || true

# Help command
help:
	@echo "Available commands:"
	@echo "  make build    - Build the Docker image"
	@echo "  make run      - Run the Docker container"
	@echo "  make clean    - Remove the Docker image" 