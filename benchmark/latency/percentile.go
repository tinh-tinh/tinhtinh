package latency

import (
	"math"
	"sort"
)

// Percentile calculates the percentile value from a sorted slice of durations
func Percentile(sortedData []float64, percentile float64) float64 {
	if len(sortedData) == 0 {
		return 0
	}

	if percentile <= 0 {
		return sortedData[0]
	}
	if percentile >= 100 {
		return sortedData[len(sortedData)-1]
	}

	index := (percentile / 100.0) * float64(len(sortedData)-1)
	lower := int(math.Floor(index))
	upper := int(math.Ceil(index))

	if lower == upper {
		return sortedData[lower]
	}

	// Linear interpolation
	fraction := index - float64(lower)
	return sortedData[lower]*(1-fraction) + sortedData[upper]*fraction
}

// LatencyStats holds latency statistics
type LatencyStats struct {
	Min    float64
	Max    float64
	Mean   float64
	P50    float64
	P95    float64
	P99    float64
	P999   float64
	StdDev float64
}

// Calculate calculates latency statistics from a slice of durations (in nanoseconds)
func Calculate(durations []float64) LatencyStats {
	if len(durations) == 0 {
		return LatencyStats{}
	}

	// Sort data
	sorted := make([]float64, len(durations))
	copy(sorted, durations)
	sort.Float64s(sorted)

	// Calculate min, max, mean
	min := sorted[0]
	max := sorted[len(sorted)-1]

	sum := 0.0
	for _, v := range sorted {
		sum += v
	}
	mean := sum / float64(len(sorted))

	// Calculate standard deviation
	variance := 0.0
	for _, v := range sorted {
		diff := v - mean
		variance += diff * diff
	}
	stdDev := math.Sqrt(variance / float64(len(sorted)))

	// Calculate percentiles
	p50 := Percentile(sorted, 50)
	p95 := Percentile(sorted, 95)
	p99 := Percentile(sorted, 99)
	p999 := Percentile(sorted, 99.9)

	return LatencyStats{
		Min:    min,
		Max:    max,
		Mean:   mean,
		P50:    p50,
		P95:    p95,
		P99:    p99,
		P999:   p999,
		StdDev: stdDev,
	}
}

// ToMilliseconds converts nanoseconds to milliseconds
func ToMilliseconds(ns float64) float64 {
	return ns / 1_000_000.0
}

// ToMicroseconds converts nanoseconds to microseconds
func ToMicroseconds(ns float64) float64 {
	return ns / 1_000.0
}
