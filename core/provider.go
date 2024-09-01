package core

type Provide string
type DynamicProvider struct {
	module *DynamicModule
}

func NewProvider(module *DynamicModule) *DynamicProvider {
	return &DynamicProvider{
		module: module,
	}
}

func (p *DynamicProvider) Get(key Provide) interface{} {
	return p.module.mapperValue[key]
}

func (p *DynamicProvider) Set(key Provide, value interface{}) {
	p.module.mapperValue[key] = value
}
