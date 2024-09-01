package core

import (
	"log"
	"net/http"
	"sync"
)

type App struct {
	pool   sync.Pool
	prefix string
	module *DynamicModule
}

type ModuleParam func() *DynamicModule

func CreateFactory(module ModuleParam) *App {
	app := &App{
		pool:   sync.Pool{},
		module: module(),
	}

	return app
}

func (app *App) SetGlobalPrefix(prefix string) {
	app.prefix = prefix
}

func (app *App) Listen(port int) {
	mux := http.NewServeMux()

	for k, v := range app.module.mux {
		route := ParseRoute(k)
		route.SetPrefix(app.prefix)

		log.Printf("[RoutesResolvers] %s\n", route.GetPath())
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
