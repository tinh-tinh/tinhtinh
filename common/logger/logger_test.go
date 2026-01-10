package logger_test

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

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
	logWithTrace := logger.Create(logger.Options{
		Max:        1,
		Rotate:     true,
		TraceDepth: 3, // Enable trace at depth 3
		Metadata: logger.Metadata{
			"svc": "AuthSvc",
		},
	})

	svc := &serviceImpl{}
	svc.DoSomething(logWithTrace)

	logWithTrace2 := logger.Create(logger.Options{
		Max:        1,
		Rotate:     true,
		TraceDepth: 2,
	})
	svc.DoSomething(logWithTrace2)

	logWithTrace3 := logger.Create(logger.Options{
		Max:        1,
		Rotate:     true,
		TraceDepth: 1,
	})
	svc.DoSomething(logWithTrace3)

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

	time.Sleep(time.Second * 3)
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

	time.Sleep(time.Second * 2)
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

func Test_LogRotation(t *testing.T) {
	t.Run("rotate option true", func(t *testing.T) {
		tmpDir := t.TempDir()
		l := logger.Create(logger.Options{
			Path:   tmpDir,
			Max:    1, // 1 MB
			Rotate: true,
		})

		l.Info("test rotation enabled")
		time.Sleep(100 * time.Millisecond)
		l.Close()

		// Verify log file was created
		files, err := os.ReadDir(tmpDir)
		require.NoError(t, err)
		require.NotEmpty(t, files)
	})

	t.Run("rotate option false", func(t *testing.T) {
		tmpDir := t.TempDir()
		l := logger.Create(logger.Options{
			Path:   tmpDir,
			Max:    1, // 1 MB
			Rotate: false,
		})

		l.Info("test rotation disabled")
		time.Sleep(100 * time.Millisecond)
		l.Close()

		// Verify log file was created
		files, err := os.ReadDir(tmpDir)
		require.NoError(t, err)
		require.NotEmpty(t, files)
	})

	t.Run("rotation file naming", func(t *testing.T) {
		tmpDir := t.TempDir()
		today := time.Now().Format("2006-01-02")

		l := logger.Create(logger.Options{
			Path:   tmpDir,
			Max:    100, // 100 MB - won't trigger actual rotation
			Rotate: true,
		})

		l.Info("test log message")
		l.Debug("debug message")
		l.Warn("warn message")
		l.Error("error message")

		time.Sleep(200 * time.Millisecond)
		l.Close()

		// Verify correct file naming: YYYY-MM-DD-level.log
		files, err := os.ReadDir(tmpDir)
		require.NoError(t, err)

		for _, file := range files {
			name := file.Name()
			require.Contains(t, name, today, "file should contain date")
			require.True(t,
				strings.Contains(name, "-info.log") ||
					strings.Contains(name, "-debug.log") ||
					strings.Contains(name, "-warn.log") ||
					strings.Contains(name, "-error.log"),
				"file should have correct level suffix: %s", name)
		}
	})

	t.Run("rotation creates numbered files when max exceeded", func(t *testing.T) {
		tmpDir := t.TempDir()
		today := time.Now().Format("2006-01-02")

		// Pre-create a log file that exceeds max size to trigger rotation
		existingLogPath := filepath.Join(tmpDir, today+"-info.log")
		err := os.WriteFile(existingLogPath, make([]byte, 1024*1024+1), 0o644)
		require.NoError(t, err)

		l := logger.Create(logger.Options{
			Path:   tmpDir,
			Max:    1, // 1 MB
			Rotate: true,
		})

		// Write new log - this should detect existing file exceeds max and create rotated file
		l.Info("message triggering rotation")
		time.Sleep(200 * time.Millisecond)
		l.Close()

		// Verify rotated file with numbered suffix was created
		files, err := os.ReadDir(tmpDir)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(files), 2, "should have original + rotated file")

		// Check for numbered rotation file: YYYY-MM-DD-info-1.log
		foundNumberedFile := false
		for _, file := range files {
			name := file.Name()
			if strings.HasPrefix(name, today+"-info-") && strings.HasSuffix(name, ".log") {
				foundNumberedFile = true
				// Verify it matches pattern like 2026-01-10-info-1.log
				require.Regexp(t, `^\d{4}-\d{2}-\d{2}-info-\d+\.log$`, name)
				break
			}
		}
		require.True(t, foundNumberedFile, "should have created numbered rotation file (e.g., %s-info-1.log)", today)
	})
}

