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
}

func Test_Scan(t *testing.T) {
	err := godotenv.Load(".env.example")
	require.Nil(t, err)
	var cfg Config
	Scan(&cfg)

	require.Equal(t, "development", cfg.NodeEnv)
	require.Equal(t, 5000, cfg.Port)
	require.Equal(t, 5*time.Minute, cfg.ExpiresIn)
}
