package core

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/tinh-tinh/tinhtinh/utils"
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
		utils.Log(
			utils.Green("[TT] "),
			utils.White(time.Now().Format("2006-01-02 15:04:05")),
			utils.Yellow(" [RoutesResolver] "),
			utils.Green(route.GetPath()+"\n"),
		)
		app.Module.mux = nil
		app.Mux.Handle(route.GetPath(), v)
	}
	return app
}

func (app *App) Listen(port int) {
	app.Module.MapperDoc = nil
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
