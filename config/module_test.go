package config

import "testing"

type TestEnv struct {
	Port int `mapstructure:"PORT"`
}

func Test_Module(t *testing.T) {
	Register[TestEnv]("")

	t.Run("Test Env", func(t *testing.T) {
		if Get[TestEnv]().Port != 5000 {
			t.Error("expect 5000, but got", Get[TestEnv]().Port)
		}
	})
}
