package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/tiago-kimura/load-tester/internal/domain"
)

type LoadTestUseCase struct {
	httpClient domain.HTTPClient
}

func NewLoadTestUseCase(httpClient domain.HTTPClient) *LoadTestUseCase {
	return &LoadTestUseCase{
		httpClient: httpClient,
	}
}

func (uc *LoadTestUseCase) Execute(ctx context.Context, config domain.LoadTestConfig) (*domain.LoadTestResult, error) {
	if !config.IsValid() {
		return nil, domain.ErrInvalidConfig
	}

	startTime := time.Now()

	requestChan := make(chan int, config.Requests)
	responseChan := make(chan domain.TestResponse, config.Requests)

	var wg sync.WaitGroup

	for i := 0; i < config.Concurrency; i++ {
		wg.Add(1)
		go uc.worker(ctx, &wg, requestChan, responseChan, config)
	}

	go func() {
		defer close(requestChan)
		for i := 0; i < config.Requests; i++ {
			select {
			case requestChan <- i:
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(responseChan)
	}()

	result := uc.collectResults(responseChan, startTime, config)

	return result, nil
}

func (uc *LoadTestUseCase) worker(ctx context.Context, wg *sync.WaitGroup, requestChan <-chan int, responseChan chan<- domain.TestResponse, config domain.LoadTestConfig) {
	defer wg.Done()

	for {
		select {
		case _, ok := <-requestChan:
			if !ok {
				return
			}

			request := domain.TestRequest{
				URL:     config.URL,
				Method:  "GET",
				Headers: config.Headers,
				Timeout: config.Timeout,
			}

			response := uc.httpClient.DoRequest(ctx, request)

			select {
			case responseChan <- response:
			case <-ctx.Done():
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

func (uc *LoadTestUseCase) collectResults(responseChan <-chan domain.TestResponse, startTime time.Time, config domain.LoadTestConfig) *domain.LoadTestResult {
	result := &domain.LoadTestResult{
		StatusCodes:  make(map[int]int),
		MinTime:      time.Hour,
		ErrorDetails: make([]string, 0),
	}

	var totalResponseTime time.Duration

	for response := range responseChan {
		result.TotalRequests++

		result.StatusCodes[response.StatusCode]++

		if response.Error != nil {
			result.TotalErrors++
			if config.Verbose {
				result.ErrorDetails = append(result.ErrorDetails, response.Error.Error())
			}
		}

		totalResponseTime += response.ResponseTime

		if response.ResponseTime < result.MinTime {
			result.MinTime = response.ResponseTime
		}

		if response.ResponseTime > result.MaxTime {
			result.MaxTime = response.ResponseTime
		}
	}

	result.TotalTime = time.Since(startTime)

	if result.TotalRequests > 0 {
		result.AverageTime = totalResponseTime / time.Duration(result.TotalRequests)
		result.RequestsPerSecond = float64(result.TotalRequests) / result.TotalTime.Seconds()
	}

	if result.MinTime == time.Hour {
		result.MinTime = 0
	}

	return result
}
