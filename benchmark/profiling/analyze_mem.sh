#!/bin/bash

# Memory Profile Analyzer
# Analyzes memory profiling data and generates reports

set -e

PROFILE_FILE="${1:-mem.prof}"

if [ ! -f "$PROFILE_FILE" ]; then
    echo "Error: Profile file '$PROFILE_FILE' not found"
    echo "Usage: $0 [profile_file]"
    echo "Run benchmarks first: go test -bench=BenchmarkMemoryProfile -memprofile=mem.prof"
    exit 1
fi

echo "========================================="
echo "Memory Profile Analysis"
echo "========================================="
echo ""

# Top 20 functions by memory allocation
echo "Top 20 functions by allocated memory:"
echo "------------------------------------"
go tool pprof -top -nodecount=20 "$PROFILE_FILE"
echo ""

# Top allocations (inuse_space)
echo "Top allocations (in-use space):"
echo "-------------------------------"
go tool pprof -top -nodecount=20 -sample_index=inuse_space "$PROFILE_FILE"
echo ""

# Top allocations (inuse_objects)
echo "Top allocations (in-use objects):"
echo "---------------------------------"
go tool pprof -top -nodecount=20 -sample_index=inuse_objects "$PROFILE_FILE"
echo ""

# Generate memory flame graph
if command -v dot &> /dev/null; then
    echo "Generating memory flame graph..."
    go tool pprof -svg "$PROFILE_FILE" > mem_flame.svg
    echo "✓ Memory flame graph saved to: mem_flame.svg"
    echo ""
fi

# Generate call graph
echo "Generating call graph..."
go tool pprof -pdf "$PROFILE_FILE" > mem_graph.pdf 2>/dev/null || echo "Note: PDF generation requires graphviz"
echo ""

# Check for potential memory leaks
echo "Checking for potential memory leaks..."
echo "-------------------------------------"
go tool pprof -list=. -sample_index=inuse_space "$PROFILE_FILE" | head -30
echo ""

# Interactive mode instructions
echo "========================================="
echo "For interactive analysis, run:"
echo "  go tool pprof $PROFILE_FILE"
echo ""
echo "Or for web UI:"
echo "  go tool pprof -http=:8080 $PROFILE_FILE"
echo ""
echo "Sample indices available:"
echo "  -sample_index=alloc_space   # Total allocated memory"
echo "  -sample_index=alloc_objects # Total allocated objects"
echo "  -sample_index=inuse_space   # Currently in-use memory"
echo "  -sample_index=inuse_objects # Currently in-use objects"
echo "========================================="
