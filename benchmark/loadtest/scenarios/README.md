# Load Test Scenarios

This directory contains load test scenario scripts for different tools.

## Available Scenarios

### Apache Bench (ab)
- `ab_simple_get.sh` - Simple GET request test
- `ab_json_post.sh` - JSON POST request test

### wrk
- `wrk_simple_get.sh` - Simple GET with latency distribution
- `wrk_json_post.sh` - JSON POST with custom Lua script

### k6
- `k6_loadtest.sh` - Multiple scenarios (constant, ramping, spike, stress, soak)

## Usage

### Apache Bench Scenarios

**Simple GET:**
```bash
cd scenarios
./ab_simple_get.sh [framework] [port] [requests] [concurrency]

# Examples:
./ab_simple_get.sh tinhtinh 3000 10000 100
./ab_simple_get.sh gin 3001 10000 100
```

**JSON POST:**
```bash
./ab_json_post.sh [framework] [port] [requests] [concurrency]

# Example:
./ab_json_post.sh tinhtinh 3000 5000 50
```

### wrk Scenarios

**Simple GET:**
```bash
./wrk_simple_get.sh [framework] [port] [duration] [threads] [connections]

# Examples:
./wrk_simple_get.sh tinhtinh 3000 30s 12 400
./wrk_simple_get.sh gin 3001 60s 12 800
```

**JSON POST:**
```bash
./wrk_json_post.sh [framework] [port] [duration] [threads] [connections]

# Example:
./wrk_json_post.sh tinhtinh 3000 30s 12 400
```

### k6 Scenarios

```bash
./k6_loadtest.sh [framework] [port]

# Examples:
./k6_loadtest.sh tinhtinh 3000
./k6_loadtest.sh gin 3001
```

## Running Complete Benchmark Suite

To benchmark all frameworks:

```bash
# Start all applications first (in separate terminals)
cd ../apps
go run tinhtinh_app.go  # Terminal 1
go run gin_app.go       # Terminal 2
go run echo_app.go      # Terminal 3
go run fiber_app.go     # Terminal 4
go run chi_app.go       # Terminal 5

# Then run benchmarks
cd ../scenarios

# Apache Bench
for fw in tinhtinh gin echo fiber chi; do
  port=$((3000 + $(echo "tinhtinh gin echo fiber chi" | tr ' ' '\n' | grep -n $fw | cut -d: -f1) - 1))
  ./ab_simple_get.sh $fw $port
done

# wrk
for fw in tinhtinh gin echo fiber chi; do
  port=$((3000 + $(echo "tinhtinh gin echo fiber chi" | tr ' ' '\n' | grep -n $fw | cut -d: -f1) - 1))
  ./wrk_simple_get.sh $fw $port
done

# k6
for fw in tinhtinh gin echo fiber chi; do
  port=$((3000 + $(echo "tinhtinh gin echo fiber chi" | tr ' ' '\n' | grep -n $fw | cut -d: -f1) - 1))
  ./k6_loadtest.sh $fw $port
done
```

## Results

All results are saved to `../results/` directory with the following naming convention:
- `ab_[framework]_[test].txt` - Apache Bench results
- `ab_[framework]_[test].tsv` - Apache Bench gnuplot data
- `wrk_[framework]_[test].txt` - wrk results
- `k6_[framework]_results.json` - k6 JSON results
- `k6_[framework]_output.txt` - k6 console output

## Prerequisites

Make sure you have the following tools installed:
- Apache Bench: `sudo apt-get install apache2-utils`
- wrk: `sudo apt-get install wrk` or build from source
- k6: `sudo gpg -k && sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69 && echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list && sudo apt-get update && sudo apt-get install k6`
