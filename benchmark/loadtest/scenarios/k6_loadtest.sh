#!/bin/bash

# k6 Load Test - Multiple Scenarios
# Modern load testing tool with JavaScript scenarios

set -e

FRAMEWORK="${1:-tinhtinh}"
PORT="${2:-3000}"

export BASE_URL="http://localhost:${PORT}"

echo "========================================="
echo "k6 - Load Test Scenarios"
echo "========================================="
echo "Framework: $FRAMEWORK"
echo "Base URL: $BASE_URL"
echo ""

# Run k6 with the configuration file
k6 run --out json="../results/k6_${FRAMEWORK}_results.json" \
       ../../config/k6.js | tee "../results/k6_${FRAMEWORK}_output.txt"

echo ""
echo "Results saved to:"
echo "  - ../results/k6_${FRAMEWORK}_results.json"
echo "  - ../results/k6_${FRAMEWORK}_output.txt"
