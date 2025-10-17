package infrastructure

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/tiago-kimura/load-tester/internal/domain"
)

type HTTPClientImpl struct {
	client *http.Client
}

func NewHTTPClient(timeout time.Duration) *HTTPClientImpl {
	return &HTTPClientImpl{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (h *HTTPClientImpl) DoRequest(ctx context.Context, request domain.TestRequest) domain.TestResponse {
	startTime := time.Now()

	var body io.Reader
	if request.Body != nil {
		body = bytes.NewReader(request.Body)
	}

	req, err := http.NewRequestWithContext(ctx, request.Method, request.URL, body)
	if err != nil {
		return domain.TestResponse{
			StatusCode:   0,
			ResponseTime: time.Since(startTime),
			Error:        err,
			BodySize:     0,
		}
	}

	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	resp, err := h.client.Do(req)
	responseTime := time.Since(startTime)

	if err != nil {
		return domain.TestResponse{
			StatusCode:   0,
			ResponseTime: responseTime,
			Error:        err,
			BodySize:     0,
		}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	bodySize := int64(0)
	if err == nil {
		bodySize = int64(len(bodyBytes))
	}

	return domain.TestResponse{
		StatusCode:   resp.StatusCode,
		ResponseTime: responseTime,
		Error:        nil,
		BodySize:     bodySize,
	}
}
