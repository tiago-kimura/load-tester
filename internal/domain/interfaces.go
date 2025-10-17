package domain

import "context"

type HTTPClient interface {
	DoRequest(ctx context.Context, request TestRequest) TestResponse
}

type LoadTester interface {
	Execute(ctx context.Context, config LoadTestConfig) (*LoadTestResult, error)
}

type Reporter interface {
	GenerateReport(result *LoadTestResult) string
}
