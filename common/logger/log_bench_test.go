package logger_test

import (
	"testing"

	"github.com/tinh-tinh/tinhtinh/v2/common/logger"
)

func BenchmarkLoggerInfo(b *testing.B) {
	tmpDir := b.TempDir()
	l := logger.Create(logger.Options{
		Path: tmpDir,
		Max:  100,
	})
	defer l.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Info("benchmark message")
	}
}

func BenchmarkLoggerInfof(b *testing.B) {
	tmpDir := b.TempDir()
	l := logger.Create(logger.Options{
		Path: tmpDir,
		Max:  100,
	})
	defer l.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Infof("benchmark message %d", i)
	}
}

func BenchmarkLoggerDebug(b *testing.B) {
	tmpDir := b.TempDir()
	l := logger.Create(logger.Options{
		Path: tmpDir,
		Max:  100,
	})
	defer l.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Debug("benchmark debug message")
	}
}

func BenchmarkLoggerWarn(b *testing.B) {
	tmpDir := b.TempDir()
	l := logger.Create(logger.Options{
		Path: tmpDir,
		Max:  100,
	})
	defer l.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Warn("benchmark warn message")
	}
}

func BenchmarkLoggerError(b *testing.B) {
	tmpDir := b.TempDir()
	l := logger.Create(logger.Options{
		Path: tmpDir,
		Max:  100,
	})
	defer l.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Error("benchmark error message")
	}
}

func BenchmarkLoggerWithMetadata(b *testing.B) {
	tmpDir := b.TempDir()
	l := logger.Create(logger.Options{
		Path: tmpDir,
		Max:  100,
	})
	defer l.Close()

	meta := logger.Metadata{
		"key1": "value1",
		"key2": 123,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Info("benchmark with metadata", meta)
	}
}

func BenchmarkLoggerWithGlobalMetadata(b *testing.B) {
	tmpDir := b.TempDir()
	l := logger.Create(logger.Options{
		Path: tmpDir,
		Max:  100,
		Metadata: logger.Metadata{
			"service": "benchmark-svc",
			"version": "1.0.0",
		},
	})
	defer l.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Info("benchmark with global metadata")
	}
}

func BenchmarkLoggerWithTraceDepth(b *testing.B) {
	tmpDir := b.TempDir()
	l := logger.Create(logger.Options{
		Path:       tmpDir,
		Max:        100,
		TraceDepth: 1, // Enable trace
	})
	defer l.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Debug("benchmark with trace")
	}
}

func BenchmarkLoggerParallel(b *testing.B) {
	tmpDir := b.TempDir()
	l := logger.Create(logger.Options{
		Path: tmpDir,
		Max:  100,
	})
	defer l.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Info("parallel benchmark message")
		}
	})
}

func BenchmarkLoggerParallelWithMetadata(b *testing.B) {
	tmpDir := b.TempDir()
	l := logger.Create(logger.Options{
		Path: tmpDir,
		Max:  100,
	})
	defer l.Close()

	meta := logger.Metadata{
		"key": "value",
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Info("parallel benchmark with metadata", meta)
		}
	})
}

func BenchmarkGetLevelName(b *testing.B) {
	levels := []logger.Level{
		logger.LevelDebug,
		logger.LevelInfo,
		logger.LevelWarn,
		logger.LevelError,
		logger.LevelFatal,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, level := range levels {
			_ = logger.GetLevelName(level)
		}
	}
}

func BenchmarkExtractAllContent(b *testing.B) {
	testStr := "Hello ${name}, your order ${orderId} is ready"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = logger.ExtractAllContent(testStr)
	}
}

func BenchmarkExtractAllContentComplex(b *testing.B) {
	testStr := "${var1} some text ${var2} more text ${var3} even more ${var4} and ${var5}"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = logger.ExtractAllContent(testStr)
	}
}

func BenchmarkLoggerMixedLevels(b *testing.B) {
	tmpDir := b.TempDir()
	l := logger.Create(logger.Options{
		Path: tmpDir,
		Max:  100,
	})
	defer l.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		switch i % 5 {
		case 0:
			l.Debug("debug message")
		case 1:
			l.Info("info message")
		case 2:
			l.Warn("warn message")
		case 3:
			l.Error("error message")
		case 4:
			l.Fatal("fatal message")
		}
	}
}
