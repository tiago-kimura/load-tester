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

type Report struct {
	TotalRequests   int
	TotalDuration   time.Duration
	StatusCounts    map[int]int
	SuccessCount    int
	ErrorCount      int
	AverageDuration time.Duration
}

func main() {
	url := flag.String("url", "", "URL to test (required)")
	requests := flag.Int("requests", 100, "Total number of requests")
	concurrency := flag.Int("concurrency", 10, "Number of concurrent requests")

	flag.Parse()

	if *url == "" {
		fmt.Println("Error: --url parameter is required")
		flag.Usage()
		os.Exit(1)
	}

	if *requests <= 0 {
		fmt.Println("Error: --requests must be greater than 0")
		os.Exit(1)
	}

	if *concurrency <= 0 {
		fmt.Println("Error: --concurrency must be greater than 0")
		os.Exit(1)
	}

	fmt.Printf("Starting load test...\n")
	fmt.Printf("URL: %s\n", *url)
	fmt.Printf("Total Requests: %d\n", *requests)
	fmt.Printf("Concurrency: %d\n\n", *concurrency)

	report := runLoadTest(*url, *requests, *concurrency)
	printReport(report)
}

func runLoadTest(url string, totalRequests, concurrency int) Report {
	startTime := time.Now()

	results := make(chan Result, totalRequests)
	var wg sync.WaitGroup

	// Create a semaphore to control concurrency
	semaphore := make(chan struct{}, concurrency)

	// Launch all requests
	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := makeRequest(url)
			results <- result
		}()
	}

	// Wait for all requests to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	statusCounts := make(map[int]int)
	successCount := 0
	errorCount := 0
	var totalDuration time.Duration

	for result := range results {
		if result.Error != nil {
			errorCount++
		} else {
			statusCounts[result.StatusCode]++
			// Consider 2xx status codes as successful
			if result.StatusCode >= 200 && result.StatusCode < 300 {
				successCount++
			}
		}
		totalDuration += result.Duration
	}

	elapsedTime := time.Since(startTime)
	avgDuration := time.Duration(0)
	if totalRequests > 0 {
		avgDuration = totalDuration / time.Duration(totalRequests)
	}

	return Report{
		TotalRequests:   totalRequests,
		TotalDuration:   elapsedTime,
		StatusCounts:    statusCounts,
		SuccessCount:    successCount,
		ErrorCount:      errorCount,
		AverageDuration: avgDuration,
	}
}

func makeRequest(url string) Result {
	startTime := time.Now()

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	duration := time.Since(startTime)

	if err != nil {
		return Result{
			Duration: duration,
			Error:    err,
		}
	}
	defer resp.Body.Close()

	return Result{
		StatusCode: resp.StatusCode,
		Duration:   duration,
		Error:      nil,
	}
}

func printReport(report Report) {
	fmt.Println("===============================================")
	fmt.Println("Load Test Report")
	fmt.Println("===============================================")
	fmt.Printf("Total Requests:       %d\n", report.TotalRequests)
	fmt.Printf("Total Time:           %v\n", report.TotalDuration)
	fmt.Printf("Average Request Time: %v\n", report.AverageDuration)
	fmt.Printf("Requests per Second:  %.2f\n\n", float64(report.TotalRequests)/report.TotalDuration.Seconds())

	fmt.Println("Status Code Distribution:")
	fmt.Println("-----------------------------------------------")
	for status, count := range report.StatusCounts {
		percentage := float64(count) / float64(report.TotalRequests) * 100
		fmt.Printf("  Status %d: %d requests (%.2f%%)\n", status, count, percentage)
	}

	if report.ErrorCount > 0 {
		percentage := float64(report.ErrorCount) / float64(report.TotalRequests) * 100
		fmt.Printf("  Errors:    %d requests (%.2f%%)\n", report.ErrorCount, percentage)
	}

	fmt.Println("\n===============================================")
	fmt.Printf("Success (HTTP 2xx): %d/%d (%.2f%%)\n",
		report.SuccessCount,
		report.TotalRequests,
		float64(report.SuccessCount)/float64(report.TotalRequests)*100)
	fmt.Println("===============================================")
}
