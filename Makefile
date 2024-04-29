# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get

# Binary name
BINARY_NAME = chat-ws

# Build target
MAIN_PATH = cmd/main.go

build:
	$(GOBUILD) -o ./build/$(BINARY_NAME) -v $(MAIN_PATH)

# Clean target
clean:
	$(GOCLEAN)
	rm -f ./build/$(BINARY_NAME)

# Test target
test:
	$(GOTEST) -v ./...

# Install dependencies
deps:
	$(GOGET) -v ./...

# Default target
default: build

docker-build:
	docker build -f ./docker/Dockerfile . -t chat-ws

docker-run:
	docker run -p 8080:8080 chat-ws

gen-proto:
	protoc  --go_out=./models --go_opt=paths=source_relative ./chat-protobuf/*.proto


.PHONY: build clean test deps