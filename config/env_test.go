package config

import (
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

type Config struct {
	NodeEnv   string        `mapstructure:"NODE_ENV"`
	Port      int           `mapstructure:"PORT"`
	ExpiresIn time.Duration `mapstructure:"EXPIRES_IN"`
	Log       bool          `mapstructure:"LOG"`
	Special   interface{}   `mapstructure:"SPECIAL"`
	Secret    string        `mapstructure:"SECRET"`
}

func Test_Scan(t *testing.T) {
	err := godotenv.Load(".env.example")
	require.Nil(t, err)
	var cfg Config
	Scan(&cfg)

	require.Equal(t, "development", cfg.NodeEnv)
	require.Equal(t, 5000, cfg.Port)
	require.Equal(t, 5*time.Minute, cfg.ExpiresIn)
	require.Equal(t, false, cfg.Log)
	require.Equal(t, "", cfg.Secret)
}

func Test_New(t *testing.T) {
	_, err := New[Config](".env.example")
	require.Nil(t, err)

	var cfg Config
	Scan(&cfg)
	require.Equal(t, "development", cfg.NodeEnv)
	require.Equal(t, 5000, cfg.Port)
	require.Equal(t, 5*time.Minute, cfg.ExpiresIn)
}

func Test_GetRaw(t *testing.T) {
	_, err := New[Config](".env.example")
	require.Nil(t, err)

	dev := GetRaw("NODE_ENV")
	require.Equal(t, "development", dev)
}
