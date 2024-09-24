package logger

import (
	"fmt"
	"testing"
)

func Test_Create(t *testing.T) {
	log := Create(Options{
		Max:    1,
		Rotate: true,
	})
	for i := 0; i < 100000; i++ {
		log.Info(fmt.Sprint(i))
	}
}
