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

type Logger struct {
	Path   string
	Rotate bool
	// Max Size in MB of each file log. Default is infinity.
	Max int64
}

type Options struct {
	// Log path. Default is "logs".
	Path string
	// Rotate log files. Default is false.
	Rotate bool
	// Max Size in MB of each file log. Default is infinity.
	Max int64
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
	return &Logger{
		Path:   opt.Path,
		Rotate: opt.Rotate,
		Max:    opt.Max,
	}
}

func (log *Logger) Info(msg string) {
	log.write(LevelInfo, msg)
}

func (log *Logger) Debug(msg string) {
	log.write(LevelDebug, msg)
}

func (log *Logger) Warn(msg string) {
	log.write(LevelWarn, msg)
}

func (log *Logger) Error(msg string) {
	log.write(LevelError, msg)
}

func (log *Logger) Fatal(msg string) {
	log.write(LevelFatal, msg)
}

func (log *Logger) Log(level Level, msg string) {
	log.write(level, msg)
}

func (log *Logger) write(level Level, msg string) {
	dir := log.Path
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		check(err)
	}

	current := time.Now().Format("2006-01-02")
	fileName := current + "-" + getLevelName(level)
	file, _ := os.OpenFile(filepath.Join(dir, fileName+".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	fi, _ := file.Stat()
	currentSize := fi.Size()

	if log.Max > 0 && currentSize > log.Max*1000*1000 {
		if !log.Rotate {
			file.Close()
			panic("Size log is hit limited storage")
		} else {
			idx := 1
			for idx > 0 {
				file, _ = os.OpenFile(filepath.Join(dir, fileName+"-"+fmt.Sprint(idx)+".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
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

	message := time.Now().Format("2006-01-02 15:04:05") + " " + msg + "\n"
	_, err := iw.Write([]byte(message))
	if err != nil {
		panic(err)
	}
}
