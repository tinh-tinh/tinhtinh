#!/bin/bash

# Load Test Runner
# Runs various load testing tools against the Tinh Tinh application

set -e

echo "========================================="
echo "Load Testing Suite"
echo "========================================="
echo ""

# Configuration
PORT="${PORT:-3000}"
BASE_URL="http://localhost:${PORT}"
RESULTS_DIR="./results"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Create results directory
mkdir -p "$RESULTS_DIR"

# Check if server is running
echo -e "${YELLOW}Checking if server is running on port ${PORT}...${NC}"
if ! curl -s "${BASE_URL}/" > /dev/null 2>&1; then
    echo -e "${RED}Error: Server not running on port ${PORT}${NC}"
    echo "Please start your application first:"
    echo "  cd apps/tinhtinh_app && go run main.go"
    exit 1
fi
echo -e "${GREEN}✓ Server is running${NC}"
echo ""

# Function to run Apache Bench
run_ab() {
    echo -e "${YELLOW}Running Apache Bench tests...${NC}"
    
    if ! command -v ab &> /dev/null; then
        echo -e "${RED}Apache Bench not installed, skipping...${NC}"
        return
    fi
    
    local output_file="${RESULTS_DIR}/ab_${TIMESTAMP}.txt"
    
    # Simple GET test
    echo "  - Simple GET (10,000 requests, concurrency 100)"
    ab -n 10000 -c 100 -g "${RESULTS_DIR}/ab_gnuplot.tsv" "${BASE_URL}/" > "$output_file" 2>&1
    
    echo -e "${GREEN}✓ Apache Bench tests complete${NC}"
    echo "  Results: $output_file"
    echo ""
}

# Function to run wrk
run_wrk() {
    echo -e "${YELLOW}Running wrk tests...${NC}"
    
    if ! command -v wrk &> /dev/null; then
        echo -e "${RED}wrk not installed, skipping...${NC}"
        return
    fi
    
    local output_file="${RESULTS_DIR}/wrk_${TIMESTAMP}.txt"
    
    # Run wrk with Lua script
    echo "  - Running 30s test with 12 threads, 400 connections"
    wrk -t12 -c400 -d30s --latency -s ../../config/wrk.lua "${BASE_URL}/" | tee "$output_file"
    
    echo -e "${GREEN}✓ wrk tests complete${NC}"
    echo "  Results: $output_file"
    echo ""
}

# Function to run k6
run_k6() {
    echo -e "${YELLOW}Running k6 tests...${NC}"
    
    if ! command -v k6 &> /dev/null; then
        echo -e "${RED}k6 not installed, skipping...${NC}"
        return
    fi
    
    local output_file="${RESULTS_DIR}/k6_${TIMESTAMP}.json"
    
    # Run k6 scenarios
    echo "  - Running k6 load test scenarios"
    BASE_URL="${BASE_URL}" k6 run --out json="${output_file}" ../../config/k6.js
    
    echo -e "${GREEN}✓ k6 tests complete${NC}"
    echo "  Results: $output_file"
    echo ""
}

# Run all tests
echo "Starting load tests against ${BASE_URL}"
echo ""

run_ab
run_wrk
run_k6

# Generate summary
echo "========================================="
echo "Load Test Summary"
echo "========================================="
echo ""
echo "All load tests complete!"
echo "Results saved to: $RESULTS_DIR"
echo ""
echo "Files generated:"
ls -lh "$RESULTS_DIR"/*_${TIMESTAMP}* 2>/dev/null || echo "No results generated"
echo ""
