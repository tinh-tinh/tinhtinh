package core

import "errors"

type Provide string
type DynamicProvider struct {
	module *DynamicModule
}

func (module *DynamicModule) NewProvider() *DynamicProvider {
	return &DynamicProvider{
		module: module,
	}
}

func (p *DynamicProvider) Get(key Provide) interface{} {
	return p.module.providers[key]
}

func (p *DynamicProvider) Set(key Provide, value interface{}) {
	p.module.providers[key] = value
}

func (p *DynamicProvider) Export(key Provide) {
	val := p.module.providers[key]
	if val == nil {
		panic(errors.New("invalid provider"))
	}
	p.module.Exports[key] = val
}
