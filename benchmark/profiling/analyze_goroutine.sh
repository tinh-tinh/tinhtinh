#!/bin/bash

# Goroutine Profile Analyzer
# Analyzes goroutine profiling data

set -e

PROFILE_FILE="${1:-goroutine.prof}"

if [ ! -f "$PROFILE_FILE" ]; then
    echo "Error: Profile file '$PROFILE_FILE' not found"
    echo "Usage: $0 [profile_file]"
    echo "Run benchmarks first: go test -bench=BenchmarkGoroutineProfile"
    exit 1
fi

echo "========================================="
echo "Goroutine Profile Analysis"
echo "========================================="
echo ""

# Display goroutine information
echo "Goroutine statistics:"
echo "--------------------"
go tool pprof -top -nodecount=20 "$PROFILE_FILE"
echo ""

# List goroutines
echo "Goroutine details:"
echo "-----------------"
go tool pprof -list=. "$PROFILE_FILE" | head -50
echo ""

# Generate visualization
if command -v dot &> /dev/null; then
    echo "Generating goroutine graph..."
    go tool pprof -svg "$PROFILE_FILE" > goroutine_graph.svg
    echo "✓ Goroutine graph saved to: goroutine_graph.svg"
    echo ""
fi

# Interactive mode instructions
echo "========================================="
echo "For interactive analysis, run:"
echo "  go tool pprof $PROFILE_FILE"
echo ""
echo "Or for web UI:"
echo "  go tool pprof -http=:8080 $PROFILE_FILE"
echo "========================================="
