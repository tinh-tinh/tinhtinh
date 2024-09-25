package user

import (
	"fmt"

	"github.com/tinh-tinh/tinhtinh/core"
)

func UserProvider(module *core.DynamicModule) *core.DynamicProvider {
	provider := module.NewProvider(core.ProviderOptions{
		Name: "user",
		Factory: func(param ...interface{}) interface{} {
			fmt.Println("param", param[1])
			return fmt.Sprintf("Root%vUser", param[1])
		},
		Inject: []core.Provide{core.REQUEST, "root", "jajaj"},
	})

	return provider
}
