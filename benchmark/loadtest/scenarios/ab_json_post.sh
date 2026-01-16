#!/bin/bash

# JSON POST Load Test using Apache Bench
# Tests JSON request/response handling

set -e

FRAMEWORK="${1:-tinhtinh}"
PORT="${2:-3000}"
REQUESTS="${3:-10000}"
CONCURRENCY="${4:-100}"

URL="http://localhost:${PORT}/api/json"
PAYLOAD='{"name":"test","value":123,"active":true}'

echo "========================================="
echo "Apache Bench - JSON POST Test"
echo "========================================="
echo "Framework: $FRAMEWORK"
echo "URL: $URL"
echo "Requests: $REQUESTS"
echo "Concurrency: $CONCURRENCY"
echo "Payload: $PAYLOAD"
echo ""

# Create temporary file with payload
echo "$PAYLOAD" > /tmp/ab_payload.json

# Run Apache Bench with POST
ab -n "$REQUESTS" -c "$CONCURRENCY" \
   -p /tmp/ab_payload.json \
   -T "application/json" \
   -g "../results/ab_${FRAMEWORK}_json_post.tsv" \
   "$URL" | tee "../results/ab_${FRAMEWORK}_json_post.txt"

# Cleanup
rm -f /tmp/ab_payload.json

echo ""
echo "Results saved to:"
echo "  - ../results/ab_${FRAMEWORK}_json_post.txt"
echo "  - ../results/ab_${FRAMEWORK}_json_post.tsv"
