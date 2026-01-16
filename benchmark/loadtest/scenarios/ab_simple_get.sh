#!/bin/bash

# Simple GET Load Test using Apache Bench
# Tests simple GET requests to measure baseline performance

set -e

FRAMEWORK="${1:-tinhtinh}"
PORT="${2:-3000}"
REQUESTS="${3:-10000}"
CONCURRENCY="${4:-100}"

URL="http://localhost:${PORT}/api/"

echo "========================================="
echo "Apache Bench - Simple GET Test"
echo "========================================="
echo "Framework: $FRAMEWORK"
echo "URL: $URL"
echo "Requests: $REQUESTS"
echo "Concurrency: $CONCURRENCY"
echo ""

# Run Apache Bench
ab -n "$REQUESTS" -c "$CONCURRENCY" -g "../results/ab_${FRAMEWORK}_simple_get.tsv" "$URL" | tee "../results/ab_${FRAMEWORK}_simple_get.txt"

echo ""
echo "Results saved to:"
echo "  - ../results/ab_${FRAMEWORK}_simple_get.txt"
echo "  - ../results/ab_${FRAMEWORK}_simple_get.tsv"
