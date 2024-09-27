package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/core"
)

func Test_Register(t *testing.T) {
	_, err := Register[Config](".env.example")
	require.Nil(t, err)

	var cfg Config
	Scan(&cfg)
	require.Equal(t, "development", cfg.NodeEnv)
	require.Equal(t, 5000, cfg.Port)
	require.Equal(t, 5*time.Minute, cfg.ExpiresIn)
}

func Test_GetRaw(t *testing.T) {
	_, err := Register[Config](".env.example")
	require.Nil(t, err)

	dev := GetRaw("NODE_ENV")
	require.Equal(t, "development", dev)
}

func Test_ForRootNil(t *testing.T) {
	appModule := core.NewModule(core.NewModuleOptions{
		Imports: []core.Module{ForRoot[Config]("")},
	})

	cfg, ok := appModule.Ref(ENV).(*Config)
	require.False(t, ok)
	require.Nil(t, cfg)
}

func Test_ForRoot(t *testing.T) {
	appModule := core.NewModule(core.NewModuleOptions{
		Imports: []core.Module{ForRoot[Config](".env.example")},
	})

	cfg, ok := appModule.Ref(ENV).(*Config)
	require.True(t, ok)
	require.NotNil(t, cfg)
	require.Equal(t, "development", cfg.NodeEnv)
}
