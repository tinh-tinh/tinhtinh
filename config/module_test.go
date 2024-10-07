package config

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/core"
)

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