// Edge case tests

func Test_LoggerClose(t *testing.T) {
	tmpDir := t.TempDir()
	l := logger.Create(logger.Options{
		Path: tmpDir,
		Max:  100,
	})

	// Write many logs to ensure some are pending when Close is called
	for i := 0; i < 1000; i++ {
		l.Infof("pending log message %d", i)
	}

	// Close should drain all pending logs and flush buffers
	require.NotPanics(t, func() {
		l.Close()
	})

	// Verify log files were created and content was flushed
	files, err := os.ReadDir(tmpDir)
	require.NoError(t, err)
	require.NotEmpty(t, files, "log files should be created")

	// Read file content to verify flush worked
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".log") {
			content, err := os.ReadFile(filepath.Join(tmpDir, file.Name()))
			require.NoError(t, err)
			require.Contains(t, string(content), "pending log message")
			break
		}
	}
}

func Test_TraceDepth(t *testing.T) {
	t.Run("trace enabled", func(t *testing.T) {
		tmpDir := t.TempDir()
		l := logger.Create(logger.Options{
			Path:       tmpDir,
			Max:        100,
			TraceDepth: 1, // Enable trace
		})

		l.Debug("testing trace enabled")
		time.Sleep(100 * time.Millisecond)
		l.Close()

		// Verify log file was created with trace info
		files, err := os.ReadDir(tmpDir)
		require.NoError(t, err)
		require.NotEmpty(t, files)

		// Read file and verify trace is present
		for _, file := range files {
			if strings.Contains(file.Name(), "debug") {
				content, err := os.ReadFile(filepath.Join(tmpDir, file.Name()))
				require.NoError(t, err)
				require.Contains(t, string(content), "[trace=")
				break
			}
		}
	})

	t.Run("trace disabled", func(t *testing.T) {
		tmpDir := t.TempDir()
		l := logger.Create(logger.Options{
			Path:       tmpDir,
			Max:        100,
			TraceDepth: 0, // Disable trace
		})

		l.Debug("testing trace disabled")
		time.Sleep(100 * time.Millisecond)
		l.Close()

		// Verify log file was created WITHOUT trace info
		files, err := os.ReadDir(tmpDir)
		require.NoError(t, err)
		require.NotEmpty(t, files)

		// Read file and verify trace is NOT present
		for _, file := range files {
			if strings.Contains(file.Name(), "debug") {
				content, err := os.ReadFile(filepath.Join(tmpDir, file.Name()))
				require.NoError(t, err)
				require.NotContains(t, string(content), "[trace=")
				break
			}
		}
	})
}

func Test_DefaultOptions(t *testing.T) {
	// Test with empty options - should use defaults
	l := logger.Create(logger.Options{})
	defer l.Close()

	require.NotPanics(t, func() {
		l.Info("testing default options")
	})
	time.Sleep(100 * time.Millisecond)
}

