package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// BenchmarkResult represents a single benchmark result
type BenchmarkResult struct {
	Name        string  `json:"name"`
	NsPerOp     float64 `json:"ns_per_op"`
	BytesPerOp  int64   `json:"bytes_per_op"`
	AllocsPerOp int64   `json:"allocs_per_op"`
	MBPerSec    float64 `json:"mb_per_sec,omitempty"`
}

// AggregatedResults holds aggregated benchmark results
type AggregatedResults struct {
	Timestamp time.Time                  `json:"timestamp"`
	Framework string                     `json:"framework"`
	Results   map[string]BenchmarkResult `json:"results"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run aggregate_results.go <results_directory>")
		os.Exit(1)
	}

	resultsDir := os.Args[1]

	// Read all result files
	files, err := filepath.Glob(filepath.Join(resultsDir, "*.txt"))
	if err != nil {
		fmt.Printf("Error reading results directory: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No result files found")
		os.Exit(1)
	}

	// Aggregate results by framework
	aggregated := make(map[string]*AggregatedResults)

	for _, file := range files {
		framework := extractFrameworkName(file)
		if framework == "" {
			continue
		}

		results, err := parseResultFile(file)
		if err != nil {
			fmt.Printf("Error parsing %s: %v\n", file, err)
			continue
		}

		aggregated[framework] = &AggregatedResults{
			Timestamp: time.Now(),
			Framework: framework,
			Results:   results,
		}
	}

	// Generate comparison report
	generateComparisonReport(aggregated)

	// Export to JSON
	exportToJSON(aggregated, filepath.Join(resultsDir, "aggregated.json"))

	fmt.Println("✓ Results aggregated successfully")
}

func extractFrameworkName(filename string) string {
	base := filepath.Base(filename)
	parts := strings.Split(base, "_")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

func parseResultFile(filename string) (map[string]BenchmarkResult, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	results := make(map[string]BenchmarkResult)
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		if !strings.HasPrefix(line, "Benchmark") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		name := fields[0]
		var nsPerOp float64
		var bytesPerOp, allocsPerOp int64

		// Parse benchmark results
		// Format: BenchmarkName-N  iterations  ns/op  B/op  allocs/op
		fmt.Sscanf(fields[2], "%f", &nsPerOp)
		if len(fields) >= 5 {
			fmt.Sscanf(fields[4], "%d", &bytesPerOp)
		}
		if len(fields) >= 7 {
			fmt.Sscanf(fields[6], "%d", &allocsPerOp)
		}

		results[name] = BenchmarkResult{
			Name:        name,
			NsPerOp:     nsPerOp,
			BytesPerOp:  bytesPerOp,
			AllocsPerOp: allocsPerOp,
		}
	}

	return results, nil
}

func generateComparisonReport(aggregated map[string]*AggregatedResults) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("Framework Comparison Report")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// Get all unique benchmark names
	benchmarkNames := make(map[string]bool)
	for _, agg := range aggregated {
		for name := range agg.Results {
			benchmarkNames[name] = true
		}
	}

	// Sort framework names
	frameworks := make([]string, 0, len(aggregated))
	for fw := range aggregated {
		frameworks = append(frameworks, fw)
	}
	sort.Strings(frameworks)

	// Print comparison for each benchmark
	for benchName := range benchmarkNames {
		fmt.Printf("### %s\n\n", benchName)
		fmt.Printf("| %-15s | %12s | %10s | %12s |\n", "Framework", "ns/op", "B/op", "allocs/op")
		fmt.Println("|" + strings.Repeat("-", 17) + "|" + strings.Repeat("-", 14) + "|" + strings.Repeat("-", 12) + "|" + strings.Repeat("-", 14) + "|")

		for _, fw := range frameworks {
			if result, ok := aggregated[fw].Results[benchName]; ok {
				fmt.Printf("| %-15s | %12.0f | %10d | %12d |\n",
					fw, result.NsPerOp, result.BytesPerOp, result.AllocsPerOp)
			}
		}
		fmt.Println()
	}
}

func exportToJSON(aggregated map[string]*AggregatedResults, filename string) {
	data, err := json.MarshalIndent(aggregated, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("Error writing JSON file: %v\n", err)
		return
	}

	fmt.Printf("✓ Results exported to: %s\n", filename)
}
