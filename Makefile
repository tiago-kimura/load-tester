.PHONY: run build clean test

run:
	@go run main.go -url=https://www.google.com -requests=10 -concurrency=2

build:
	@go build -o load-tester main.go

clean:
	@rm -f load-tester

test:
	@go test ./...
