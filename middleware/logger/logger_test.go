package logger_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/middleware/logger"
)

func Test_Create(t *testing.T) {
	log := logger.Create(logger.Options{
		Max:    1,
		Rotate: true,
	})
	for i := 0; i < 1000; i++ {
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

	require.Panics(t, func() {
		for i := 0; i < 2; i++ {
			log2.Info(randomBigStr())
		}
	})

	log3 := logger.Create(logger.Options{
		Path:   "logs/test2",
		Max:    1,
		Rotate: true,
	})

	for i := 0; i < 2; i++ {
		log3.Info(randomBigStr())
	}
}

func randomBigStr() string {
	var bigString strings.Builder
	// Define the number of repetitions
	repeat := 1000000
	smallString := "Hello, Go! "

	// Append the small string multiple times
	for i := 0; i < repeat; i++ {
		bigString.WriteString(smallString)
	}

	// Convert the builder to a string
	result := bigString.String()
	return result
}
