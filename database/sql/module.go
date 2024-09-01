package sql

import (
	"fmt"
	"reflect"

	"github.com/tinh-tinh/tinhtinh/core"
	"gorm.io/gorm"
)

const (
	ConnectDB core.Provide = "ConnectDB"
)

type RegistryOptions struct {
	Dialect gorm.Dialector
	Models  []interface{}
}

func Registry(opt RegistryOptions) core.Module {
	conn, err := gorm.Open(opt.Dialect, &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	fmt.Println("connected to database migrating...")
	err = conn.AutoMigrate(opt.Models...)
	if err != nil {
		panic(err)
	}
	fmt.Println("Migrated successful")
	provider := core.NewProvider(ConnectDB, conn)

	dbModule := core.NewModule(core.NewModuleOptions{
		Global: true,
		Providers: []core.Provider{func(module *core.DynamicModule) *core.DynamicProvider {
			return provider
		}},
	})

	return func(module *core.DynamicModule) *core.DynamicModule {
		return dbModule
	}
}

func GetModel[M any](name ...string) core.Provider {
	var model M
	var provide string

	if len(name) == 0 {
		provide = GetStructName(model)
	} else {
		provide = name[0]
	}

	return func(module *core.DynamicModule) *core.DynamicProvider {
		conn := module.Ref(ConnectDB).(*gorm.DB)
		return core.NewProvider(core.Provide(provide), conn.Model(&model))
	}
}

func RegistryModel[M any](name ...string) core.Provider {
	var model M
	var provide string

	if len(name) == 0 {
		provide = GetStructName(model)
	} else {
		provide = name[0]
	}

	return func(module *core.DynamicModule) *core.DynamicProvider {
		conn := module.Ref(ConnectDB).(*gorm.DB)
		return core.NewProvider(core.Provide(provide), conn.Model(&model))
	}
}

func GetStructName(val interface{}) string {
	if t := reflect.TypeOf(val); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}
