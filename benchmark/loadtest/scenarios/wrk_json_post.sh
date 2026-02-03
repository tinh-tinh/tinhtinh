#!/bin/bash

# wrk Load Test - JSON POST
# Tests JSON handling with wrk

set -e

FRAMEWORK="${1:-tinhtinh}"
PORT="${2:-3000}"
DURATION="${3:-30s}"
THREADS="${4:-12}"
CONNECTIONS="${5:-400}"

URL="http://localhost:${PORT}/api/json"

echo "========================================="
echo "wrk - JSON POST Test"
echo "========================================="
echo "Framework: $FRAMEWORK"
echo "URL: $URL"
echo "Duration: $DURATION"
echo "Threads: $THREADS"
echo "Connections: $CONNECTIONS"
echo ""

# Run wrk with Lua script
wrk -t"$THREADS" -c"$CONNECTIONS" -d"$DURATION" --latency \
    -s "../../config/wrk.lua" \
    "$URL" | tee "../results/wrk_${FRAMEWORK}_json_post.txt"

echo ""
echo "Results saved to: ../results/wrk_${FRAMEWORK}_json_post.txt"
