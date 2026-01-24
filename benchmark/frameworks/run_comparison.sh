#!/bin/bash

# Framework Comparison Benchmark Runner
# This script runs benchmarks for all frameworks and generates a comparison report

set -e

echo "========================================="
echo "Framework Comparison Benchmark"
echo "========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
BENCH_TIME="${BENCH_TIME:-5s}"
BENCH_COUNT="${BENCH_COUNT:-5}"
OUTPUT_DIR="./results"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Create output directory
mkdir -p "$OUTPUT_DIR"

echo -e "${GREEN}Configuration:${NC}"
echo "  Benchmark time: $BENCH_TIME"
echo "  Iterations: $BENCH_COUNT"
echo "  Output directory: $OUTPUT_DIR"
echo ""

# Function to run benchmarks
run_benchmark() {
    local framework=$1
    local output_file="$OUTPUT_DIR/${framework}_${TIMESTAMP}.txt"
    
    echo -e "${YELLOW}Running $framework benchmarks...${NC}"
    go test -bench="Benchmark${framework}" \
        -benchmem \
        -benchtime="$BENCH_TIME" \
        -count="$BENCH_COUNT" \
        -run=^$ \
        | tee "$output_file"
    
    echo -e "${GREEN}✓ $framework benchmarks complete${NC}"
    echo ""
}

# Run benchmarks for each framework
echo "Starting benchmark runs..."
echo ""

run_benchmark "TinhTinh"
run_benchmark "Gin"
run_benchmark "Echo"
run_benchmark "Fiber"
run_benchmark "Chi"

# Generate comparison report
echo -e "${YELLOW}Generating comparison report...${NC}"

REPORT_FILE="$OUTPUT_DIR/comparison_${TIMESTAMP}.md"

cat > "$REPORT_FILE" << 'EOF'
# Framework Comparison Report

Generated: $(date)

## Benchmark Results

### Simple GET Request

| Framework | ns/op | B/op | allocs/op |
|-----------|-------|------|-----------|
EOF

# Parse results and add to report
# This is a simplified version - in production you'd use benchstat or similar
for framework in TinhTinh Gin Echo Fiber Chi; do
    result_file="$OUTPUT_DIR/${framework}_${TIMESTAMP}.txt"
    if [ -f "$result_file" ]; then
        # Extract SimpleGET benchmark results (simplified parsing)
        grep "Benchmark${framework}_SimpleGET" "$result_file" | tail -1 | \
            awk -v fw="$framework" '{printf "| %s | %s | %s | %s |\n", fw, $3, $5, $7}' >> "$REPORT_FILE"
    fi
done

cat >> "$REPORT_FILE" << 'EOF'

### JSON Response

| Framework | ns/op | B/op | allocs/op |
|-----------|-------|------|-----------|
EOF

for framework in TinhTinh Gin Echo Fiber Chi; do
    result_file="$OUTPUT_DIR/${framework}_${TIMESTAMP}.txt"
    if [ -f "$result_file" ]; then
        grep "Benchmark${framework}_JSONResponse" "$result_file" | tail -1 | \
            awk -v fw="$framework" '{printf "| %s | %s | %s | %s |\n", fw, $3, $5, $7}' >> "$REPORT_FILE"
    fi
done

cat >> "$REPORT_FILE" << 'EOF'

### Path Parameter Parsing

| Framework | ns/op | B/op | allocs/op |
|-----------|-------|------|-----------|
EOF

for framework in TinhTinh Gin Echo Fiber Chi; do
    result_file="$OUTPUT_DIR/${framework}_${TIMESTAMP}.txt"
    if [ -f "$result_file" ]; then
        grep "Benchmark${framework}_PathParam" "$result_file" | tail -1 | \
            awk -v fw="$framework" '{printf "| %s | %s | %s | %s |\n", fw, $3, $5, $7}' >> "$REPORT_FILE"
    fi
done

cat >> "$REPORT_FILE" << 'EOF'

### Parallel Requests

| Framework | ns/op | B/op | allocs/op |
|-----------|-------|------|-----------|
EOF

for framework in TinhTinh Gin Echo Fiber Chi; do
    result_file="$OUTPUT_DIR/${framework}_${TIMESTAMP}.txt"
    if [ -f "$result_file" ]; then
        grep "Benchmark${framework}_ParallelRequests" "$result_file" | tail -1 | \
            awk -v fw="$framework" '{printf "| %s | %s | %s | %s |\n", fw, $3, $5, $7}' >> "$REPORT_FILE"
    fi
done

cat >> "$REPORT_FILE" << EOF

## Notes

- Lower values are better for all metrics
- ns/op: Nanoseconds per operation
- B/op: Bytes allocated per operation
- allocs/op: Number of allocations per operation

## Raw Results

Raw benchmark results are available in:
$(ls -1 $OUTPUT_DIR/*_${TIMESTAMP}.txt | sed 's/^/- /')

EOF

echo -e "${GREEN}✓ Comparison report generated: $REPORT_FILE${NC}"
echo ""

# Display summary
echo "========================================="
echo "Benchmark Summary"
echo "========================================="
cat "$REPORT_FILE"

echo ""
echo -e "${GREEN}All benchmarks complete!${NC}"
echo "Results saved to: $OUTPUT_DIR"
