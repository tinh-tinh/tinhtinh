#!/bin/bash

# TinhTinh Profiling Script
# Generates profile data and runs analysis in one step
# Usage: ./run_profiling.sh [cpu|mem|goroutine|all] [output_file]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

PROFILE_TYPE="${1:-all}"
OUTPUT_FILE="${2:-profiling_results.txt}"

print_header() {
    echo "========================================"
    echo "$1"
    echo "========================================"
}

run_cpu_profile() {
    print_header "CPU Profiling"
    echo "Generating CPU profile..."
    go test -bench=BenchmarkCPUProfile -benchtime=3s -run=^$ .
    
    if [ -f "cpu.prof" ]; then
        echo ""
        echo "Top 20 functions by CPU time:"
        echo "-----------------------------"
        go tool pprof -top -nodecount=20 cpu.prof
        
        if command -v dot &> /dev/null; then
            echo ""
            echo "Generating CPU flame graph..."
            go tool pprof -svg cpu.prof > cpu_flame.svg 2>/dev/null && \
                echo "✓ Flame graph saved to: cpu_flame.svg"
        fi
    fi
    echo ""
}

run_mem_profile() {
    print_header "Memory Profiling"
    echo "Generating memory profile..."
    go test -bench=BenchmarkMemoryProfile -benchtime=3s -run=^$ .
    
    if [ -f "mem.prof" ]; then
        echo ""
        echo "Top 20 functions by memory allocation:"
        echo "---------------------------------------"
        go tool pprof -top -nodecount=20 mem.prof
        
        echo ""
        echo "Heap profile (inuse_space):"
        echo "---------------------------"
        go tool pprof -top -nodecount=10 -inuse_space mem.prof
    fi
    echo ""
}

run_goroutine_profile() {
    print_header "Goroutine Profiling"
    echo "Generating goroutine profile..."
    go test -bench=BenchmarkGoroutineProfile -benchtime=3s -run=^$ .
    
    if [ -f "goroutine.prof" ]; then
        echo ""
        echo "Goroutine analysis:"
        echo "-------------------"
        go tool pprof -top -nodecount=20 goroutine.prof
    fi
    echo ""
}

run_all_benchmarks() {
    print_header "Running All Profiling Benchmarks"
    echo ""
    go test -bench=. -benchmem -benchtime=2s -run=^$ -v .
    echo ""
}

# Main execution
{
    print_header "TinhTinh Profiling Report"
    echo "Date: $(date)"
    echo "Profile Type: $PROFILE_TYPE"
    echo ""

    case "$PROFILE_TYPE" in
        cpu)
            run_cpu_profile
            ;;
        mem|memory)
            run_mem_profile
            ;;
        goroutine)
            run_goroutine_profile
            ;;
        all)
            run_all_benchmarks
            run_cpu_profile
            run_mem_profile
            run_goroutine_profile
            ;;
        *)
            echo "Unknown profile type: $PROFILE_TYPE"
            echo "Usage: $0 [cpu|mem|goroutine|all] [output_file]"
            exit 1
            ;;
    esac

    print_header "Profiling Complete"
    echo "Generated files:"
    ls -lh *.prof *.svg *.pdf 2>/dev/null || echo "  (profile files in current directory)"
    echo ""
    echo "For interactive analysis, run:"
    echo "  go tool pprof -http=:8080 cpu.prof"
    echo "  go tool pprof -http=:8080 mem.prof"

} 2>&1 | tee "$OUTPUT_FILE"

echo ""
echo "✓ Results saved to: $OUTPUT_FILE"
