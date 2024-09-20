package core

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/tinh-tinh/tinhtinh/middleware"
	"github.com/tinh-tinh/tinhtinh/utils"
)

type App struct {
	pool   sync.Pool
	Prefix string
	log    bool
	Mux    *http.ServeMux
	Module *DynamicModule
	cors   *middleware.Cors
	hooks  []*Hook
}

type ModuleParam func() *DynamicModule

func CreateFactory(module ModuleParam, prefix string) *App {
	app := &App{
		pool:   sync.Pool{},
		Module: module(),
		Prefix: prefix,
		Mux:    http.NewServeMux(),
	}
	utils.Log(
		utils.Green("[TT] "),
		utils.White(time.Now().Format("2006-01-02 15:04:05")),
		utils.Yellow(" [Module Initializer] "),
		utils.Green(utils.GetFunctionName(module)+"\n"),
	)

	app.Module.init()
	for _, r := range app.Module.Routers {
		route := ParseRoute(r.Path)
		route.SetPrefix(app.Prefix)
		utils.Log(
			utils.Green("[TT] "),
			utils.White(time.Now().Format("2006-01-02 15:04:05")),
			utils.Yellow(" [RoutesResolver] "),
			utils.Green(route.GetPath()+"\n"),
		)
		app.Mux.Handle(route.GetPath(), r.Handler)
	}

	app.Mux.Handle(IfSlashPrefixString(app.Prefix), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, "API is running")
		if err != nil {
			log.Fatalf("error when running server %v", err)
		}
	}))
	return app
}

func (app *App) Log() *App {
	app.log = true
	return app
}

func (app *App) EnableCors(opt middleware.CorsOptions) *App {
	app.cors = middleware.NewCors(opt)
	return app
}

func (app *App) Listen(port int) {
	server := http.Server{
		Addr:    ":" + IntToString(port),
		Handler: app.Mux,
	}
	if app.log {
		loggedRouter := logRequests(server.Handler)
		server.Handler = loggedRouter
	}
	if app.cors != nil {
		corsHandler := app.cors.Handler(server.Handler)
		server.Handler = corsHandler
	}

	log.Printf("Server running on http://localhost:%d/%s\n", port, app.Prefix)

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error when running server %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	for _, hook := range app.hooks {
		if hook.RunAt == BEFORE_SHUTDOWN {
			hook.fnc()
		}
	}

	err := server.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatalf("error when shutdown server %v", err)
	}
	log.Println("Server shutdown")
	for _, hook := range app.hooks {
		if hook.RunAt == AFTER_SHUTDOWN {
			hook.fnc()
		}
	}
}
