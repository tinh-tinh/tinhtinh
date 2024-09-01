package sql

import (
	"fmt"

	"github.com/tinh-tinh/tinhtinh/core"
	"gorm.io/gorm"
)

type ConnectionOptions struct {
	Dialect gorm.Dialector
}

const ConnectDB core.Provide = "ConnectDB"

func ForRoot(opt ConnectionOptions) core.Module {
	conn, err := gorm.Open(opt.Dialect, &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	fmt.Println("connected to database")

	dbModule := core.NewModule(core.NewModuleOptions{
		Global: true,
		Providers: []core.Provider{func(module *core.DynamicModule) *core.DynamicProvider {
			provider := core.NewProvider(module)
			provider.Set(ConnectDB, conn)

			return provider
		}},
	})

	return func(module *core.DynamicModule) *core.DynamicModule {
		return dbModule
	}
}

func ForFeature(models ...interface{}) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		conn := module.Ref(ConnectDB).(*gorm.DB)
		fmt.Println("Migrating...")
		err := conn.AutoMigrate(models...)
		if err != nil {
			panic(err)
		}
		fmt.Println("Migrated successful")
		return module
	}
}
