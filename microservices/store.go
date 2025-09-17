package microservices

import "github.com/tinh-tinh/tinhtinh/v2/core"

const STORE core.Provide = "STORE"

type Store struct {
	Subscribers map[EventType][]*SubscribeHandler
}

func Register(transports ...string) core.Modules {
	return func(module core.Module) core.Module {
		handlerModule := module.New(core.NewModuleOptions{})

		handlerModule.NewProvider(core.ProviderOptions{
			Name:  STORE,
			Value: &Store{Subscribers: make(map[EventType][]*SubscribeHandler)},
		})
		handlerModule.Export(STORE)

		if len(transports) > 0 {
			for _, transport := range transports {
				name := ToTransport(transport)
				handlerModule.NewProvider(core.ProviderOptions{
					Name:  name,
					Value: &Store{Subscribers: make(map[EventType][]*SubscribeHandler)},
				})
				handlerModule.Export(name)
			}
		}

		return handlerModule
	}
}

func (store *Store) GetRPC() []*SubscribeHandler {
	return store.Subscribers[RPC]
}

func (store *Store) GetPubSub() []*SubscribeHandler {
	return store.Subscribers[PubSub]
}

func ToTransport(transport string) core.Provide {
	return STORE + core.Provide(transport)
}
