package logger_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/logger"
)

type serviceImpl struct{}

func (s *serviceImpl) DoSomething(log *logger.Logger) {
	log.Debug("Doing something in Service")
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

	loggerAnother := logger.Create(logger.Options{
		Max:        1,
		Rotate:     true,
		TraceDepth: 5,
	})
	svc.DoSomething(loggerAnother)
}

func Test_Create(t *testing.T) {
	log := logger.Create(logger.Options{
		Max:    1,
		Rotate: true,
	})
	for i := range 1000 {
		val := strconv.Itoa(i)
		if i%2 == 0 {
			log.Info(val)
		} else if i%3 == 0 {
			log.Warn(val)
		} else if i%5 == 0 {
			log.Error(val)
		} else if i%7 == 0 {
			log.Fatal(val)
		} else {
			log.Debug(val)
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

	log = logger.Create(logger.Options{
		Max:    1,
		Rotate: true,
	})
	for i := range 100 {
		if i%2 == 0 {
			log.Infof("The value is %d", i)
		} else if i%3 == 0 {
			log.Warnf("The value is %d", i)
		} else if i%5 == 0 {
			log.Errorf("The value is %d", i)
		} else if i%7 == 0 {
			log.Fatalf("The value is %d", i)
		} else {
			log.Debugf("The value is %d", i)
		}
		log.Logf(logger.LevelDebug, "alayws have ata %d", i)
	}

	for i := range 1000 {
		val := strconv.Itoa(i)
		if i%2 == 0 {
			log.Info(val, logger.Metadata{
				"function": "Test",
			})
		} else if i%3 == 0 {
			log.Warn(val, logger.Metadata{
				"function": "Test",
			})
		} else if i%5 == 0 {
			log.Error(val, logger.Metadata{
				"function": "Test",
			})
		} else if i%7 == 0 {
			log.Fatal(val, logger.Metadata{
				"function": "Test",
			})
		} else {
			log.Debug(val, logger.Metadata{
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
