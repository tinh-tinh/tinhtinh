#!/bin/bash

# Goroutine Profile Analyzer
# Generates goroutine profile and analyzes the data
# Usage: ./analyze_goroutine.sh [output_file]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

OUTPUT_FILE="${1:-goroutine_results.txt}"
PROFILE_FILE="goroutine.prof"

{
    echo "========================================="
    echo "Goroutine Profile Analysis"
    echo "========================================="
    echo "Date: $(date)"
    echo ""

    # Generate goroutine profile first
    echo "Generating goroutine profile..."
    echo "-------------------------------"
    go test -bench=BenchmarkGoroutineProfile -benchtime=3s -run=^$ .
    echo ""

    if [ ! -f "$PROFILE_FILE" ]; then
        echo "Error: Profile file '$PROFILE_FILE' was not generated"
        exit 1
    fi

    # Display goroutine information
    echo "Goroutine statistics:"
    echo "--------------------"
    go tool pprof -top -nodecount=20 "$PROFILE_FILE"
    echo ""

    # List goroutines
    echo "Goroutine details:"
    echo "-----------------"
    go tool pprof -list=. "$PROFILE_FILE" 2>/dev/null | head -50 || true
    echo ""

    # Generate visualization
    if command -v dot &> /dev/null; then
        echo "Generating goroutine graph..."
        go tool pprof -svg "$PROFILE_FILE" > goroutine_graph.svg
        echo "✓ Goroutine graph saved to: goroutine_graph.svg"
        echo ""
    fi

    # Summary
    echo "========================================="
    echo "Analysis Complete"
    echo "========================================="
    echo "Generated files:"
    ls -lh "$PROFILE_FILE" goroutine_graph.svg 2>/dev/null || true
    echo ""
    echo "For interactive analysis, run:"
    echo "  go tool pprof -http=:8080 $PROFILE_FILE"
    echo "========================================="

} 2>&1 | tee "$OUTPUT_FILE"

echo ""
echo "✓ Results saved to: $OUTPUT_FILE"
