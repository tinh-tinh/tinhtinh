package core

import (
	"github.com/tinh-tinh/tinhtinh/utils"
)

type Provide string
type DynamicProvider struct {
	Name    Provide
	Value   interface{}
	Factory ReqProvider
}

func (module *DynamicModule) NewProvider(val interface{}, name ...Provide) *DynamicProvider {
	if val == nil {
		return nil
	}
	var providerName Provide
	if len(name) > 0 {
		providerName = name[0]
	} else {
		providerName = Provide(utils.GetNameStruct(val))
	}
	provider := &DynamicProvider{
		Value: val,
		Name:  providerName,
	}

	module.providers = append(module.providers, provider)
	return provider
}

type ReqProvider func(ctx Ctx) interface{}

func (module *DynamicModule) NewReqProvider(name string, reqFnc ReqProvider) *DynamicProvider {
	provider := &DynamicProvider{
		Name:    Provide(name),
		Factory: reqFnc,
	}

	module.providers = append(module.providers, provider)
	return provider
}
