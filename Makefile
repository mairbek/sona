.PHONY: build run proto deps

# Variables
IMAGE_NAME = sona
PORT = 8080
GOPATH ?= $(shell go env GOPATH)
BUF_BIN = $(GOPATH)/bin/buf
PROTOC_GEN_GO = $(GOPATH)/bin/protoc-gen-go
PROTOC_GEN_CONNECT_GO = $(GOPATH)/bin/protoc-gen-connect-go
PATH := $(GOPATH)/bin:$(PATH)

# Install dependencies
deps:
	go mod download
	go install github.com/bufbuild/buf/cmd/buf@v1.28.1
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
	go mod tidy
	@echo "Please ensure $(GOPATH)/bin is in your PATH"

# Generate protobuf files
proto: $(PROTOC_GEN_GO) $(PROTOC_GEN_CONNECT_GO)
	PATH=$(GOPATH)/bin:$(PATH) $(BUF_BIN) generate proto
	go mod tidy

# Build the Docker image
build: proto
	go build -o sona main.go

run:
	go run main.go

test:
	go test -v -parallel=4 -timeout 30s ./...