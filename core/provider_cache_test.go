package core

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test structs
type TestProvider struct {
	Value string
}

type CustomNameProvider struct {
	Value string
}

func (CustomNameProvider) ProvideName() string {
	return "CustomName"
}

func TestGetProvideNameCaching(t *testing.T) {
	// Clear cache before test
	providerNameCache = sync.Map{}

	t.Run("should cache struct name", func(t *testing.T) {
		provider1 := &TestProvider{Value: "test1"}
		name1 := getProvideName(provider1)
		require.Equal(t, "TestProvider", name1)

		// Second call should hit cache
		provider2 := &TestProvider{Value: "test2"}
		name2 := getProvideName(provider2)
		require.Equal(t, "TestProvider", name2)
		require.Equal(t, name1, name2)
	})

	t.Run("should cache custom ProvideName", func(t *testing.T) {
		provider1 := &CustomNameProvider{Value: "test1"}
		name1 := getProvideName(provider1)
		require.Equal(t, "CustomName", name1)

		// Second call should hit cache
		provider2 := &CustomNameProvider{Value: "test2"}
		name2 := getProvideName(provider2)
		require.Equal(t, "CustomName", name2)
		require.Equal(t, name1, name2)
	})

	t.Run("should handle different types separately", func(t *testing.T) {
		testProvider := &TestProvider{Value: "test"}
		customProvider := &CustomNameProvider{Value: "custom"}

		testName := getProvideName(testProvider)
		customName := getProvideName(customProvider)

		require.Equal(t, "TestProvider", testName)
		require.Equal(t, "CustomName", customName)
		require.NotEqual(t, testName, customName)
	})
}

func TestGetProvideNameConcurrency(t *testing.T) {
	// Clear cache before test
	providerNameCache = sync.Map{}

	t.Run("should be thread-safe", func(t *testing.T) {
		var wg sync.WaitGroup
		iterations := 100

		// Run concurrent calls
		for i := 0; i < iterations; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				provider := &TestProvider{Value: "concurrent"}
				name := getProvideName(provider)
				assert.Equal(t, "TestProvider", name)
			}()
		}

		wg.Wait()
	})

	t.Run("should handle multiple types concurrently", func(t *testing.T) {
		var wg sync.WaitGroup
		iterations := 50

		for i := 0; i < iterations; i++ {
			wg.Add(2)

			go func() {
				defer wg.Done()
				provider := &TestProvider{Value: "test"}
				name := getProvideName(provider)
				assert.Equal(t, "TestProvider", name)
			}()

			go func() {
				defer wg.Done()
				provider := &CustomNameProvider{Value: "custom"}
				name := getProvideName(provider)
				assert.Equal(t, "CustomName", name)
			}()
		}

		wg.Wait()
	})
}

func BenchmarkGetProvideNameWithCache(b *testing.B) {
	// Clear cache before benchmark
	providerNameCache = sync.Map{}

	provider := &TestProvider{Value: "benchmark"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getProvideName(provider)
	}
}

func BenchmarkGetProvideNameWithCacheParallel(b *testing.B) {
	// Clear cache before benchmark
	providerNameCache = sync.Map{}

	provider := &TestProvider{Value: "benchmark"}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = getProvideName(provider)
		}
	})
}

func BenchmarkGetProvideNameCustomName(b *testing.B) {
	// Clear cache before benchmark
	providerNameCache = sync.Map{}

	provider := &CustomNameProvider{Value: "benchmark"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getProvideName(provider)
	}
}
