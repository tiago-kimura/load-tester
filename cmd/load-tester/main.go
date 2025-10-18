package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/tiago-kimura/load-tester/internal/domain"
	"github.com/tiago-kimura/load-tester/internal/infrastructure"
	"github.com/tiago-kimura/load-tester/internal/usecase"
)

func main() {
	var (
		url         = flag.String("url", "", "URL of the service to be tested")
		requests    = flag.Int("requests", 100, "Total number of requests")
		concurrency = flag.Int("concurrency", 10, "Number of concurrent requests")
		timeout     = flag.Duration("timeout", 30*time.Second, "Timeout per request")
		headers     = flag.String("headers", "", "HTTP headers in format 'key1:value1,key2:value2'")
		verbose     = flag.Bool("verbose", false, "Show error details and debugging")
		help        = flag.Bool("help", false, "Show this help message")
	)

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *url == "" {
		fmt.Fprintf(os.Stderr, "Error: URL is required\n")
		showUsage()
		os.Exit(1)
	}

	if *requests <= 0 {
		fmt.Fprintf(os.Stderr, "Error: Number of requests must be greater than 0\n")
		showUsage()
		os.Exit(1)
	}

	if *concurrency <= 0 {
		fmt.Fprintf(os.Stderr, "Error: Concurrency must be greater than 0\n")
		showUsage()
		os.Exit(1)
	}

	if *concurrency > *requests {
		*concurrency = *requests
	}

	parsedHeaders := parseHeaders(*headers)

	config := domain.LoadTestConfig{
		URL:         *url,
		Requests:    *requests,
		Concurrency: *concurrency,
		Timeout:     *timeout,
		Headers:     parsedHeaders,
		Verbose:     *verbose,
	}

	httpClient := infrastructure.NewHTTPClient(*timeout)
	loadTestUseCase := usecase.NewLoadTestUseCase(httpClient)
	reporter := infrastructure.NewConsoleReporter()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nReceived interruption signal, cancelling test...")
		cancel()
	}()

	fmt.Printf("Starting load test...\n")
	fmt.Printf("URL: %s\n", config.URL)
	fmt.Printf("Requests: %d\n", config.Requests)
	fmt.Printf("Concurrency: %d\n", config.Concurrency)
	fmt.Printf("Timeout: %v\n", config.Timeout)
	fmt.Println("=====================================")

	result, err := loadTestUseCase.Execute(ctx, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing test: %v\n", err)
		os.Exit(1)
	}

	report := reporter.GenerateReport(result)
	fmt.Print(report)
}

func showHelp() {
	fmt.Println("Load Tester - CLI tool for load testing")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  load-tester --url=<URL> --requests=<NUM> --concurrency=<NUM> [OPTIONS]")
	fmt.Println()
	fmt.Println("REQUIRED FLAGS:")
	fmt.Println("  --url string        URL of the service to be tested")
	fmt.Println()
	fmt.Println("OPTIONAL FLAGS:")
	fmt.Println("  --requests int      Total number of requests (default: 100)")
	fmt.Println("  --concurrency int   Number of concurrent requests (default: 10)")
	fmt.Println("  --timeout duration  Timeout per request (default: 30s)")
	fmt.Println("  --headers string    HTTP headers in format 'key1:value1,key2:value2'")
	fmt.Println("  --verbose          Show error details and debugging")
	fmt.Println("  --help             Show this help message")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  load-tester --url=http://google.com --requests=1000 --concurrency=10")
	fmt.Println("  load-tester --url=https://api.example.com --requests=500 --concurrency=20 --timeout=10s")
	fmt.Println("  load-tester --url=http://localhost:8080/api/test --headers='api_key:abc123' --requests=100")
	fmt.Println()
	fmt.Println("DOCKER USAGE:")
	fmt.Println("  docker run <image> --url=http://google.com --requests=1000 --concurrency=10")
}

func showUsage() {
	fmt.Fprintf(os.Stderr, "\nUse 'load-tester --help' to see all available options.\n")
}

func parseHeaders(headerString string) map[string]string {
	headers := make(map[string]string)

	if headerString == "" {
		return headers
	}

	pairs := strings.Split(headerString, ",")

	for _, pair := range pairs {
		parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key != "" && value != "" {
				headers[key] = value
			}
		}
	}

	return headers
}
