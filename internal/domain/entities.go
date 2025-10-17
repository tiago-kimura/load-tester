package domain

import (
	"time"
)

type TestRequest struct {
	URL     string
	Method  string
	Headers map[string]string
	Body    []byte
	Timeout time.Duration
}

type TestResponse struct {
	StatusCode   int
	ResponseTime time.Duration
	Error        error
	BodySize     int64
}

type LoadTestConfig struct {
	URL         string
	Requests    int
	Concurrency int
	Timeout     time.Duration
	Headers     map[string]string
	Verbose     bool
}

type LoadTestResult struct {
	TotalRequests     int
	TotalTime         time.Duration
	StatusCodes       map[int]int
	TotalErrors       int
	AverageTime       time.Duration
	MinTime           time.Duration
	MaxTime           time.Duration
	RequestsPerSecond float64
	ErrorDetails      []string
}

func (c *LoadTestConfig) IsValid() bool {
	return c.URL != "" && c.Requests > 0 && c.Concurrency > 0
}
