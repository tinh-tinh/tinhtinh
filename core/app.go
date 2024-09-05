package core

import (
	"log"
	"net/http"
	"sync"
)

type App struct {
	pool   sync.Pool
	Prefix string
	Mux    *http.ServeMux
	Module *DynamicModule
}

type ModuleParam func() *DynamicModule

func CreateFactory(module ModuleParam, prefix string) *App {
	app := &App{
		pool:   sync.Pool{},
		Module: module(),
		Prefix: prefix,
		Mux:    http.NewServeMux(),
	}

	for k, v := range app.Module.mux {
		route := ParseRoute(k)
		route.SetPrefix(app.Prefix)
		log.Printf("[RoutesResolvers] %s\n", route.GetPath())
		app.Mux.Handle(route.GetPath(), v)
	}
	return app
}

func (app *App) Listen(port int) {
	server := http.Server{
		Addr:    ":" + IntToString(port),
		Handler: app.Mux,
	}
	log.Printf("Server running on http://localhost:%d\n", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("error when running server %v", err)
	}
}
