package cache

import (
	"testing"
	"time"
)

func Test_Set(t *testing.T) {
	Register()
	t.Run("test case set", func(t *testing.T) {
		Set("key", "value", 15*time.Second)
		val := Get("key")
		if val != "value" {
			t.Error("expect value, but got", val)
		}
		time.Sleep(15 * time.Second)
		val = Get("key")
		if val != nil {
			t.Error("expect nil, but got", val)
		}
	})
}
