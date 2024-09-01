package core

type Provide string
type DynamicProvider struct {
	Name  Provide
	Value interface{}
}

func NewProvider(name Provide, value interface{}) *DynamicProvider {
	return &DynamicProvider{
		Name:  name,
		Value: value,
	}
}
