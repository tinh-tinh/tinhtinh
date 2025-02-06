package microservices

import "github.com/tinh-tinh/tinhtinh/v2/core"

const STORE core.Provide = "STORE"

type Store struct {
	Subscribers map[string][]*SubscribeHandler
}

func Register() core.Modules {
	return func(module core.Module) core.Module {
		handlerModule := module.New(core.NewModuleOptions{})

		handlerModule.NewProvider(core.ProviderOptions{
			Name:  STORE,
			Value: &Store{Subscribers: make(map[string][]*SubscribeHandler)},
		})

		handlerModule.Export(STORE)
		return handlerModule
	}
}

func (store *Store) GetRPC() []*SubscribeHandler {
	return store.Subscribers[string(RPC)]
}

func (store *Store) GetPubSub() []*SubscribeHandler {
	return store.Subscribers[string(PubSub)]
}
