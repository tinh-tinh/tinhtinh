package api

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type App struct {
	pool   sync.Pool
	prefix string
	routes map[string]http.Handler
}

func New(module *Module) *App {
	var routes = make(map[string]http.Handler)

	for k, v := range module.mux {
		routes[k] = v
	}
	module = nil
	return &App{
		pool:   sync.Pool{},
		routes: routes,
	}
}

func (app *App) SetGlobalPrefix(prefix string) {
	app.prefix = prefix
}

func (app *App) Listen(port int) {
	mux := http.NewServeMux()

	for k, v := range app.routes {
		route := ParseRoute(k)
		route.SetPrefix(app.prefix)

		fmt.Printf("[RoutesResolvers] %s\n", route.GetPath())
		mux.Handle(route.GetPath(), v)
	}

	server := http.Server{
		Addr:    ":" + IntToString(port),
		Handler: mux,
	}
	log.Printf("Server running on http://localhost:%d\n", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("error when running server %v", err)
	}
}
