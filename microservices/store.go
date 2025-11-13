package microservices

import "github.com/tinh-tinh/tinhtinh/v2/core"

const STORE core.Provide = "STORE"

type Store struct {
	Subscribers []*SubscribeHandler
	Rpcs        []RpcHandler
}

func Register(transports ...string) core.Modules {
	return func(module core.Module) core.Module {
		handlerModule := module.New(core.NewModuleOptions{})

		handlerModule.NewProvider(core.ProviderOptions{
			Name:  STORE,
			Value: &Store{},
		})
		handlerModule.Export(STORE)

		if len(transports) > 0 {
			for _, transport := range transports {
				name := ToTransport(transport)
				handlerModule.NewProvider(core.ProviderOptions{
					Name:  name,
					Value: &Store{},
				})
				handlerModule.Export(name)
			}
		}

		return handlerModule
	}
}

func ToTransport(transport string) core.Provide {
	return STORE + core.Provide(transport)
}
