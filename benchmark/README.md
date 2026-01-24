# Benchmark Infrastructure

Comprehensive benchmarking suite for the Tinh Tinh framework, including framework comparisons, load testing, profiling, and performance regression detection.

## Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [Framework Comparisons](#framework-comparisons)
- [Load Testing](#load-testing)
- [Profiling](#profiling)
- [Latency Analysis](#latency-analysis)
- [Concurrent Request Testing](#concurrent-request-testing)
- [CI Integration](#ci-integration)
- [Interpreting Results](#interpreting-results)
- [Contributing](#contributing)

## Overview

This benchmark suite provides:

- **Framework Comparisons**: Performance benchmarks comparing Tinh Tinh with Gin, Echo, Fiber, and Chi
- **Load Testing**: Scenarios using Apache Bench, wrk, and k6
- **Profiling**: CPU, memory, and goroutine profiling via pprof
- **Latency Measurements**: p50, p95, p99 percentile tracking
- **Concurrent Request Tests**: Testing at various concurrency levels (10, 100, 1000, 10000)
- **CI Integration**: Automated performance regression detection
- **Weekly Reports**: Automated benchmark reports with visualizations

## Quick Start

### Prerequisites

Install required tools:

```bash
# Apache Bench (usually comes with Apache)
sudo apt-get install apache2-utils  # Ubuntu/Debian
brew install httpd                   # macOS

# wrk
sudo apt-get install wrk            # Ubuntu/Debian
brew install wrk                     # macOS

# k6
sudo gpg -k
sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6

# Or use the install script
brew install k6  # macOS
```

### Running All Benchmarks

```bash
# Run framework comparison benchmarks
cd benchmark/frameworks
./run_comparison.sh

# Run load tests
cd ../loadtest
./run_loadtest.sh

# Run profiling tests
cd ../profiling
go test -bench=. -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof
./analyze_cpu.sh
./analyze_mem.sh

# Run latency tests
cd ../latency
go test -bench=. -benchmem

# Run concurrent request tests
cd ../concurrent
go test -bench=. -benchmem
```

## Framework Comparisons

### Running Comparisons

```bash
cd benchmark/frameworks
./run_comparison.sh
```

This will run identical benchmarks across all frameworks and generate a comparison report.

### Available Benchmarks

- **Simple GET**: Basic route handling
- **JSON Response**: JSON serialization
- **JSON Request/Response**: Full JSON round-trip
- **Path Parameters**: Dynamic route parsing
- **Query Parameters**: Query string parsing
- **Middleware Chain**: Multiple middleware execution

### Example Output

```
Framework Comparison Results
============================

Simple GET Request (10000 iterations):
Framework    | ns/op    | B/op  | allocs/op
-------------|----------|-------|----------
Tinh Tinh    | 12345    | 1024  | 5
Gin          | 13456    | 1152  | 6
Echo         | 11234    | 896   | 4
Fiber        | 10123    | 768   | 3
Chi          | 14567    | 1280  | 7
```

## Load Testing

### Apache Bench (ab)

Simple HTTP load testing:

```bash
cd benchmark/loadtest/scenarios
./simple_get.sh
```

Configuration in `benchmark/config/ab.conf`:
- 100,000 requests
- Concurrency levels: 10, 50, 100, 500, 1000

### wrk

Advanced HTTP benchmarking:

```bash
cd benchmark/loadtest/scenarios
wrk -t12 -c400 -d30s --latency -s ../../config/wrk.lua http://localhost:3000/
```

Features:
- Custom Lua scripting
- Latency distribution
- Request customization

### k6

Modern load testing:

```bash
cd benchmark/loadtest/scenarios
k6 run ../../config/k6.js
```

Scenarios:
- Ramping VUs (gradual load increase)
- Constant load
- Spike tests
- Stress tests
- Soak tests (sustained load)

## Profiling

### CPU Profiling

```bash
cd benchmark/profiling
go test -bench=BenchmarkCPU -cpuprofile=cpu.prof
./analyze_cpu.sh
```

This generates:
- CPU flame graph
- Top functions by CPU time
- Call graph visualization

### Memory Profiling

```bash
go test -bench=BenchmarkMemory -memprofile=mem.prof
./analyze_mem.sh
```

This shows:
- Memory allocations
- Allocation hotspots
- Potential memory leaks

### Goroutine Profiling

```bash
go test -bench=BenchmarkGoroutine
./analyze_goroutine.sh
```

This tracks:
- Goroutine creation
- Goroutine lifecycle
- Potential goroutine leaks

### Analyzing Profiles

```bash
# Interactive analysis
go tool pprof cpu.prof

# Generate flame graph
go tool pprof -http=:8080 cpu.prof

# Top 10 functions
go tool pprof -top cpu.prof

# List specific function
go tool pprof -list=FunctionName cpu.prof
```

## Latency Analysis

### Running Latency Tests

```bash
cd benchmark/latency
go test -bench=. -benchmem -v
```

### Percentile Metrics

The latency tests track:
- **p50 (median)**: 50% of requests complete within this time
- **p95**: 95% of requests complete within this time
- **p99**: 99% of requests complete within this time
- **p99.9**: 99.9% of requests complete within this time

### Example Output

```
Latency Distribution
====================
Min:    1.2ms
p50:    5.3ms
p95:    12.7ms
p99:    25.4ms
p99.9:  45.2ms
Max:    102.3ms
Mean:   6.8ms
```

### Interpreting Latency

- **p50**: Typical user experience
- **p95**: Most users' experience
- **p99**: Worst-case for most users
- **p99.9**: Absolute worst-case scenarios

## Concurrent Request Testing

### Running Concurrent Tests

```bash
cd benchmark/concurrent
go test -bench=. -benchmem
```

### Concurrency Levels

Tests run at:
- 10 concurrent requests
- 100 concurrent requests
- 1,000 concurrent requests
- 10,000 concurrent requests

### Throughput Measurement

```bash
go test -bench=BenchmarkThroughput -benchmem
```

Measures:
- Requests per second (RPS)
- Connection pooling efficiency
- Keep-alive performance

## CI Integration

### GitHub Actions Workflow

The benchmark workflow runs:
- On every push to main
- On pull requests
- Weekly (scheduled)

### Performance Regression Detection

The CI automatically:
1. Runs all benchmarks
2. Compares with baseline results
3. Detects regressions >5%
4. Posts results as PR comment
5. Fails PR if critical regressions found

### Viewing Results

Results are available:
- As PR comments
- In GitHub Actions artifacts
- In the performance dashboard

### Baseline Management

Update baseline after approved performance changes:

```bash
cd benchmark/scripts
go run update_baseline.go
```

## Interpreting Results

### Benchmark Metrics

- **ns/op**: Nanoseconds per operation (lower is better)
- **B/op**: Bytes allocated per operation (lower is better)
- **allocs/op**: Number of allocations per operation (lower is better)

### What's Good Performance?

For web frameworks:
- **Latency**: p99 < 50ms for simple requests
- **Throughput**: >10,000 RPS on modern hardware
- **Memory**: <5KB per request
- **Allocations**: <10 per request

### Performance Regression Thresholds

- **Minor**: 5-10% degradation (warning)
- **Moderate**: 10-20% degradation (review required)
- **Critical**: >20% degradation (block merge)

## Contributing

### Adding New Benchmarks

1. Create benchmark function in appropriate test file:

```go
func BenchmarkNewFeature(b *testing.B) {
    // Setup
    app := setupApp()
    
    b.ResetTimer()
    b.ReportAllocs()
    
    for i := 0; i < b.N; i++ {
        // Benchmark code
    }
}
```

2. Add to comparison script if comparing frameworks

3. Update documentation

### Best Practices

1. **Isolate benchmarks**: Each benchmark should test one thing
2. **Use b.ResetTimer()**: Reset after setup code
3. **Use b.ReportAllocs()**: Track memory allocations
4. **Avoid I/O**: Mock external dependencies
5. **Run multiple times**: Use `-count=10` for statistical significance
6. **Disable CPU scaling**: For consistent results

### Running Benchmarks Properly

```bash
# Disable CPU frequency scaling (Linux)
sudo cpupower frequency-set --governor performance

# Run with multiple iterations
go test -bench=. -benchmem -count=10 -benchtime=5s

# Re-enable CPU scaling
sudo cpupower frequency-set --governor powersave
```

## Automated Reporting

Weekly reports are automatically generated and include:
- Framework comparison trends
- Performance improvements/regressions
- Latency percentile charts
- Memory usage graphs
- Top performance issues

Reports are posted to GitHub Discussions every Monday.

## Troubleshooting

### High Variance in Results

- Disable CPU frequency scaling
- Close other applications
- Run multiple iterations
- Use `-benchtime=10s` for longer runs

### Out of Memory Errors

- Reduce concurrency levels
- Increase available memory
- Check for memory leaks in code

### Connection Refused Errors

- Ensure test app is running
- Check port availability
- Verify firewall settings

## License

Same as Tinh Tinh framework - see [LICENSE](../LICENSE)

## Support

- Issues: [GitHub Issues](https://github.com/tinh-tinh/tinhtinh/issues)
- Discussions: [GitHub Discussions](https://github.com/tinh-tinh/tinhtinh/discussions)
- Documentation: [Main README](../README.md)
