# TinhTinh Performance Summary

## Quick Links

| Benchmark | Script | Description |
|-----------|--------|-------------|
| **All** | `./run_all.sh` | Complete performance report |
| **Latency** | `./latency/run_latency.sh` | Response time percentiles |
| **Concurrent** | `./concurrent/run_concurrent.sh` | Multi-goroutine scaling |
| **CPU Profile** | `./profiling/analyze_cpu.sh` | CPU usage analysis |
| **Memory Profile** | `./profiling/analyze_mem.sh` | Memory allocation analysis |
| **Goroutine Profile** | `./profiling/analyze_goroutine.sh` | Goroutine analysis |
| **Profiling All** | `./profiling/run_profiling.sh` | All profiling combined |
| **Framework Comparison** | `./frameworks/run_comparison.sh` | Compare with Gin, Echo, Fiber, Chi |

## Benchmark Categories

### 1. Latency Benchmarks (`benchmark/latency/`)
Measures request response time with percentile distribution:
- **p50 (median)**: Typical user experience
- **p95**: Most users' experience  
- **p99**: Worst-case for most users
- **p99.9**: Absolute worst-case

### 2. Concurrent Benchmarks (`benchmark/concurrent/`)
Tests performance under concurrent load:
- 10, 100, 1,000, 10,000 concurrent goroutines
- Context pooling efficiency
- Sustained load testing

### 3. Throughput Benchmarks (`benchmark/concurrent/`)
Measures requests per second (RPS):
- Simple GET throughput
- JSON response throughput
- Parallel throughput
- Keep-alive performance
- Burst load handling

### 4. Profiling (`benchmark/profiling/`)
Detailed performance analysis:
- **CPU**: Identify slow functions
- **Memory**: Find allocation hotspots
- **Goroutine**: Detect goroutine leaks

### 5. Framework Comparisons (`benchmark/frameworks/`)
Compare TinhTinh with other Go web frameworks:
- Gin
- Echo
- Fiber
- Chi

## Performance Targets

| Metric | Target | Description |
|--------|--------|-------------|
| **Latency p99** | < 50ms | 99th percentile response time |
| **Throughput** | > 10,000 RPS | Requests per second |
| **Memory** | < 5KB/req | Memory per request |
| **Allocations** | < 10/req | Allocations per request |

## Key Metrics Explained

- **ns/op**: Nanoseconds per operation (lower is better)
- **B/op**: Bytes allocated per operation (lower is better)
- **allocs/op**: Number of allocations (lower is better)
- **req/s**: Requests per second (higher is better)

## Quick Start

```bash
# Run complete benchmark suite
cd benchmark
./run_all.sh

# Run specific benchmarks
./latency/run_latency.sh
./concurrent/run_concurrent.sh
./profiling/run_profiling.sh all

# View interactive CPU profile
go tool pprof -http=:8080 profiling/cpu.prof
```

## Generated Files

After running benchmarks, you'll find:
- `*.txt` - Text reports in each directory
- `*.prof` - Profile files for pprof analysis
- `*.svg` - Flame graphs and call graphs
