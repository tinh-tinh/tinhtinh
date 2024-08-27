package config

import "testing"

type TestEnv struct {
	Port string `mapstructure:"PORT"`
}

func Test_Module(t *testing.T) {
	configModule := New[TestEnv]("")

	t.Run("Test Env", func(t *testing.T) {
		if configModule.Get("Port") == 5000 {
			t.Error("error")
		}
	})
}
