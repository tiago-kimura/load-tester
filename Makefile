.PHONY: build run test clean docker-build docker-run help

BINARY_NAME=load-tester
DOCKER_IMAGE=load-tester:latest
MAIN_PATH=./cmd/load-tester

# Default commands
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building application..."
	go build -o $(BINARY_NAME) $(MAIN_PATH)

run: build ## Build and run the application with example parameters
	@echo "Running application..."
	./$(BINARY_NAME) --url=https://httpbin.org/get --requests=10 --concurrency=3

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

clean: ## Remove build artifacts
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

docker-run: docker-build ## Build and run application via Docker
	@echo "Running application via Docker..."
	docker run --rm $(DOCKER_IMAGE) --url=https://httpbin.org/get --requests=10 --concurrency=3

# Examples with different parameters
example-google: build ## Test Google with 100 requests
	./$(BINARY_NAME) --url=https://google.com --requests=100 --concurrency=10

example-httpbin: build ## Test httpbin.org with 50 requests
	./$(BINARY_NAME) --url=https://httpbin.org/status/200 --requests=50 --concurrency=5