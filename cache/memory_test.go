package cache

import (
	"fmt"
	"testing"
	"time"
)

func Test_NewInMemory(t *testing.T) {
	memory := NewInMemory[string, int](MemoryOptions{
		Max: 10,
		Ttl: 10 * time.Second,
	})

	for i := 0; i < 13; i++ {
		key := fmt.Sprint(i)
		memory.Set(key, i)
	}

	t.Run("Test_Get", func(t *testing.T) {
		ab := memory.Get("10")
		if *ab != 10 {
			t.Error("ab should equal 10")
		}
	})
}
