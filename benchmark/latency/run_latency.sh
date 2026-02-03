#!/bin/bash

# Script to run latency benchmarks and save results to a text file
# Usage: ./run_benchmark.sh [output_file]

# Default output file
OUTPUT_FILE="${1:-benchmark_results.txt}"

# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Change to the benchmark directory
cd "$SCRIPT_DIR" || exit 1

# Print header
echo "Running Latency Benchmarks..."
echo "Results will be saved to: $OUTPUT_FILE"
echo ""

# Run benchmarks and save to file
{
    echo "======================================"
    echo "TinhTinh Latency Benchmark Results"
    echo "======================================"
    echo "Date: $(date)"
    echo "Directory: $SCRIPT_DIR"
    echo ""
    echo "======================================"
    echo ""
    
    # Run all benchmarks with verbose output
    go test -bench=. -benchmem -benchtime=1s -v
    
    echo ""
    echo "======================================"
    echo "Benchmark completed at: $(date)"
    echo "======================================"
} 2>&1 | tee "$OUTPUT_FILE"

echo ""
echo "✓ Benchmark results saved to: $OUTPUT_FILE"
