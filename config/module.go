package config

import (
	"sync"

	"github.com/joho/godotenv"
)

type Module struct {
	sync.Pool
}

var pool sync.Pool

func Register[E any](path string) {
	if path == "" {
		path = ".env"
	}
	err := godotenv.Load(path)
	if err != nil {
		panic(err)
	}

	pool = sync.Pool{
		New: func() any {
			var env E
			Scan(&env)
			return env
		},
	}
}

func Get[E any]() E {
	return pool.Get().(E)
}
