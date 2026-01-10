package logger_test

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/common/logger"
)

type serviceImpl struct {
	log *flog
}

func (s *serviceImpl) DoSomething(l *logger.Logger) {
	l.Debug("Doing something in Service")
}

func (s *serviceImpl) LogInternal() {
	s.log.DoSomething()
}

type flog struct {
	log *logger.Logger
}

func NewFlog(l *logger.Logger) *flog {
	return &flog{
		log: l,
	}
}

func (f *flog) DoSomething() {
	f.log.Debug("Doing something in Service")
}

func Test_HappyLog(t *testing.T) {
	logFullpath := logger.Create(logger.Options{
		Max:        1,
		Rotate:     true,
		TraceDepth: logger.TracerFullPath,
		Metadata: logger.Metadata{
			"svc": "AuthSvc",
		},
	})

	svc := &serviceImpl{}
	svc.DoSomething(logFullpath)

	logEntryfile := logger.Create(logger.Options{
		Max:        1,
		Rotate:     true,
		TraceDepth: logger.TracerEntryFile,
	})
	svc.DoSomething(logEntryfile)

	logFuncOnly := logger.Create(logger.Options{
		Max:        1,
		Rotate:     true,
		TraceDepth: logger.TraceOnlyFunc,
	})
	svc.DoSomething(logFuncOnly)

	logAnother := logger.Create(logger.Options{
		Max:        1,
		Rotate:     true,
		TraceDepth: 4,
	})
	f := NewFlog(logAnother)
	svc2 := &serviceImpl{
		log: f,
	}
	svc2.LogInternal()
}

func Test_Create(t *testing.T) {
	l := logger.Create(logger.Options{
		Max:    1,
		Rotate: true,
	})
	for i := range 1000 {
		val := strconv.Itoa(i)
		if i%2 == 0 {
			l.Info(val)
		} else if i%3 == 0 {
			l.Warn(val)
		} else if i%5 == 0 {
			l.Error(val)
		} else if i%7 == 0 {
			l.Fatal(val)
		} else {
			l.Debug(val)
		}
	}

	log2 := logger.Create(logger.Options{
		Path:   "logs/test",
		Max:    1,
		Rotate: false,
	})

	require.NotPanics(t, func() {
		for range 2 {
			log2.Info(randomBigStr())
		}
	})

	log3 := logger.Create(logger.Options{
		Path:   "logs/test2",
		Max:    1,
		Rotate: true,
	})

	for range 2 {
		log3.Info(randomBigStr())
	}

	l = logger.Create(logger.Options{
		Max:    1,
		Rotate: true,
	})
	for i := range 100 {
		if i%2 == 0 {
			l.Infof("The value is %d", i)
		} else if i%3 == 0 {
			l.Warnf("The value is %d", i)
		} else if i%5 == 0 {
			l.Errorf("The value is %d", i)
		} else if i%7 == 0 {
			l.Fatalf("The value is %d", i)
		} else {
			l.Debugf("The value is %d", i)
		}
		l.Logf(logger.LevelDebug, "alayws have ata %d", i)
	}

	for i := range 1000 {
		val := strconv.Itoa(i)
		if i%2 == 0 {
			l.Info(val, logger.Metadata{
				"function": "Test",
			})
		} else if i%3 == 0 {
			l.Warn(val, logger.Metadata{
				"function": "Test",
			})
		} else if i%5 == 0 {
			l.Error(val, logger.Metadata{
				"function": "Test",
			})
		} else if i%7 == 0 {
			l.Fatal(val, logger.Metadata{
				"function": "Test",
			})
		} else {
			l.Debug(val, logger.Metadata{
				"function": "Test",
			})
		}
	}
}

func randomBigStr() string {
	var bigString strings.Builder
	// Define the number of repetitions
	repeat := 100000
	smallString := "Hello, Go! "

	// Append the small string multiple times
	for i := 0; i < repeat; i++ {
		bigString.WriteString(smallString)
	}

	// Convert the builder to a string
	result := bigString.String()
	return result
}

func Test_MkdirError(t *testing.T) {
	// Test with invalid path that cannot be created (e.g., path with null byte)
	l := logger.Create(logger.Options{
		Path: "/dev/null/invalid_path",
		Max:  1,
	})

	// Should not panic, just gracefully handle the error
	require.NotPanics(t, func() {
		l.Info("test message")
	})
}

func Test_FileOpenError(t *testing.T) {
	// Test with a path that exists but is not writable (directory as file)
	l := logger.Create(logger.Options{
		Path: "/proc", // exists but cannot create files here
		Max:  1,
	})

	// Should not panic, just gracefully handle the error
	require.NotPanics(t, func() {
		l.Info("test message")
	})
}

func Test_InvalidPathError(t *testing.T) {
	// Test with path containing invalid characters
	l := logger.Create(logger.Options{
		Path: string([]byte{0x00}), // null byte in path
		Max:  1,
	})

	// Should not panic, just gracefully handle the error
	require.NotPanics(t, func() {
		l.Info("test message")
	})
}

func Test_MkdirPermissionDenied(t *testing.T) {
	// Test mkdir error when trying to create directory in a read-only location
	// /sys/firmware is a sysfs path that doesn't allow directory creation
	// This triggers the os.Mkdir error path specifically (line 118-122)
	l := logger.Create(logger.Options{
		Path: "/sys/firmware/new_log_dir", // sysfs - cannot create directories here
		Max:  1,
	})

	// Should not panic, just gracefully handle the mkdir error
	require.NotPanics(t, func() {
		l.Info("test message - mkdir permission denied")
	})
}

func Test_WriteToReadOnlyFile(t *testing.T) {
	// Create a temporary directory for this test
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "readonly_test")

	// Create the log directory first
	err := os.MkdirAll(logPath, 0o755)
	require.NoError(t, err)

	// Create a logger that will write to this directory
	l := logger.Create(logger.Options{
		Path: logPath,
		Max:  1,
	})

	// Write first log to create the file
	l.Info("first message")

	// Make the directory read-only to trigger write errors on subsequent writes
	err = os.Chmod(logPath, 0o444)
	require.NoError(t, err)

	// Restore permissions after test
	defer func() {
		_ = os.Chmod(logPath, 0o755)
	}()

	// Try to write again - this won't trigger iw.Write error but tests the path
	// The file is already open, so it will still work, but good coverage test
	require.NotPanics(t, func() {
		l.Info("second message after chmod")
	})
}
