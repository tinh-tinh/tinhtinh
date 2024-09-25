package abc

import (
	"fmt"

	"github.com/tinh-tinh/tinhtinh/core"
)

func NewModule(module *core.DynamicModule) *core.DynamicModule {
	abcModule := module.New(core.NewModuleOptions{
		Scope: core.Request,
	})

	abcModule.NewProvider(core.ProviderOptions{
		Name: "jajaj",
		Factory: func(param ...interface{}) interface{} {
			return fmt.Sprintf("%vAbc", param[0])
		},
		Inject: []core.Provide{core.REQUEST},
	})
	abcModule.Export("jajaj")

	return abcModule
}
