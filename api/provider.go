package api

type Provider struct {
	Name  string
	Value interface{}
}

type GlobalProvider struct {
	Name string
}

func NewProvider(name string, value interface{}) *Provider {
	return &Provider{
		Name:  name,
		Value: value,
	}
}
