package sql

import (
	"fmt"

	"github.com/tinh-tinh/tinhtinh/core"
	"gorm.io/gorm"
)

type ConnectionOptions struct {
	Dialect gorm.Dialector
	Factory func(module *core.DynamicModule) gorm.Dialector
}

const ConnectDB core.Provide = "ConnectDB"

func ForRoot(opt ConnectionOptions) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		var dialector gorm.Dialector
		if opt.Factory != nil {
			dialector = opt.Factory(module)
		} else {
			dialector = opt.Dialect
		}
		conn, err := gorm.Open(dialector, &gorm.Config{})
		if err != nil {
			panic(err)
		}
		conn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
		fmt.Println("connected to database")

		sqlModule := module.New(core.NewModuleOptions{})
		sqlModule.NewProvider(conn, ConnectDB)
		sqlModule.Export(ConnectDB)

		return sqlModule
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

func InjectGorm(module *core.DynamicModule) *gorm.DB {
	return module.Ref(ConnectDB).(*gorm.DB)
}
