package app

import (
	"fmt"
	"os"

	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/database/sql"
	"github.com/tinh-tinh/tinhtinh/example/app/user"
	"gorm.io/driver/postgres"
)

// func NewModule() *api.Module {
// 	appModule := api.NewModule(api.NewModuleOptions{
// 		Import: []api.ModuleParam{user.Module},
// 	})

// 	return appModule
// }

func NewModule() *core.DynamicModule {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	appModule := core.NewModule(core.NewModuleOptions{
		Global: true,
		Imports: []core.Module{
			sql.ForRoot(sql.ConnectionOptions{
				Dialect: postgres.Open(dsn),
			}),
			user.Module,
		},
	})

	return appModule
}
