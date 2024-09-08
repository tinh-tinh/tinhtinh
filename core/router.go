package core

import "net/http"

type Router struct {
	Tag      string
	Path     string
	Handler  http.Handler
	Dtos     []Pipe
	Security []string
}
