package app

import (
	"github.com/tinh-tinh/tinhtinh/core"
)

// type Config struct {
// 	Port    int    `mapstructure:"PORT"`
// 	NodeEnv string `mapstructure:"NODE_ENV"`

// 	DBHost string `mapstructure:"DB_HOST"`
// 	DBPort int    `mapstructure:"DB_PORT"`
// 	DBUser string `mapstructure:"DB_USER"`
// 	DBPass string `mapstructure:"DB_PASS"`
// 	DBName string `mapstructure:"DB_NAME"`
// }

func NewModule() *core.DynamicModule {
	appModule := core.NewModule(core.NewModuleOptions{
		Global: true,
	})

	return appModule
}
