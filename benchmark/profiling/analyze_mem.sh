#!/bin/bash

# Memory Profile Analyzer
# Generates memory profile and analyzes the data
# Usage: ./analyze_mem.sh [output_file]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

OUTPUT_FILE="${1:-mem_results.txt}"
PROFILE_FILE="mem.prof"

{
    echo "========================================="
    echo "Memory Profile Analysis"
    echo "========================================="
    echo "Date: $(date)"
    echo ""

    # Generate memory profile first
    echo "Generating memory profile..."
    echo "-----------------------------"
    go test -bench=BenchmarkMemoryProfile -benchtime=3s -run=^$ .
    echo ""

    if [ ! -f "$PROFILE_FILE" ]; then
        echo "Error: Profile file '$PROFILE_FILE' was not generated"
        exit 1
    fi

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

    # Check for potential memory leaks
    echo "Checking for potential memory leaks..."
    echo "-------------------------------------"
    go tool pprof -list=. -sample_index=inuse_space "$PROFILE_FILE" 2>/dev/null | head -30 || true
    echo ""

    # Summary
    echo "========================================="
    echo "Analysis Complete"
    echo "========================================="
    echo "Generated files:"
    ls -lh "$PROFILE_FILE" mem_flame.svg 2>/dev/null || true
    echo ""
    echo "For interactive analysis, run:"
    echo "  go tool pprof -http=:8080 $PROFILE_FILE"
    echo ""
    echo "Sample indices available:"
    echo "  -sample_index=alloc_space   # Total allocated memory"
    echo "  -sample_index=alloc_objects # Total allocated objects"
    echo "  -sample_index=inuse_space   # Currently in-use memory"
    echo "  -sample_index=inuse_objects # Currently in-use objects"
    echo "========================================="

} 2>&1 | tee "$OUTPUT_FILE"

echo ""
echo "✓ Results saved to: $OUTPUT_FILE"
