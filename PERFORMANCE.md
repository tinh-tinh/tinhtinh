# Performance Factors for Backend Frameworks

This document lists the key performance factors for backend frameworks and explains how the Tinh Tinh framework addresses each of them.

## Table of Contents

- [1. Memory Management](#1-memory-management)
- [2. Concurrency and Parallelism](#2-concurrency-and-parallelism)
- [3. Request/Response Handling](#3-requestresponse-handling)
- [4. Routing Efficiency](#4-routing-efficiency)
- [5. Dependency Injection and Provider Scoping](#5-dependency-injection-and-provider-scoping)
- [6. Logging and Observability](#6-logging-and-observability)
- [7. Caching](#7-caching)
- [8. Serialization/Deserialization](#8-serializationdeserialization)
- [9. Middleware Pipeline](#9-middleware-pipeline)
- [10. Resource Cleanup](#10-resource-cleanup)

---

## 1. Memory Management

### Performance Factor
Efficient memory management is critical to avoid excessive garbage collection pauses and memory leaks. Frameworks should minimize allocations and reuse objects where possible.

### How Tinh Tinh Addresses This

**Object Pooling with sync.Pool**

Tinh Tinh uses `sync.Pool` to reuse context objects, reducing memory allocations per request:

```go
// core/app.go
type App struct {
    pool sync.Pool
    // ...
}

func CreateFactory(module ModuleParam, opt ...AppOptions) *App {
    app := &App{
        // ...
    }
    app.pool = sync.Pool{
        New: func() any {
            return NewCtx(app)
        },
    }
    // ...
}
```

**Context Reuse in Request Handling**

```go
// core/ctx.go
func ParseCtx(app *App, router *Router) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := app.pool.Get().(*DefaultCtx)
        defer app.pool.Put(ctx)
        // handle request...
    })
}
```

**Garbage Collection Optimization**

The framework explicitly triggers garbage collection after route registration to clean up temporary data:

```go
// core/router.go
func (app *App) free() {
    app.Module.free()
    runtime.GC()
}
```

---

## 2. Concurrency and Parallelism

### Performance Factor
Backend frameworks must handle concurrent requests efficiently, using proper synchronization to prevent race conditions while maximizing throughput.

### How Tinh Tinh Addresses This

**Mutex-Protected Shared State**

The framework uses `sync.RWMutex` for thread-safe access to shared data:

```go
// common/memory/memory.go
type Store struct {
    ttl  time.Duration
    data map[string]item
    sync.RWMutex
}

func (m *Store) Get(key string) interface{} {
    m.RLock()
    v, ok := m.data[key]
    m.RUnlock()
    // ...
}

func (m *Store) Set(key string, val interface{}, ttl ...time.Duration) {
    // ...
    m.Lock()
    m.data[key] = i
    m.Unlock()
}
```

**Asynchronous Logging**

The logger uses buffered channels for non-blocking log writes:

```go
// common/logger/logger.go
type Logger struct {
    // ...
    logCh      chan *logEntry
    bufferSize int
}

func (log *Logger) write(level Level, msg string, meta ...Metadata) {
    entry := &logEntry{
        level: level,
        msg:   msg,
        meta:  meta,
        time:  time.Now(),
    }

    // Non-blocking send
    select {
    case log.logCh <- entry:
    default:
        // Channel full, log warning to stderr
        fmt.Fprintf(os.Stderr, "[WARN] Log buffer full...")
    }
}
```

**Goroutine-Based Microservices**

Microservices are started as goroutines for parallel execution:

```go
// core/app.go
func (app *App) StartAllMicroservices() {
    for _, svc := range app.Services {
        go svc.Listen()
    }
}
```

---

## 3. Request/Response Handling

### Performance Factor
Efficient handling of HTTP requests and responses, including proper streaming, buffering, and status code management.

### How Tinh Tinh Addresses This

**Safe Response Writer**

Tinh Tinh prevents duplicate `WriteHeader` calls which can cause performance issues:

```go
// core/ctx.go
type SafeResponseWriter struct {
    http.ResponseWriter
    wroteHeader bool
}

func (w *SafeResponseWriter) WriteHeader(code int) {
    if !w.wroteHeader {
        w.wroteHeader = true
        w.ResponseWriter.WriteHeader(code)
    }
}
```

**Chainable Response Methods**

The framework supports method chaining for cleaner and more efficient response handling:

```go
// core/ctx.go
func (ctx *DefaultCtx) Status(statusCode int) Ctx {
    ctx.statusCode = statusCode
    return ctx
}

// Usage: ctx.Status(http.StatusOK).JSON(data)
```

**Timeout Handling**

Built-in request timeout support prevents long-running requests from consuming resources:

```go
// core/app.go
if app.timeout != 0 {
    handler = http.TimeoutHandler(handler, app.timeout, "timeout")
}
```

---

## 4. Routing Efficiency

### Performance Factor
Route matching should be fast, with minimal overhead for path parsing and parameter extraction.

### How Tinh Tinh Addresses This

**Standard Library ServeMux**

Tinh Tinh leverages Go's standard library `http.ServeMux` which provides efficient route matching:

```go
// core/app.go
type App struct {
    Mux *http.ServeMux
    // ...
}

func CreateFactory(module ModuleParam, opt ...AppOptions) *App {
    app := &App{
        Mux: http.NewServeMux(),
        // ...
    }
    // ...
}
```

**Route Registration Cleanup**

Routes are cleaned up after registration to free memory:

```go
// core/router.go
func (app *App) registerRoutes() {
    routes := make(map[string][]*Router)
    // ... route registration ...
    for k, v := range routes {
        app.Mux.Handle(k, app.versionMiddleware(v))
        delete(routes, k)  // Clean up immediately after registration
    }
    app.free()
}
```

**Path Parsing Utilities**

Efficient string operations for path formatting:

```go
// core/router.go
func IfSlashPrefixString(s string) string {
    if s == "" {
        return s
    }
    s = strings.TrimSuffix(s, "/")
    if strings.HasPrefix(s, "/") {
        return ToFormat(s)
    }
    return "/" + ToFormat(s)
}
```

---

## 5. Dependency Injection and Provider Scoping

### Performance Factor
Dependency injection should minimize overhead, with appropriate scoping to balance memory usage and initialization costs.

### How Tinh Tinh Addresses This

**Three Provider Scopes**

Tinh Tinh supports three scopes for optimal resource management:

```go
// core/module.go
const (
    Global    Scope = "global"    // Singleton - created once, shared across all requests
    Request   Scope = "request"   // Created per request
    Transient Scope = "transient" // Created each time it's injected
)
```

**Lazy Initialization for Request Scope**

Request-scoped providers are only created when needed:

```go
// core/module.go
func requestMiddleware(module *DynamicModule) Middleware {
    return func(ctx Ctx) error {
        for _, p := range module.getRequest() {
            if p.GetValue() == nil {
                var values []interface{}
                for _, p := range p.GetInject() {
                    values = append(values, module.Ref(p, ctx))
                }
                factory := p.GetFactory()
                value := factory(values...)
                ctx.Set(p.GetName(), value)
            }
        }
        return ctx.Next()
    }
}
```

**Efficient Provider Lookup**

Using `slices.IndexFunc` for fast provider lookup:

```go
// core/module.go
func (m *DynamicModule) Ref(name Provide, ctx ...Ctx) interface{} {
    idx := slices.IndexFunc(m.DataProviders, func(e Provider) bool {
        return e.GetName() == name
    })
    if idx == -1 {
        return nil
    }
    // ...
}
```

---

## 6. Logging and Observability

### Performance Factor
Logging should have minimal impact on request latency, with options for buffering and async writes.

### How Tinh Tinh Addresses This

**Buffered File Writing**

Uses 256KB buffers for efficient file I/O:

```go
// common/logger/logger.go
const (
    defaultBufSize = 256 * 1024 // 256KB buffer
)

func (log *Logger) getOrCreateFileWriter(filepath string, level Level, t time.Time) *fileWriter {
    // ...
    fw = &fileWriter{
        file:      file,
        writer:    bufio.NewWriterSize(file, log.bufferSize),
        // ...
    }
    // ...
}
```

**Async Log Processing**

Logs are processed asynchronously to avoid blocking request handlers:

```go
// common/logger/logger.go
func Create(opt Options) *Logger {
    l := &Logger{
        // ...
        logCh:      make(chan *logEntry, channelSize), // Default: 100,000
        // ...
    }
    // Start async log process
    l.wg.Add(1)
    go l.processLog()
    // Start periodic flusher
    l.wg.Add(1)
    go l.periodicFlush()
    return l
}
```

**Periodic Flushing**

Logs are flushed periodically (every 1 second) to balance durability and performance:

```go
// common/logger/logger.go
func (log *Logger) periodicFlush() {
    defer log.wg.Done()
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    // ...
}
```

---

## 7. Caching

### Performance Factor
In-memory caching reduces database and external service calls, with efficient expiration and garbage collection.

### How Tinh Tinh Addresses This

**Time-Based Expiration**

The memory store supports TTL-based expiration:

```go
// common/memory/memory.go
type Store struct {
    ttl  time.Duration
    data map[string]item
    sync.RWMutex
}

func (m *Store) Set(key string, val interface{}, ttl ...time.Duration) {
    var exp uint32
    if len(ttl) > 0 {
        exp = uint32(ttl[0].Seconds()) + era.Timestamp()
    } else {
        exp = uint32(m.ttl.Seconds()) + era.Timestamp()
    }
    // ...
}
```

**Background Garbage Collection**

Expired items are cleaned up in the background:

```go
// common/memory/memory.go
func (m *Store) gc(sleep time.Duration) {
    ticker := time.NewTimer(sleep)
    defer ticker.Stop()
    var expired []string

    for range ticker.C {
        ts := era.Timestamp()
        expired = expired[:0]  // Reuse slice to avoid allocations
        m.RLock()
        for key, v := range m.data {
            if v.e != 0 && v.e <= ts {
                expired = append(expired, key)
            }
        }
        m.RUnlock()
        // Delete expired items...
    }
}
```

**Optimized Timestamp Calculation**

Pre-calculated timestamps for better performance:

```go
// common/era/time.go
// Using pre-calculated timestamps instead of determining at runtime each time
```

---

## 8. Serialization/Deserialization

### Performance Factor
JSON and other format encoding/decoding should be efficient, with options for custom encoders.

### How Tinh Tinh Addresses This

**Pluggable Encoders/Decoders**

The framework allows custom encoders for optimization:

```go
// core/app.go
type App struct {
    encoder Encode
    decoder Decode
    // ...
}

type AppOptions struct {
    Encoder Encode
    Decoder Decode
    // ...
}
```

**Default JSON Encoding**

Uses Go's standard library JSON encoding by default:

```go
// core/app.go
func CreateFactory(module ModuleParam, opt ...AppOptions) *App {
    app := &App{
        encoder:      json.Marshal,
        decoder:      json.Unmarshal,
        // ...
    }
    // ...
}
```

**Multiple Output Formats**

Support for JSON, XML, and template rendering:

```go
// core/ctx.go
func (ctx *DefaultCtx) JSON(data any) error { /* ... */ }
func (ctx *DefaultCtx) XML(data any) error { /* ... */ }
func (ctx *DefaultCtx) Render(name string, bind Map, layouts ...string) error { /* ... */ }
```

---

## 9. Middleware Pipeline

### Performance Factor
Middleware should be composable with minimal overhead, using efficient chaining patterns.

### How Tinh Tinh Addresses This

**Handler Chaining**

Middleware is applied in reverse order for correct execution:

```go
// core/router.go
func (r *Router) getHandler(app *App) http.Handler {
    var mergeHandler http.Handler
    // ...
    for i := len(r.Middlewares) - 1; i >= 0; i-- {
        v := r.Middlewares[i]
        mid := ParseCtxMiddleware(app, v, r)
        mergeHandler = mid(mergeHandler)
    }
    return mergeHandler
}
```

**Middleware Inheritance**

Modules inherit middleware from parent modules:

```go
// core/module.go
func (m *DynamicModule) New(opt NewModuleOptions) Module {
    newMod := &DynamicModule{isRoot: false}
    newMod.Middlewares = append(newMod.Middlewares, m.Middlewares...)
    // ...
}
```

**Guard Integration**

Guards are converted to middleware for unified processing:

```go
// core/module.go
for _, g := range opt.Guards {
    if g == nil {
        continue
    }
    mid := module.ParseGuard(g)
    module.Middlewares = append(module.Middlewares, mid)
}
```

---

## 10. Resource Cleanup

### Performance Factor
Proper cleanup of resources (connections, file handles, etc.) prevents memory leaks and resource exhaustion.

### How Tinh Tinh Addresses This

**Graceful Shutdown**

The framework supports graceful shutdown with configurable timeout:

```go
// core/app.go
func (app *App) Listen(port int) {
    // ...
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
    defer shutdownRelease()

    for _, hook := range app.hooks {
        if hook.RunAt == BEFORE_SHUTDOWN {
            hook.fnc()
        }
    }

    err := server.Shutdown(shutdownCtx)
    // ...
}
```

**Lifecycle Hooks**

Support for before and after shutdown hooks:

```go
// core/hook.go
const (
    BEFORE_SHUTDOWN RunAt = "beforeShutdown"
    AFTER_SHUTDOWN  RunAt = "afterShutdown"
)
```

**Logger Cleanup**

The logger properly drains and closes all resources:

```go
// common/logger/logger.go
func (log *Logger) Close() {
    close(log.stopCh)
    log.wg.Wait()

    // Close all file writers
    log.cacheMu.Lock()
    for _, fw := range log.fileCache {
        fw.writer.Flush()
        fw.file.Close()
    }
    log.fileCache = nil
    log.cacheMu.Unlock()
}
```

---

## Benchmarking

Tinh Tinh includes benchmark tests to measure and track performance. Run benchmarks with:

```bash
make benchmark
```

This executes:

```bash
go test ./... -benchmem -bench=. -run=^Benchmark_$
```

Key benchmarks are available in:
- `core/app_test.go` - Application startup and request handling
- `core/ctx_test.go` - Context operations
- `core/pipe_test.go` - Validation pipeline
- `core/provider_test.go` - Dependency injection
- `common/logger/log_bench_test.go` - Logging operations
- `common/memory/memory_test.go` - In-memory caching
- `dto/validator/bench_test.go` - DTO validation

---

## Recommendations for Optimal Performance

1. **Use Global Scope for Stateless Services**: Services that don't hold request-specific state should use `Global` scope to avoid per-request instantiation overhead.

2. **Enable Connection Pooling**: When using databases or external services, configure connection pooling to reuse connections.

3. **Configure Logger Buffer Size**: For high-throughput applications, increase the logger buffer size:
   ```go
   logger.Create(logger.Options{
       BufferSize: 500000, // Default is 100,000
   })
   ```

4. **Use Request Timeout**: Always configure a request timeout to prevent resource exhaustion:
   ```go
   core.CreateFactory(module, core.AppOptions{
       Timeout: 30 * time.Second,
   })
   ```

5. **Implement Graceful Shutdown Hooks**: Register cleanup hooks to properly release resources:
   ```go
   app.AddHook(core.BEFORE_SHUTDOWN, func() {
       // Close database connections, etc.
   })
   ```

---

## Contributing

If you identify additional performance improvements, please:
1. Open an issue describing the performance concern
2. Include benchmark results if possible
3. Submit a PR with the proposed optimization

See [CONTRIBUTING.md](CONTRIBUTING.md) for more details.
