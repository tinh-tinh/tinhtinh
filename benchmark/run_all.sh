#!/bin/bash

# TinhTinh Complete Benchmark Suite
# Runs all benchmarks and generates a comprehensive performance report
# Usage: ./run_all.sh [output_file]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

OUTPUT_FILE="${1:-PERFORMANCE_REPORT.txt}"
TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")

print_header() {
    echo ""
    echo "╔══════════════════════════════════════════════════════════════════╗"
    echo "║  $1"
    echo "╚══════════════════════════════════════════════════════════════════╝"
    echo ""
}

print_section() {
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "  $1"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
}

{
    print_header "TINHTINH FRAMEWORK PERFORMANCE REPORT"
    echo "Generated: $(date)"
    echo "System: $(uname -a)"
    echo "Go Version: $(go version)"
    echo "CPU: $(grep -m1 'model name' /proc/cpuinfo 2>/dev/null | cut -d: -f2 | xargs || echo 'Unknown')"
    echo "Memory: $(free -h 2>/dev/null | grep Mem | awk '{print $2}' || echo 'Unknown')"
    echo ""

    # ═══════════════════════════════════════════════════════════════════════
    # LATENCY BENCHMARKS
    # ═══════════════════════════════════════════════════════════════════════
    print_section "1. LATENCY BENCHMARKS"
    echo "Testing request latency with percentile distribution (p50, p95, p99)"
    echo ""
    cd "$SCRIPT_DIR/latency"
    go test -bench=. -benchmem -benchtime=2s -run=^$ . 2>&1 || echo "Latency tests completed with warnings"
    cd "$SCRIPT_DIR"

    # ═══════════════════════════════════════════════════════════════════════
    # CONCURRENT BENCHMARKS
    # ═══════════════════════════════════════════════════════════════════════
    print_section "2. CONCURRENT REQUEST BENCHMARKS"
    echo "Testing with various concurrency levels (10, 100, 1000, 10000 goroutines)"
    echo ""
    cd "$SCRIPT_DIR/concurrent"
    go test -bench=BenchmarkConcurrent -benchmem -benchtime=1s -run=^$ . 2>&1 || echo "Concurrent tests completed with warnings"
    cd "$SCRIPT_DIR"

    # ═══════════════════════════════════════════════════════════════════════
    # THROUGHPUT BENCHMARKS
    # ═══════════════════════════════════════════════════════════════════════
    print_section "3. THROUGHPUT BENCHMARKS"
    echo "Measuring requests per second (RPS)"
    echo ""
    cd "$SCRIPT_DIR/concurrent"
    go test -bench=BenchmarkThroughput -benchmem -benchtime=2s -run=^$ . 2>&1 || echo "Throughput tests completed with warnings"
    cd "$SCRIPT_DIR"

    # ═══════════════════════════════════════════════════════════════════════
    # PROFILING BENCHMARKS
    # ═══════════════════════════════════════════════════════════════════════
    print_section "4. PROFILING BENCHMARKS"
    echo "CPU, Memory, and Goroutine profiling"
    echo ""
    cd "$SCRIPT_DIR/profiling"
    go test -bench=. -benchmem -benchtime=1s -run=^$ . 2>&1 || echo "Profiling tests completed with warnings"
    cd "$SCRIPT_DIR"

    # ═══════════════════════════════════════════════════════════════════════
    # SUMMARY
    # ═══════════════════════════════════════════════════════════════════════
    print_header "PERFORMANCE SUMMARY"
    
    echo "Benchmark Categories:"
    echo "  ✓ Latency       - Request response time percentiles"
    echo "  ✓ Concurrent    - Multi-goroutine performance scaling"
    echo "  ✓ Throughput    - Requests per second capacity"
    echo "  ✓ Profiling     - CPU, memory, goroutine analysis"
    echo ""
    
    echo "Key Metrics Explained:"
    echo "  • ns/op      - Nanoseconds per operation (lower is better)"
    echo "  • B/op       - Bytes allocated per operation (lower is better)"
    echo "  • allocs/op  - Number of allocations per operation (lower is better)"
    echo "  • req/s      - Requests per second (higher is better)"
    echo ""
    
    echo "Performance Targets:"
    echo "  • Latency p99 < 50ms for simple requests"
    echo "  • Throughput > 10,000 RPS on modern hardware"
    echo "  • Memory < 5KB per request"
    echo "  • Allocations < 10 per request"
    echo ""
    
    print_header "REPORT COMPLETE"
    echo "Completed at: $(date)"
    echo "Results saved to: $OUTPUT_FILE"

} 2>&1 | tee "$OUTPUT_FILE"

echo ""
echo "╔══════════════════════════════════════════════════════════════════╗"
echo "║  ✓ Performance report saved to: $OUTPUT_FILE"
echo "╚══════════════════════════════════════════════════════════════════╝"
