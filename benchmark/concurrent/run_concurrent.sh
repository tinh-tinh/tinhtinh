#!/bin/bash

# Concurrent Benchmark Runner
# Runs all concurrent benchmarks and saves results to a text file
# Usage: ./run_concurrent.sh [output_file]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

OUTPUT_FILE="${1:-concurrent_results.txt}"

{
    echo "========================================="
    echo "TinhTinh Concurrent Benchmark Report"
    echo "========================================="
    echo "Date: $(date)"
    echo "Directory: $SCRIPT_DIR"
    echo ""

    # Run all concurrent benchmarks
    echo "Running concurrent benchmarks..."
    echo "---------------------------------"
    echo ""

    go test -bench=. -benchmem -benchtime=2s -v -run=^$ .

    echo ""
    echo "========================================="
    echo "Benchmark Summary"
    echo "========================================="
    echo ""
    echo "Concurrency Levels Tested:"
    echo "  - 10 concurrent goroutines"
    echo "  - 100 concurrent goroutines"
    echo "  - 1000 concurrent goroutines"
    echo "  - 10000 concurrent goroutines"
    echo "  - RunParallel (GOMAXPROCS)"
    echo "  - With Contention (atomic counter)"
    echo "  - Context Pooling"
    echo "  - Sustained Load (5s duration)"
    echo ""
    echo "========================================="
    echo "Completed at: $(date)"
    echo "========================================="

} 2>&1 | tee "$OUTPUT_FILE"

echo ""
echo "✓ Results saved to: $OUTPUT_FILE"
