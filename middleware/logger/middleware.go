package logger

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type MiddlewareOptions struct {
	Path   string
	Rotate bool
	// Max Size in MB of each file log. Default is infinity.
	Max    int64
	Format string
	Level  Level
}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func Middleware(opt MiddlewareOptions) func(http.Handler) http.Handler {
	logger := Create(Options{
		Path:   opt.Path,
		Rotate: opt.Rotate,
		Max:    opt.Max,
	})
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
			logger.Info(content)
		})
	}
}
