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
	LevelFatal Level = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)

type Logger struct {
	Path   string
	Rotate bool
	// Max Size in MB of each file log. Default is infinity.
	Max int64
}

type Options struct {
	Path   string
	Rotate bool
	// Max Size in MB of each file log. Default is infinity.
	Max int64
}

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
			fmt.Println("Size log is hit limited storage")
			file.Close()
			return
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
