package core

type DynamicProvider struct {
	Name  string
	Value interface{}
}

func NewProvider(name string, value interface{}) *DynamicProvider {
	return &DynamicProvider{
		Name:  name,
		Value: value,
	}
}
