package infrastructure

import (
	"fmt"
	"strings"

	"github.com/tiago-kimura/load-tester/internal/domain"
)

type ConsoleReporter struct{}

func NewConsoleReporter() *ConsoleReporter {
	return &ConsoleReporter{}
}

func (r *ConsoleReporter) GenerateReport(result *domain.LoadTestResult) string {
	var report strings.Builder

	report.WriteString("\n")
	report.WriteString("===============================================\n")
	report.WriteString("           LOAD TEST REPORT\n")
	report.WriteString("===============================================\n\n")

	report.WriteString("GENERAL STATISTICS:\n")
	report.WriteString(fmt.Sprintf("• Total requests:       %d\n", result.TotalRequests))
	report.WriteString(fmt.Sprintf("• Total time:           %v\n", result.TotalTime))
	report.WriteString(fmt.Sprintf("• Requests per second:  %.2f\n", result.RequestsPerSecond))
	report.WriteString(fmt.Sprintf("• Total errors:         %d\n", result.TotalErrors))
	report.WriteString("\n")

	report.WriteString("RESPONSE TIME STATISTICS:\n")
	report.WriteString(fmt.Sprintf("• Average time:         %v\n", result.AverageTime))
	report.WriteString(fmt.Sprintf("• Minimum time:         %v\n", result.MinTime))
	report.WriteString(fmt.Sprintf("• Maximum time:         %v\n", result.MaxTime))
	report.WriteString("\n")

	report.WriteString("HTTP STATUS CODE DISTRIBUTION:\n")

	if count, exists := result.StatusCodes[200]; exists {
		report.WriteString(fmt.Sprintf("• HTTP 200 (Success):   %d requests\n", count))
	}

	for statusCode, count := range result.StatusCodes {
		if statusCode != 200 {
			statusName := getStatusCodeName(statusCode)
			report.WriteString(fmt.Sprintf("• HTTP %d (%s): %d requests\n", statusCode, statusName, count))
		}
	}

	successCount := result.StatusCodes[200]
	successRate := float64(successCount) / float64(result.TotalRequests) * 100
	report.WriteString(fmt.Sprintf("\nSUCCESS RATE: %.2f%% (%d out of %d requests)\n",
		successRate, successCount, result.TotalRequests))

	if len(result.ErrorDetails) > 0 {
		report.WriteString("\nERROR DETAILS:\n")
		for i, errorDetail := range result.ErrorDetails {
			if i >= 10 {
				report.WriteString(fmt.Sprintf("... and %d more errors\n", len(result.ErrorDetails)-10))
				break
			}
			report.WriteString(fmt.Sprintf("• %s\n", errorDetail))
		}
	}

	report.WriteString("\n===============================================\n")

	return report.String()
}

func getStatusCodeName(statusCode int) string {
	switch statusCode {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	case 502:
		return "Bad Gateway"
	case 503:
		return "Service Unavailable"
	case 504:
		return "Gateway Timeout"
	default:
		return "Unknown"
	}
}
