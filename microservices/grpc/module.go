package grpc

import (
	"github.com/tinh-tinh/tinhtinh/microservices"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func Register() core.Modules {
	return func(module core.Module) core.Module {
		return module.New(core.NewModuleOptions{
			Imports: []core.Modules{microservices.Register()},
		})
	}
}
