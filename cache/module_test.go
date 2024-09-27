package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/core"
)

func Test_Register(t *testing.T) {
	appModule := core.NewModule(core.NewModuleOptions{
		Imports: []core.Module{Register(Options{
			Ttl: 5 * time.Minute,
			Max: 100,
		})},
	})

	cache := appModule.Ref(CACHE_MANAGER).(Store)

	cache.Set("john", "doe")
	require.Equal(t, 1, cache.Count())
	require.Equal(t, "doe", cache.Get("john"))

	cache.Set("alice", "doe", 0)
	require.Nil(t, cache.Get("alice"))
}
