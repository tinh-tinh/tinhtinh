package logger

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

const (
	MiB            = 1 << 20    // 1 MiB
	defaultBufSize = 256 * 1024 // 256KB buffer
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

type fileWriter struct {
	file      *os.File
	writer    *bufio.Writer
	path      string
	lastFlush time.Time
	mu        sync.Mutex
	size      int64
}

type logEntry struct {
	level Level
	msg   string
	meta  []Metadata
	time  time.Time
}

type Logger struct {
	Options
	// Performance factors
	mu         sync.Mutex
	fileCache  map[string]*fileWriter
	cacheMu    sync.RWMutex
	stopCh     chan struct{}
	wg         sync.WaitGroup
	logCh      chan *logEntry
	bufferSize int
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
	l := &Logger{
		Options:    opt,
		fileCache:  make(map[string]*fileWriter),
		stopCh:     make(chan struct{}),
		logCh:      make(chan *logEntry, 100000),
		bufferSize: defaultBufSize,
	}

	if err := os.MkdirAll(l.Path, 0o755); err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
	}

	// Start async log process
	l.wg.Add(1)
	go l.processLog()

	// Start periodic flusher
	l.wg.Add(1)
	go l.periodicFlush()

	return l
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

func (log *Logger) writeEntryLog(entry *logEntry) {
	fileName := entry.time.Format("2006-01-02") + "-" + GetLevelName(entry.level) + ".log"
	filePath := filepath.Join(log.Path, fileName)

	fw := log.getOrCreateFileWriter(filePath, entry.level, entry.time)
	if fw == nil {
		return
	}

	merged := appendMetadata(log.Metadata, entry.meta...)
	if log.TraceDepth > 0 {
		pc, _, _, ok := runtime.Caller(log.TraceDepth)
		if ok {
			fn := runtime.FuncForPC(pc)
			merged["trace"] = fn.Name()
		}
	}

	metaStr := ""
	for k, v := range merged {
		metaStr += fmt.Sprintf("[%s=%v] ", k, v)
	}

	message := fmt.Sprintf("%s [%s] %s%s\n",
		entry.time.Format("2006-01-02 15:04:05"),
		GetLevelName(entry.level),
		metaStr,
		entry.msg,
	)

	fw.mu.Lock()
	n, err := fw.writer.WriteString(message)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write log: %v\n", err)
	} else {
		fw.size += int64(n)
	}
	fw.mu.Unlock()

	// Also write to stdout (unbuffered for immediate visibility)
	fmt.Print(message)

	if log.Max > 0 && fw.size >= log.Max*MiB {
		log.rotateFile(filePath, entry.level, entry.time)
	}
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
		// Channel full, log to stderr
		fmt.Fprintf(os.Stderr, "LOG CHANNEL FULL: [%s] %s\n", GetLevelName(level), msg)
	}
}

func (log *Logger) processLog() {
	defer log.wg.Done()

	for {
		select {
		case entry := <-log.logCh:
			log.writeEntryLog(entry)
		case <-log.stopCh:
			// Drain remaining logs
			for len(log.logCh) > 0 {
				entry := <-log.logCh
				log.writeEntryLog(entry)
			}
			return
		}
	}
}

func (log *Logger) periodicFlush() {
	defer log.wg.Done()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.flushAll()
		case <-log.stopCh:
			log.flushAll()
			return
		}
	}
}

func (log *Logger) flushAll() {
	log.cacheMu.RLock()
	defer log.cacheMu.RUnlock()

	for _, writer := range log.fileCache {
		writer.mu.Lock()
		writer.writer.Flush()
		writer.lastFlush = time.Now()
		writer.mu.Unlock()
	}
}

func (log *Logger) getOrCreateFileWriter(filepath string, level Level, t time.Time) *fileWriter {
	log.cacheMu.RLock()
	fw, exists := log.fileCache[filepath]
	log.cacheMu.RUnlock()

	if exists {
		return fw
	}

	// If need create new file writer
	log.cacheMu.Lock()
	defer log.cacheMu.Unlock()

	// Double-check after acquiring write lock
	fw, exists = log.fileCache[filepath]
	if exists {
		return fw
	}

	flags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	// Check if file needs rotation before opening
	if log.Max > 0 {
		if fi, err := os.Stat(filepath); err == nil {
			if fi.Size() >= log.Max*MiB {
				filepath = log.getRotatedPath(filepath, level, t)
				if !log.Rotate {
					flags = os.O_TRUNC | os.O_CREATE | os.O_WRONLY
				}
			}
		}
	}

	file, err := os.OpenFile(filepath, flags, 0o666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
		return nil
	}

	var initialSize int64
	if fi, err := file.Stat(); err == nil {
		initialSize = fi.Size()
	}

	fw = &fileWriter{
		file:      file,
		writer:    bufio.NewWriterSize(file, log.bufferSize),
		path:      filepath,
		lastFlush: time.Now(),
		size:      initialSize,
	}

	log.fileCache[filepath] = fw
	return fw
}

func (log *Logger) getRotatedPath(basepath string, level Level, t time.Time) string {
	if !log.Rotate {
		return basepath
	}

	dir := filepath.Dir(basepath)
	current := t.Format(time.DateOnly)
	levelName := GetLevelName(level)

	for idx := 1; ; idx++ {
		fileName := fmt.Sprintf("%s-%s-%d.log", current, levelName, idx)
		newPath := filepath.Join(dir, fileName)
		if checkAvailableFile(newPath, log.Max) {
			return newPath
		}
	}
}

func (log *Logger) rotateFile(oldPath string, level Level, t time.Time) {
	log.cacheMu.Lock()
	defer log.cacheMu.Unlock()

	if fw, exists := log.fileCache[oldPath]; exists {
		fw.mu.Lock()
		fw.writer.Flush()
		fw.file.Close()
		fw.mu.Unlock()
		delete(log.fileCache, oldPath)
	}
}

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
