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

	"github.com/tinh-tinh/tinhtinh/middleware/cors"
	"github.com/tinh-tinh/tinhtinh/utils"
)

type App struct {
	pool       sync.Pool
	Prefix     string
	Mux        *http.ServeMux
	Module     *DynamicModule
	cors       *cors.Cors
	hooks      []*Hook
	Middleware []Middleware
}

type ModuleParam func() *DynamicModule

// CreateFactory is a function that creates an App instance with a DynamicModule
// and a specified prefix. The DynamicModule is created by calling the given
// module function, and the prefix is set on the App instance. The App instance's
// Mux is set to a new http.ServeMux, and the Module is initialized. The routes
// on the Module are resolved and added to the Mux. The Mux is then set to
// handle the root path with a handler that writes "API is running" to the
// response writer. Finally, the App instance is returned.
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

// EnableCors enables CORS on the API server. The passed in options are used
// to configure the CORS middleware.
func (app *App) EnableCors(opt cors.CorsOptions) *App {
	app.cors = cors.NewCors(opt)
	return app
}

// Use appends the given middleware functions to the App instance's list of
// middleware handlers. The middleware handlers are run in the order they are
// added to the App instance. The middleware handlers are run before the
// App instance's handlers. The App instance's middleware handlers are run
// after the module's middleware handlers. The module middleware handlers are
// run after the module's parent middleware handlers. The module middleware
// handlers are run before the module's controllers. The App instance's
// middleware handlers are run before the App instance's handlers.
func (app *App) Use(middleware ...Middleware) *App {
	app.Middleware = append(app.Middleware, middleware...)
	return app
}

// Listen starts the App instance's HTTP server on the given port. The App
// instance's Mux is used as the handler for the server. If the App instance's
// CORS middleware is not nil, it is added to the server's handler. If the App
// instance's middleware handlers are not empty, they are added to the server's
// handler. The server is then started in a goroutine. The App instance's
// shutdown hooks are run in the order they were added. The server is then shut
// down gracefully. The App instance's shutdown hooks are run in the order they
// were added.
func (app *App) Listen(port int) {
	server := http.Server{
		Addr:    ":" + IntToString(port),
		Handler: app.Mux,
	}
	if app.cors != nil {
		corsHandler := app.cors.Handler(server.Handler)
		server.Handler = corsHandler
	}

	if len(app.Middleware) > 0 {
		for _, m := range app.Middleware {
			server.Handler = m(server.Handler)
		}
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
