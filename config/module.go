package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func New[E any](path string) *Module {
	if path == "" {
		path = ".emv"
	}
	godotenv.Load(path)

	var env E
	return &Module{
		mapper: Scan(&env),
	}
}

type Module struct {
	mapper map[string]interface{}
}

func (m *Module) Get(key string) interface{} {
	fmt.Print(m.mapper)
	return m.mapper[key]
}
