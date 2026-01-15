#!/bin/bash

# CPU Profile Analyzer
# Analyzes CPU profiling data and generates reports

set -e

PROFILE_FILE="${1:-cpu.prof}"

if [ ! -f "$PROFILE_FILE" ]; then
    echo "Error: Profile file '$PROFILE_FILE' not found"
    echo "Usage: $0 [profile_file]"
    echo "Run benchmarks first: go test -bench=BenchmarkCPUProfile -cpuprofile=cpu.prof"
    exit 1
fi

echo "========================================="
echo "CPU Profile Analysis"
echo "========================================="
echo ""

# Top 20 functions by CPU time
echo "Top 20 functions by CPU time:"
echo "----------------------------"
go tool pprof -top -nodecount=20 "$PROFILE_FILE"
echo ""

# Generate flame graph (requires graphviz)
if command -v dot &> /dev/null; then
    echo "Generating flame graph..."
    go tool pprof -svg "$PROFILE_FILE" > cpu_flame.svg
    echo "✓ Flame graph saved to: cpu_flame.svg"
    echo ""
fi

# Generate call graph
echo "Generating call graph..."
go tool pprof -pdf "$PROFILE_FILE" > cpu_graph.pdf 2>/dev/null || echo "Note: PDF generation requires graphviz"
echo ""

# Interactive mode instructions
echo "========================================="
echo "For interactive analysis, run:"
echo "  go tool pprof $PROFILE_FILE"
echo ""
echo "Or for web UI:"
echo "  go tool pprof -http=:8080 $PROFILE_FILE"
echo "========================================="
