#!/bin/bash

# CPU Profile Analyzer
# Generates CPU profile and analyzes the data
# Usage: ./analyze_cpu.sh [output_file]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

OUTPUT_FILE="${1:-cpu_results.txt}"
PROFILE_FILE="cpu.prof"

{
    echo "========================================="
    echo "CPU Profile Analysis"
    echo "========================================="
    echo "Date: $(date)"
    echo ""

    # Generate CPU profile first
    echo "Generating CPU profile..."
    echo "-------------------------"
    go test -bench=BenchmarkCPUProfile -benchtime=3s -run=^$ .
    echo ""

    if [ ! -f "$PROFILE_FILE" ]; then
        echo "Error: Profile file '$PROFILE_FILE' was not generated"
        exit 1
    fi

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

    # Summary
    echo "========================================="
    echo "Analysis Complete"
    echo "========================================="
    echo "Generated files:"
    ls -lh "$PROFILE_FILE" cpu_flame.svg 2>/dev/null || true
    echo ""
    echo "For interactive analysis, run:"
    echo "  go tool pprof $PROFILE_FILE"
    echo ""
    echo "Or for web UI:"
    echo "  go tool pprof -http=:8080 $PROFILE_FILE"
    echo "========================================="

} 2>&1 | tee "$OUTPUT_FILE"

echo ""
echo "✓ Results saved to: $OUTPUT_FILE"
