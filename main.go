package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type Result struct {
	StatusCode int
	Duration   time.Duration
	Error      error
}

func main() {
	url := flag.String("url", "", "URL to test (required)")
	requests := flag.Int("requests", 100, "Number of requests to send")
	concurrency := flag.Int("concurrency", 10, "Number of concurrent workers")
	flag.Parse()

	if *url == "" {
		fmt.Println("Error: URL is required")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("Load Testing: %s\n", *url)
	fmt.Printf("Total Requests: %d\n", *requests)
	fmt.Printf("Concurrency: %d\n\n", *concurrency)

	results := runLoadTest(*url, *requests, *concurrency)
	printResults(results)
}

func runLoadTest(url string, totalRequests, concurrency int) []Result {
	results := make([]Result, totalRequests)
	var wg sync.WaitGroup
	requestChan := make(chan int, totalRequests)

	// Start workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range requestChan {
				results[idx] = makeRequest(url)
			}
		}()
	}

	// Send requests
	for i := 0; i < totalRequests; i++ {
		requestChan <- i
	}
	close(requestChan)

	wg.Wait()
	return results
}

func makeRequest(url string) Result {
	start := time.Now()
	resp, err := http.Get(url)
	duration := time.Since(start)

	if err != nil {
		return Result{
			StatusCode: 0,
			Duration:   duration,
			Error:      err,
		}
	}
	defer resp.Body.Close()

	return Result{
		StatusCode: resp.StatusCode,
		Duration:   duration,
		Error:      nil,
	}
}

func printResults(results []Result) {
	var totalDuration time.Duration
	statusCodes := make(map[int]int)
	errorCount := 0

	for _, result := range results {
		totalDuration += result.Duration
		if result.Error != nil {
			errorCount++
		} else {
			statusCodes[result.StatusCode]++
		}
	}

	fmt.Println("Results:")
	fmt.Println("--------")
	fmt.Printf("Total requests: %d\n", len(results))
	fmt.Printf("Average response time: %v\n", totalDuration/time.Duration(len(results)))
	fmt.Printf("Total time: %v\n\n", totalDuration)

	fmt.Println("Status codes:")
	for code, count := range statusCodes {
		fmt.Printf("  %d: %d\n", code, count)
	}

	if errorCount > 0 {
		fmt.Printf("\nErrors: %d\n", errorCount)
	}
}
