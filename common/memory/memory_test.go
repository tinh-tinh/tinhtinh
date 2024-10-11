package memory

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_NewInMemory(t *testing.T) {
	t.Parallel()

	store := New(Options{
		Max: 100,
		Ttl: 1 * time.Hour,
	})

	var (
		key = "john-internal"
		val = []byte("doe")
		exp = 1 * time.Second
	)
	// Set key with value
	store.Set(key, val)
	result := store.Get(key)
	require.Equal(t, val, result)

	// Get non-existing key
	result = store.Get("empty")
	require.Nil(t, result)

	// Set key with value and ttl
	store.Set(key, val, exp)
	time.Sleep(1100 * time.Millisecond)
	result = store.Get(key)
	require.Nil(t, result)

	// Set key with value and no expiration
	store.Set(key, val)
	result = store.Get(key)
	require.Equal(t, val, result)

	// Delete key
	store.Delete(key)
	result = store.Get(key)
	require.Nil(t, result)

	// Reset all keys
	store.Set("john-reset", val, 0)
	store.Set("doe-reset", val, 0)
	store.Clear()

	// Check if all keys are deleted
	result = store.Get("john-reset")
	require.Nil(t, result)
	result = store.Get("doe-reset")
	require.Nil(t, result)

	// Count the number of keys
	count := store.Count()
	require.Equal(t, 0, count)

	// Check if key exists
	store.Set(key, val)
	require.True(t, store.Has(key))
	store.Delete(key)
	require.False(t, store.Has(key))

	// Get all keys
	store.Set(key, val)
	require.Len(t, store.Keys(), 1)
	store.Clear()
	require.Len(t, store.Keys(), 0)
}

func Test_Options(t *testing.T) {

	store := New(Options{
		Max: 10,
		Ttl: 1 * time.Hour,
	})

	for i := 0; i < 15; i++ {
		store.Set(fmt.Sprintf("john-%d", i), []byte("doe"))
	}

	require.Equal(t, 10, store.Count())
}
