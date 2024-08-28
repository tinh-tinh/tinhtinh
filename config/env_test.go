package config

import (
	"fmt"
	"testing"
	"time"
)

type Config struct {
	NodeEnv   string        `mapstructure:"NODE_ENV"`
	Port      int           `mapstructure:"PORT"`
	ExpiresIn time.Duration `mapstructure:"EXPIRES_IN"`
}

func Test_Scan(t *testing.T) {
	Register[Config](".env")
	t.Run("test case", func(t *testing.T) {
		var cfg Config
		Scan(&cfg)
		fmt.Printf("Config is %v", cfg)
	})
}
