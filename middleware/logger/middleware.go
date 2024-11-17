package logger

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type MiddlewareOptions struct {
	Path               string
	Rotate             bool
	SeparateBaseStatus bool
	// Max Size in MB of each file log. Default is infinity.
	Max    int64
	Format string
	Level  Level
}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

// Handler returns a middleware that logs the request and response.
//
// The middleware will log the request method, path, remote address, response status
// and latency. The format of the log message is configurable with the Format
// option. The format string can contain the following placeholders:
// - ${method}: the request method
// - ${path}: the request path
// - ${ip}: the remote address
// - ${status}: the response status
// - ${latency}: the latency of the request
//
// The Path option specifies the path of the log files. The Rotate option specifies
// whether the log files should be rotated. The Max option specifies the maximum
// size of each log file. The unit of the size is MB. The default value is
// infinity. The Level option specifies the level of the log messages. The level
// can be one of the following:
// - LevelFatal: log messages with a fatal level
// - LevelError: log messages with an error level
// - LevelWarn: log messages with a warning level
// - LevelInfo: log messages with an info level
// - LevelDebug: log messages with a debug level
//
// The middleware will log the request and response with the specified level.
func Handler(opt MiddlewareOptions) func(http.Handler) http.Handler {
	logger := Create(Options{
		Path:   opt.Path,
		Rotate: opt.Rotate,
		Max:    opt.Max,
	})
	level := opt.Level

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &wrappedWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}
			next.ServeHTTP(wrapped, r)
			content := opt.Format
			specs := extractAllContent(opt.Format)
			for _, spec := range specs {
				switch spec {
				case "status":
					content = strings.Replace(content, "${status}", fmt.Sprint(wrapped.statusCode), 1)
				case "method":
					content = strings.Replace(content, "${method}", r.Method, 1)
				case "path":
					content = strings.Replace(content, "${path}", r.URL.Path, 1)
				case "ip":
					content = strings.Replace(content, "${ip}", r.RemoteAddr, 1)
				case "latency":
					elapsed := time.Since(start)
					content = strings.Replace(content, "${latency}", elapsed.String(), 1)
				}
			}
			if opt.SeparateBaseStatus {
				level = SeparateBaseStatus(wrapped.statusCode)
			}
			logger.Log(level, content)
		})
	}
}

func SeparateBaseStatus(statusCode int) Level {
	if statusCode < 300 {
		return LevelInfo
	} else if statusCode < 400 {
		return LevelWarn
	} else if statusCode < 500 {
		return LevelError
	} else {
		return LevelFatal
	}
}
