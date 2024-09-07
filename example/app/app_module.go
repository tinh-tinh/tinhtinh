package app

import (
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/example/app/post"
)

type Config struct {
	Port    int    `mapstructure:"PORT"`
	NodeEnv string `mapstructure:"NODE_ENV"`

	DBHost string `mapstructure:"DB_HOST"`
	DBPort int    `mapstructure:"DB_PORT"`
	DBUser string `mapstructure:"DB_USER"`
	DBPass string `mapstructure:"DB_PASS"`
	DBName string `mapstructure:"DB_NAME"`
}

func NewModule() *core.DynamicModule {
	appModule := core.NewModule(core.NewModuleOptions{
		Global: true,
		Imports: []core.Module{
			// config.ForRoot[Config](),
			// sql.ForRoot(sql.ConnectionOptions{
			// 	Factory: func(module *core.DynamicModule) gorm.Dialector {
			// 		env := module.Ref(config.ENV).(Config)
			// 		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", env.DBHost, env.DBUser, env.DBPass, env.DBName, env.DBPort)

			// 		return postgres.Open(dsn)
			// 	},
			// }),
			// token.Register(token.Options{
			// 	Alg:    jwt.SigningMethodHS256,
			// 	Secret: "1234567890krj3k4brub45uybf874847g2f345uy",
			// 	Exp:    time.Hour,
			// }),
			// user.Module,
			post.Module,
		},
	})

	return appModule
}
