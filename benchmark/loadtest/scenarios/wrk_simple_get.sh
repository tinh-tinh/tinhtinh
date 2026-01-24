#!/bin/bash

# wrk Load Test - Simple GET
# High-performance HTTP benchmarking tool

set -e

FRAMEWORK="${1:-tinhtinh}"
PORT="${2:-3000}"
DURATION="${3:-30s}"
THREADS="${4:-12}"
CONNECTIONS="${5:-400}"

URL="http://localhost:${PORT}/api/"

echo "========================================="
echo "wrk - Simple GET Test"
echo "========================================="
echo "Framework: $FRAMEWORK"
echo "URL: $URL"
echo "Duration: $DURATION"
echo "Threads: $THREADS"
echo "Connections: $CONNECTIONS"
echo ""

# Run wrk
wrk -t"$THREADS" -c"$CONNECTIONS" -d"$DURATION" --latency "$URL" | tee "../results/wrk_${FRAMEWORK}_simple_get.txt"

echo ""
echo "Results saved to: ../results/wrk_${FRAMEWORK}_simple_get.txt"
