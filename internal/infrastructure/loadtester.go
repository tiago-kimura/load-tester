package infrastructure

import (
	"context"

	"github.com/tiago-kimura/load-tester/internal/domain"
)

type LoadTesterImpl struct {
	useCase domain.LoadTester
}

func NewLoadTester(useCase domain.LoadTester) *LoadTesterImpl {
	return &LoadTesterImpl{
		useCase: useCase,
	}
}

func (lt *LoadTesterImpl) Execute(ctx context.Context, config domain.LoadTestConfig) (*domain.LoadTestResult, error) {
	return lt.useCase.Execute(ctx, config)
}
