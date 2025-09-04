package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

type Metadata map[string]any
type Logger struct {
	Path   string
	Rotate bool
	// Max Size in MB of each file log. Default is infinity.
	Max      int64
	Metadata Metadata
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
		Path:     opt.Path,
		Rotate:   opt.Rotate,
		Max:      opt.Max,
		Metadata: opt.Metadata,
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
		check(err)
	}

	current := time.Now().Format("2006-01-02")
	fileName := current + "-" + getLevelName(level)
	file, _ := os.OpenFile(filepath.Join(dir, fileName+".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	fi, _ := file.Stat()
	currentSize := fi.Size()

	if log.Max > 0 && currentSize > log.Max*1000*1000 {
		if !log.Rotate {
			file.Close()
			panic("Size log is hit limited storage")
		} else {
			idx := 1
			for idx > 0 {
				file, _ = os.OpenFile(filepath.Join(dir, fileName+"-"+fmt.Sprint(idx)+".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
				fi, _ := file.Stat()
				currentSize = fi.Size()
				if currentSize < log.Max*1000*1000 {
					break
				}
				idx++
			}
		}
	}

	iw := io.MultiWriter(os.Stdout, file)

	// Merge default and per-call metadata
	merged := make(Metadata)
	for k, v := range log.Metadata {
		merged[k] = v
	}
	if len(meta) > 0 {
		for k, v := range meta[0] {
			merged[k] = v
		}
	}
	metaStr := ""
	for k, v := range merged {
		metaStr += fmt.Sprintf("[%s=%s] ", k, v)
	}

	message := fmt.Sprintf("%s [%s] %s%s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		getLevelName(level),
		metaStr,
		msg,
	)
	_, err := iw.Write([]byte(message))
	if err != nil {
		panic(err)
	}
}
