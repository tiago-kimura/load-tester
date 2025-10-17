# Load Tester

CLI system in Go for performing load tests on web services, developed following SOLID principles and Clean Architecture.

## 📋 Features

### Input Parameters
- `--url`: URL of the service to be tested (required)
- `--requests`: Total number of requests (default: 100)
- `--concurrency`: Number of concurrent calls (default: 10)
- `--timeout`: Timeout per request (default: 30s)
- `--headers`: HTTP headers in format 'key1:value1,key2:value2'
- `--verbose`: Show error details and debugging

### Generated Report
- Total execution time
- Total number of requests made
- Requests with HTTP 200 status
- Distribution of HTTP status codes
- Response time statistics (average, minimum, maximum)
- Success rate and requests per second

## 🛠️ Installation and Usage

### Local Compilation

```bash
# Clone the repository
git clone https://github.com/tiago-kimura/load-tester.git
cd load-tester

# Build the application
make build

# Run a test
./load-tester --url=https://google.com --requests=1000 --concurrency=10
```

### Docker

```bash
# Build the image
docker build -t load-tester .

# Run a test
docker run --rm load-tester --url=https://google.com --requests=1000 --concurrency=10
```

## 📖 Usage Examples

### Basic Usage
```bash
./load-tester --url=https://httpbin.org/get --requests=100 --concurrency=5
```

### Test with Custom Headers
```bash
./load-tester --url=http://localhost:8080/api/test --headers='api_key:abc123,User-Agent:MyApp/1.0' --requests=100
```

### High Concurrency Test
```bash
./load-tester --url=https://api.example.com --requests=5000 --concurrency=50 --timeout=10s
```

### Debugging with Verbose
```bash
./load-tester --url=http://localhost:8080/api/test --headers='api_key:abc123' --verbose --requests=10
```

### Via Docker
```bash
docker run --rm load-tester --url=https://google.com --requests=1000 --concurrency=10
```

### Makefile (Useful Commands)
```bash
make help              # List all available commands
make run               # Build and run with example parameters
make docker-run        # Run via Docker
make example-google    # Quick test on Google
make example-httpbin   # Test on httpbin.org
```

## 🏗️ Architecture

The project follows Clean Architecture and SOLID principles:

```
├── cmd/load-tester/     # Application entry point
├── internal/
│   ├── domain/          # Domain entities and interfaces
│   ├── usecase/         # Business logic
│   └── infrastructure/ # Concrete implementations
├── pkg/                 # Reusable packages (future)
├── Dockerfile          # Containerization
└── Makefile           # Build automation
```

## 🧪 Tests

```bash
# Run tests
make test
