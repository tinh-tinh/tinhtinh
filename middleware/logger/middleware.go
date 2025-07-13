package logger

import (
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"
)

const (
	Dev      string = "${ip} - ${method} ${path} ${status} ${latency}"
	Common   string = "${ip}:${date} - '${method} ${path} ${http-version}' ${status} ${content-length}"
	Combined string = "${ip}:${date} - '${method} ${path} ${http-version}' ${status} ${content-length} ${latency} - ${referer}:${user-agent}"
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

// Handler returns a middleware that logs the request and response with the
// given options.
//
// The level of the log message is determined by the StatusCode of the response.
// If SeparateBaseStatus is true, the level is determined by the base status
// (100-500) of the response. Otherwise, the level is set to the given level.
//
// The format of the log message can be customized with the Format option. The
// following formats are available:
//
// - dev: ${ip} - ${method} ${path} ${status} ${latency}
// - common: ${ip}:${date} - '${method} ${path} ${http-version}' ${status} ${content-length}
// - combined: ${ip}:${date} - '${method} ${path} ${http-version}' ${status} ${content-length} ${latency} - ${referer}:${user-agent}
//
// The format can be customized with variables in the format string. The
// following variables are available:
//
// - ${http-version}: the HTTP version of the request
// - ${user-agent}: the User-Agent header of the request
// - ${referer}: the Referer header of the request
// - ${status}: the StatusCode of the response
// - ${method}: the method of the request
// - ${path}: the path of the request
// - ${ip}: the IP address of the client
// - ${content-length}: the Content-Length header of the response
// - ${latency}: the latency of the request in milliseconds
// - ${date}: the current date and time in the format 2006-01-02 15:04:05
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
			format := opt.Format
			if format == "" {
				format = Dev
			}
			content := format
			specs := extractAllContent(format)
			for _, spec := range specs {
				switch spec {
				case "http-version":
					content = strings.Replace(content, "${http-version}", r.Proto, 1)
				case "user-agent":
					content = strings.Replace(content, "${user-agent}", html.EscapeString(r.UserAgent()), 1)
				case "referer":
					content = strings.Replace(content, "${referer}", r.Referer(), 1)
				case "status":
					content = strings.Replace(content, "${status}", fmt.Sprint(wrapped.statusCode), 1)
				case "method":
					content = strings.Replace(content, "${method}", r.Method, 1)
				case "path":
					content = strings.Replace(content, "${path}", r.URL.Path, 1)
				case "ip":
					content = strings.Replace(content, "${ip}", r.RemoteAddr, 1)
				case "content-length":
					content = strings.Replace(content, "${content-length}", fmt.Sprint(wrapped.Header().Get("Content-Length")), 1)
				case "latency":
					elapsed := time.Since(start)
					content = strings.Replace(content, "${latency}", elapsed.String(), 1)
				case "date":
					content = strings.Replace(content, "${date}", time.Now().Format("2006-01-02 15:04:05"), 1)
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
