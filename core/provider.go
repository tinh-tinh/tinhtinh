package core

import (
	"slices"

	"github.com/tinh-tinh/tinhtinh/utils"
)

type Provide string

type ProvideStatus string

const (
	PUBLIC  ProvideStatus = "public"
	PRIVATE ProvideStatus = "private"
)

type Factory func(param ...interface{}) interface{}

type DynamicProvider struct {
	Status ProvideStatus
	Name   Provide
	Value  interface{}
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

	existProvider := slices.IndexFunc(module.DataProviders, func(e *DynamicProvider) bool {
		return e.Name == providerName
	})
	if existProvider > -1 {
		return module.DataProviders[existProvider]
	}

	provider := &DynamicProvider{
		Value:  val,
		Name:   providerName,
		Status: PRIVATE,
	}

	module.DataProviders = append(module.DataProviders, provider)
	return provider
}

type FactoryOptions struct {
	Name    Provide
	Factory Factory
	Inject  []Provide
}

func (module *DynamicModule) NewFactoryProvider(opt FactoryOptions) *DynamicProvider {
	var values []interface{}
	for _, p := range opt.Inject {
		values = append(values, module.Ref(p))
	}

	provider := &DynamicProvider{
		Value:  opt.Factory(values...),
		Name:   opt.Name,
		Status: PRIVATE,
	}
	module.DataProviders = append(module.DataProviders, provider)
	return provider
}

func (module *DynamicModule) getExports() []*DynamicProvider {
	exports := make([]*DynamicProvider, 0)
	for _, v := range module.DataProviders {
		if v.Status == PUBLIC {
			exports = append(exports, v)
		}
	}

	return exports
}
