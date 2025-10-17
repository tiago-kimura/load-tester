.PHONY: build run test clean docker-build docker-run help

BINARY_NAME=load-tester
DOCKER_IMAGE=load-tester:latest
MAIN_PATH=./cmd/load-tester

help: 
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: 
	@echo "Compilando aplicação..."
	go build -o $(BINARY_NAME) $(MAIN_PATH)

run: build 
	@echo "Executando aplicação..."
	./$(BINARY_NAME) --url=https://httpbin.org/get --requests=10 --concurrency=3

test: 
	@echo "Executando testes..."
	go test -v ./...

clean: 
	@echo "Limpando arquivos de build..."
	rm -f $(BINARY_NAME)

docker-build: 
	@echo "Construindo imagem Docker..."
	docker build -t $(DOCKER_IMAGE) .

docker-run: docker-build 
	@echo "Executando aplicação via Docker..."
	docker run --rm $(DOCKER_IMAGE) --url=https://httpbin.org/get --requests=10 --concurrency=3

example-google: build 
	./$(BINARY_NAME) --url=https://google.com --requests=100 --concurrency=10

example-httpbin: build 
	./$(BINARY_NAME) --url=https://httpbin.org/status/200 --requests=50 --concurrency=5