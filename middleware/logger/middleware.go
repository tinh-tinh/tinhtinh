package logger

import (
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	clogger "github.com/tinh-tinh/tinhtinh/v2/common/logger"
)

const (
	Dev      string = "${ip} - ${method} ${path} ${status} ${latency}"
	Common   string = "${ip}:${date} - '${method} ${path} ${http-version}' ${status} ${content-length}"
	Combined string = "${ip}:${date} - '${method} ${path} ${http-version}' ${status} ${content-length} ${latency} - ${referer}:${user-agent}"
)

// LogContext contains all the information available for custom log formatting.
type LogContext struct {
	// Request is the original HTTP request
	Request *http.Request
	// StatusCode is the response status code
	StatusCode int
	// Latency is the time taken to process the request
	Latency time.Duration
	// ResponseHeaders contains the response headers
	ResponseHeaders http.Header
	// StartTime is when the request started
	StartTime time.Time
}

// CustomFormatter is a function type for custom log formatting.
// It receives LogContext and returns the formatted log string.
type CustomFormatter func(ctx LogContext) string

type MiddlewareOptions struct {
	Path               string
	Rotate             bool
	SeparateBaseStatus bool
	// Max Size in MB of each file log. Default is infinity.
	Max int64
	// Format is the log message template format (e.g., Dev, Common, Combined).
	Format string
	// OutputFormat specifies the log output format (text or json).
	// Uses clogger.FormatText or clogger.FormatJSON.
	OutputFormat clogger.Format
	Level        clogger.Level
	// CustomFormatter allows fully custom log formatting.
	// When set, this takes precedence over Format.
	CustomFormatter CustomFormatter
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
//
// Alternatively, use CustomFormatter for fully custom log formatting:
//
//	logger.Handler(logger.MiddlewareOptions{
//	    CustomFormatter: func(ctx logger.LogContext) string {
//	        return fmt.Sprintf("[%s] %s %s - %d (%s)",
//	            ctx.Request.RemoteAddr,
//	            ctx.Request.Method,
//	            ctx.Request.URL.Path,
//	            ctx.StatusCode,
//	            ctx.Latency,
//	        )
//	    },
//	})
//
// The OutputFormat option specifies the log output format:
// - clogger.FormatText (default): key=value text format
// - clogger.FormatJSON: JSON format
func Handler(opt MiddlewareOptions) func(http.Handler) http.Handler {
	logger := clogger.Create(clogger.Options{
		Path:   opt.Path,
		Rotate: opt.Rotate,
		Max:    opt.Max,
		Format: opt.OutputFormat,
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

			ctx := LogContext{
				Request:         r,
				StatusCode:      wrapped.statusCode,
				Latency:         time.Since(start),
				ResponseHeaders: wrapped.Header(),
				StartTime:       start,
			}

			content := formatLogMessage(opt, ctx)

			if opt.SeparateBaseStatus {
				level = SeparateBaseStatus(wrapped.statusCode)
			}
			logger.Log(level, content)
		})
	}
}

// formatLogMessage formats the log message based on options and context.
func formatLogMessage(opt MiddlewareOptions, ctx LogContext) string {
	if opt.CustomFormatter != nil {
		return opt.CustomFormatter(ctx)
	}
	return formatTemplate(opt.Format, ctx)
}

// formatTemplate formats the log message using template placeholders.
func formatTemplate(format string, ctx LogContext) string {
	if format == "" {
		format = Dev
	}

	content := format
	specs := clogger.ExtractAllContent(format)

	for _, spec := range specs {
		content = replaceSpec(content, spec, ctx)
	}

	return content
}

// replaceSpec replaces all occurrences of a placeholder spec with its value.
func replaceSpec(content, spec string, ctx LogContext) string {
	r := ctx.Request

	switch spec {
	case "http-version":
		return strings.ReplaceAll(content, "${http-version}", r.Proto)
	case "user-agent":
		return strings.ReplaceAll(content, "${user-agent}", html.EscapeString(r.UserAgent()))
	case "referer":
		return strings.ReplaceAll(content, "${referer}", html.EscapeString(r.Referer()))
	case "status":
		return strings.ReplaceAll(content, "${status}", fmt.Sprint(ctx.StatusCode))
	case "method":
		return strings.ReplaceAll(content, "${method}", r.Method)
	case "path":
		return strings.ReplaceAll(content, "${path}", html.EscapeString(r.URL.Path))
	case "ip":
		return strings.ReplaceAll(content, "${ip}", r.RemoteAddr)
	case "content-length":
		return strings.ReplaceAll(content, "${content-length}", ctx.ResponseHeaders.Get("Content-Length"))
	case "latency":
		return strings.ReplaceAll(content, "${latency}", ctx.Latency.String())
	case "date":
		return strings.ReplaceAll(content, "${date}", ctx.StartTime.Format("2006-01-02 15:04:05"))
	default:
		return content
	}
}

// SeparateBaseStatus returns the log level based on HTTP status code.
func SeparateBaseStatus(statusCode int) clogger.Level {
	if statusCode < 300 {
		return clogger.LevelInfo
	} else if statusCode < 400 {
		return clogger.LevelWarn
	} else if statusCode < 500 {
		return clogger.LevelError
	} else {
		return clogger.LevelFatal
	}
}
