package cache

// import (
// 	"fmt"
// 	"testing"
// 	"time"
// )

// func Test_NewInMemory(t *testing.T) {
// 	memory := NewInMemory[string, int](MemoryOptions{
// 		Max: 10,
// 		Ttl: 10 * time.Second,
// 	})

// 	for i := 0; i < 13; i++ {
// 		key := fmt.Sprint(i)
// 		time.Sleep(10 * time.Millisecond)
// 		memory.Set(key, i)
// 	}

// 	t.Run("Test_Get", func(t *testing.T) {
// 		ab := memory.Get("10")
// 		fmt.Printf("ab is %v", ab)
// 		if *ab != 10 {
// 			t.Error("ab should equal 10")
// 		}

// 		first := memory.Get("1")
// 		if first != nil {
// 			t.Error("first should out of cache")
// 		}
// 	})

// 	time.Sleep(15 * time.Second)
// 	t.Run("Test_Get", func(t *testing.T) {
// 		ab := memory.Get("12")
// 		if ab != nil {
// 			t.Error("ab should out of cache")
// 		}
// 	})
// }