func Test_LogAndLogfMethods(t *testing.T) {
	tmpDir := t.TempDir()
	l := logger.Create(logger.Options{
		Path: tmpDir,
		Max:  100,
	})
	defer l.Close()

	// Test Log method with all levels
	l.Log(logger.LevelDebug, "log debug")
	l.Log(logger.LevelInfo, "log info")
	l.Log(logger.LevelWarn, "log warn")
	l.Log(logger.LevelError, "log error")
	l.Log(logger.LevelFatal, "log fatal")

	// Test Logf method with all levels
	l.Logf(logger.LevelDebug, "logf debug %d", 1)
	l.Logf(logger.LevelInfo, "logf info %d", 2)
	l.Logf(logger.LevelWarn, "logf warn %d", 3)
	l.Logf(logger.LevelError, "logf error %d", 4)
	l.Logf(logger.LevelFatal, "logf fatal %d", 5)

	time.Sleep(100 * time.Millisecond)

	// Verify log files were created
	files, err := os.ReadDir(tmpDir)
	require.NoError(t, err)
	require.NotEmpty(t, files)
}

func Test_GetLevelNameEdgeCases(t *testing.T) {
	tests := []struct {
		level    logger.Level
		expected string
	}{
		{logger.LevelDebug, "debug"},
		{logger.LevelInfo, "info"},
		{logger.LevelWarn, "warn"},
		{logger.LevelError, "error"},
		{logger.LevelFatal, "fatal"},
		{logger.Level(999), ""}, // Unknown level
		{logger.Level(-1), ""},  // Negative level
	}

	for _, tt := range tests {
		result := logger.GetLevelName(tt.level)
		require.Equal(t, tt.expected, result, "level %d", tt.level)
	}
}

func Test_ExtractAllContentEdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"no placeholders", nil},
		{"${single}", []string{"single"}},
		{"${first} ${second}", []string{"first", "second"}},
		{"${}", []string{""}},                                   // Empty placeholder
		{"${ spaced }", []string{" spaced "}},                   // Placeholder with spaces
		{"${a}${b}${c}", []string{"a", "b", "c"}},               // Consecutive
		{"text ${var} more text", []string{"var"}},              // Mixed content
		{"${123}", []string{"123"}},                             // Numeric content
		{"${special-chars_123}", []string{"special-chars_123"}}, // Special chars
	}

	for _, tt := range tests {
		result := logger.ExtractAllContent(tt.input)
		require.Equal(t, tt.expected, result, "input: %s", tt.input)
	}
}

func Test_MetadataAppending(t *testing.T) {
	tmpDir := t.TempDir()
	l := logger.Create(logger.Options{
		Path: tmpDir,
		Max:  100,
		Metadata: logger.Metadata{
			"global1": "value1",
			"global2": "value2",
		},
	})
	defer l.Close()

	// Test with additional metadata
	l.Info("with extra metadata", logger.Metadata{
		"local1": "localValue1",
	})

	// Test with multiple metadata maps
	l.Info("with multiple metadata", logger.Metadata{
		"extra1": "extraValue1",
	}, logger.Metadata{
		"extra2": "extraValue2",
	})

	// Test with empty metadata
	l.Info("with empty metadata", logger.Metadata{})

	// Test with nil slice (no metadata)
	l.Info("with no metadata")

	time.Sleep(100 * time.Millisecond)
}

func Test_ConcurrentLogging(t *testing.T) {
	tmpDir := t.TempDir()
	l := logger.Create(logger.Options{
		Path: tmpDir,
		Max:  100,
	})
	defer l.Close()

	// Start multiple goroutines logging concurrently
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				l.Infof("goroutine %d message %d", id, j)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	time.Sleep(500 * time.Millisecond)
}

func Test_LoggerWithEmptyMetadata(t *testing.T) {
	tmpDir := t.TempDir()
	l := logger.Create(logger.Options{
		Path:     tmpDir,
		Max:      100,
		Metadata: logger.Metadata{}, // Empty global metadata
	})
	defer l.Close()

	l.Info("message with empty global metadata")
	time.Sleep(100 * time.Millisecond)
}

func Test_LoggerWithNilMetadata(t *testing.T) {
	tmpDir := t.TempDir()
	l := logger.Create(logger.Options{
		Path:     tmpDir,
		Max:      100,
		Metadata: nil, // Nil global metadata
	})
	defer l.Close()

	l.Info("message with nil global metadata")
	time.Sleep(100 * time.Millisecond)
}
