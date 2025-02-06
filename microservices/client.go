package microservices

import (
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

const CLIENT core.Provide = "CLIENT"

func RegisterClient(client ClientProxy) core.Modules {
	return func(module core.Module) core.Module {
		clientModule := module.New(core.NewModuleOptions{})

		clientModule.NewProvider(core.ProviderOptions{
			Name:  CLIENT,
			Value: client,
		})

		clientModule.Export(CLIENT)
		return clientModule
	}
}

func Inject(module core.Module) ClientProxy {
	conn, ok := module.Ref(CLIENT).(ClientProxy)
	if !ok {
		return nil
	}
	return conn
}
