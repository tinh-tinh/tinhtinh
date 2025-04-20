package microservices

import (
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

type ClientOptions struct {
	Name      core.Provide
	Transport ClientProxy
}

func RegisterClient(options ...ClientOptions) core.Modules {
	return func(module core.Module) core.Module {
		clientModule := module.New(core.NewModuleOptions{})

		for _, option := range options {
			clientModule.NewProvider(core.ProviderOptions{
				Name:  option.Name,
				Value: option.Transport,
			})

			clientModule.Export(option.Name)
		}

		return clientModule
	}
}

type ClientFactory func(ref core.RefProvider) []ClientOptions

func RegisterClientFactory(factory ClientFactory) core.Modules {
	return func(module core.Module) core.Module {
		options := factory(module)
		clientModule := module.New(core.NewModuleOptions{})

		for _, option := range options {
			clientModule.NewProvider(core.ProviderOptions{
				Name:  option.Name,
				Value: option.Transport,
			})

			clientModule.Export(option.Name)
		}

		return clientModule
	}
}

func InjectClient(ref core.RefProvider, name core.Provide) ClientProxy {
	conn, ok := ref.Ref(name).(ClientProxy)
	if !ok {
		return nil
	}
	return conn
}
