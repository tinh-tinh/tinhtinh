package latency

import (
	"fmt"
	"strings"
)

// GenerateReport generates a formatted latency report
func GenerateReport(stats LatencyStats, title string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("=== %s ===\n\n", title))
	sb.WriteString("Latency Distribution (milliseconds):\n")
	sb.WriteString("------------------------------------\n")
	sb.WriteString(fmt.Sprintf("Min:    %.3f ms\n", ToMilliseconds(stats.Min)))
	sb.WriteString(fmt.Sprintf("p50:    %.3f ms\n", ToMilliseconds(stats.P50)))
	sb.WriteString(fmt.Sprintf("p95:    %.3f ms\n", ToMilliseconds(stats.P95)))
	sb.WriteString(fmt.Sprintf("p99:    %.3f ms\n", ToMilliseconds(stats.P99)))
	sb.WriteString(fmt.Sprintf("p99.9:  %.3f ms\n", ToMilliseconds(stats.P999)))
	sb.WriteString(fmt.Sprintf("Max:    %.3f ms\n", ToMilliseconds(stats.Max)))
	sb.WriteString(fmt.Sprintf("Mean:   %.3f ms\n", ToMilliseconds(stats.Mean)))
	sb.WriteString(fmt.Sprintf("StdDev: %.3f ms\n", ToMilliseconds(stats.StdDev)))
	sb.WriteString("\n")

	return sb.String()
}

// GenerateMarkdownTable generates a markdown table of latency stats
func GenerateMarkdownTable(results map[string]LatencyStats) string {
	var sb strings.Builder

	sb.WriteString("| Benchmark | Min (ms) | p50 (ms) | p95 (ms) | p99 (ms) | p99.9 (ms) | Max (ms) | Mean (ms) |\n")
	sb.WriteString("|-----------|----------|----------|----------|----------|------------|----------|----------|\n")

	for name, stats := range results {
		sb.WriteString(fmt.Sprintf("| %s | %.3f | %.3f | %.3f | %.3f | %.3f | %.3f | %.3f |\n",
			name,
			ToMilliseconds(stats.Min),
			ToMilliseconds(stats.P50),
			ToMilliseconds(stats.P95),
			ToMilliseconds(stats.P99),
			ToMilliseconds(stats.P999),
			ToMilliseconds(stats.Max),
			ToMilliseconds(stats.Mean),
		))
	}

	return sb.String()
}

// GenerateHistogram generates a simple ASCII histogram
func GenerateHistogram(durations []float64, buckets int) string {
	if len(durations) == 0 {
		return "No data"
	}

	// Handle single data point or all same values
	if len(durations) == 1 {
		return fmt.Sprintf("\nLatency Histogram:\n------------------\n%.2f ms [%5d]: %s\n",
			ToMilliseconds(durations[0]), 1, strings.Repeat("█", 50))
	}

	stats := Calculate(durations)
	min := stats.Min
	max := stats.Max

	// If all values are the same
	if min == max {
		return fmt.Sprintf("\nLatency Histogram:\n------------------\n%.2f ms [%5d]: %s\n",
			ToMilliseconds(min), len(durations), strings.Repeat("█", 50))
	}

	bucketSize := (max - min) / float64(buckets)

	// Create buckets
	counts := make([]int, buckets)
	for _, d := range durations {
		bucketIndex := int((d - min) / bucketSize)
		if bucketIndex >= buckets {
			bucketIndex = buckets - 1
		}
		if bucketIndex < 0 {
			bucketIndex = 0
		}
		counts[bucketIndex]++
	}

	// Find max count for scaling
	maxCount := 0
	for _, count := range counts {
		if count > maxCount {
			maxCount = count
		}
	}

	// Generate histogram
	var sb strings.Builder
	sb.WriteString("\nLatency Histogram:\n")
	sb.WriteString("------------------\n")

	const maxBarWidth = 50
	for i, count := range counts {
		bucketStart := min + float64(i)*bucketSize
		bucketEnd := bucketStart + bucketSize

		barWidth := 0
		if maxCount > 0 {
			barWidth = (count * maxBarWidth) / maxCount
		}

		bar := strings.Repeat("█", barWidth)
		sb.WriteString(fmt.Sprintf("%.2f-%.2f ms [%5d]: %s\n",
			ToMilliseconds(bucketStart),
			ToMilliseconds(bucketEnd),
			count,
			bar,
		))
	}

	return sb.String()
}
