package config

import (
	"fmt"
	"testing"
)

type Config struct {
	NodeEnv string `mapstructure:"NODE_ENV"`
	Port    string `mapstructure:"PORT"`
}

func Test_Scan(t *testing.T) {
	t.Run("test case", func(t *testing.T) {
		mapper := Scan(&Config{})
		fmt.Println(mapper)
	})
}
