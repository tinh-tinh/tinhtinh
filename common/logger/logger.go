package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const MiB = 1 << 20 // 1 MiB

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

const (
	TraceOnlyFunc = iota + 1
	TracerEntryFile
	TracerFullPath
)

type Metadata map[string]any
type Logger struct {
	Options
}

type Options struct {
	// Log path. Default is "logs".
	Path string
	// Rotate log files. Default is false.
	Rotate bool
	// Max Size in MB of each file log. Default is infinity.
	Max int64
	// metadata
	Metadata Metadata
	// TraceDepth enables including the caller function name in debug logs.
	TraceDepth int
}

// Create a new Logger with the specified options.
//
// The created logger will have the given path for the log files. If the path is
// empty, the default value "logs" is used. The logger will rotate log files if
// the Rotate option is true. The maximum size of each log file can be set with
// the Max option. The default value is infinity.
func Create(opt Options) *Logger {
	if opt.Path == "" {
		opt.Path = "logs"
	}
	if opt.Max == 0 {
		opt.Max = 20
	}
	return &Logger{
		Options: opt,
	}
}

func (log *Logger) Info(msg string, meta ...Metadata) {
	log.write(LevelInfo, msg, meta...)
}

func (log *Logger) Infof(msg string, args ...any) {
	log.write(LevelInfo, fmt.Sprintf(msg, args...))
}

func (log *Logger) Debug(msg string, meta ...Metadata) {
	log.write(LevelDebug, msg, meta...)
}

func (log *Logger) Debugf(msg string, args ...any) {
	log.write(LevelDebug, fmt.Sprintf(msg, args...))
}

func (log *Logger) Warn(msg string, meta ...Metadata) {
	log.write(LevelWarn, msg, meta...)
}

func (log *Logger) Warnf(msg string, args ...any) {
	log.write(LevelWarn, fmt.Sprintf(msg, args...))
}

func (log *Logger) Error(msg string, meta ...Metadata) {
	log.write(LevelError, msg, meta...)
}

func (log *Logger) Errorf(msg string, args ...any) {
	log.write(LevelError, fmt.Sprintf(msg, args...))
}

func (log *Logger) Fatal(msg string, meta ...Metadata) {
	log.write(LevelFatal, msg, meta...)
}

func (log *Logger) Fatalf(msg string, args ...any) {
	log.write(LevelFatal, fmt.Sprintf(msg, args...))
}

func (log *Logger) Log(level Level, msg string, meta ...Metadata) {
	log.write(level, msg, meta...)
}

func (log *Logger) Logf(level Level, msg string, args ...any) {
	log.write(level, fmt.Sprintf(msg, args...))
}

func (log *Logger) write(level Level, msg string, meta ...Metadata) {
	dir := log.Path
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0o755)
		if err != nil {
			fmt.Printf("Failed to create log directory: %v\n", err)
			return
		}
	}

	current := time.Now().Format("2006-01-02")
	fileName := current + "-" + GetLevelName(level) + ".log"
	filePath := filepath.Join(dir, fileName)

	flags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	if !checkAvailableFile(filePath, log.Max) {
		if log.Rotate {
			idx := 1
			for idx > 0 {
				fileName = current + "-" + GetLevelName(level) + "-" + fmt.Sprint(idx) + ".log"
				filePath = filepath.Join(dir, fileName)
				if checkAvailableFile(filePath, log.Max) {
					break
				}
				idx++
			}
		} else {
			flags = os.O_TRUNC | os.O_CREATE | os.O_WRONLY
		}
	}

	file, err := os.OpenFile(filePath, flags, 0o666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}
	defer file.Close()

	// Use io.MultiWriter to write to both stdout and the file
	iw := io.MultiWriter(os.Stdout, file)

	// Merge default and per-call metadata
	merged := appendMetadata(log.Metadata, meta...)
	if log.TraceDepth > 0 {
		pc, _, _, ok := runtime.Caller(log.TraceDepth)
		if ok {
			fn := runtime.FuncForPC(pc)
			merged["trace"] = traceDepthName(log.TraceDepth, fn)
		}
	}

	metaStr := ""
	for k, v := range merged {
		metaStr += fmt.Sprintf("[%s=%s] ", k, v)
	}

	message := fmt.Sprintf("%s [%s] %s%s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		GetLevelName(level),
		metaStr,
		msg,
	)
	_, err = iw.Write([]byte(message))
	if err != nil {
		fmt.Printf("Failed to write log: %v\n", err)
		return
	}
}

func checkAvailableFile(filename string, max int64) bool {
	if max <= 0 {
		return true
	}
	fi, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// no file yet = available
		return true
	}
	if err != nil {
		fmt.Printf("Failed to check log file: %v\n", err)
		return false
	}
	return fi.Size() < max*MiB
}

func appendMetadata(base Metadata, extra ...Metadata) Metadata {
	merged := make(Metadata)
	for k, v := range base {
		merged[k] = v
	}
	if len(extra) > 0 {
		for _, m := range extra {
			for k, v := range m {
				merged[k] = v
			}
		}
	}

	return merged
}

func traceDepthName(depth int, fn *runtime.Func) string {
	fullName := fn.Name()
	switch depth {
	case TraceOnlyFunc:
		splits := strings.Split(fullName, ".")
		shortName := splits[len(splits)-1]
		return strings.SplitN(shortName, "(", 2)[0]
	case TracerEntryFile:
		entryIndex := strings.LastIndex(fullName, "/")
		entryFile := fullName[entryIndex+1:]
		return entryFile
	case TracerFullPath:
		return fullName
	default:
		// Fallback to just the file name
		return fullName
	}
}
