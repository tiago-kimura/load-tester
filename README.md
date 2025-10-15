# load-tester
CLI system in Go to perform load testing on a web service

## Description

A command-line load testing tool written in Go that allows you to test web services by sending concurrent HTTP requests and generating detailed performance reports.

## Features

- Concurrent HTTP requests with configurable concurrency level
- Detailed performance metrics (response times, request rates)
- Status code distribution analysis
- Success/failure tracking
- Docker support for easy deployment

## Prerequisites

- Go 1.24+ (for local development)
- Docker (for containerized execution)

## Installation

### Option 1: Build from source

```bash
go build -o load-tester main.go
```

### Option 2: Using Docker

```bash
docker build -t load-tester .
```

## Usage

### Command Line Parameters

- `--url` (required): The URL to test
- `--requests` (optional): Total number of requests to make (default: 100)
- `--concurrency` (optional): Number of concurrent requests (default: 10)

### Examples

#### Running locally:

```bash
# Basic usage
./load-tester --url=http://example.com --requests=100 --concurrency=10

# Test with 1000 requests and 50 concurrent connections
./load-tester --url=https://api.example.com/endpoint --requests=1000 --concurrency=50

# Quick test with defaults
./load-tester --url=http://localhost:8080
```

#### Running with Docker:

```bash
# Basic usage
docker run load-tester --url=http://google.com --requests=1000 --concurrency=10

# Testing a local service (use host.docker.internal on Mac/Windows or host network on Linux)
docker run --network="host" load-tester --url=http://localhost:8080 --requests=500 --concurrency=5

# Test external API
docker run load-tester --url=https://api.github.com --requests=100 --concurrency=5
```

## Output

The tool provides a comprehensive report including:

- Total number of requests made
- Total execution time
- Average request time
- Requests per second
- Status code distribution with percentages
- Success rate (HTTP 200 responses)
- Error count

### Sample Output:

```
Starting load test...
URL: http://example.com
Total Requests: 100
Concurrency: 10

===============================================
Load Test Report
===============================================
Total Requests:       100
Total Time:           2.5s
Average Request Time: 25ms
Requests per Second:  40.00

Status Code Distribution:
-----------------------------------------------
  Status 200: 95 requests (95.00%)
  Status 500: 5 requests (5.00%)

===============================================
Success (HTTP 200): 95/100 (95.00%)
===============================================
```

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o load-tester main.go
```

## License

This project is open source and available under the MIT License.

