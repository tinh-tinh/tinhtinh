package core

import (
	"context"
	"errors"
	"fmt"
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
	version    *Version
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

// prepareBeforeListen is a helper function that prepares the App instance's
// HTTP handler before listening. It registers the routes from the App
// instance's Module, and adds a handler that writes "API is running" to the
// request writer. It also adds the App instance's CORS middleware if it is
// not nil. Finally, it chains the App instance's middleware handlers and
// returns the final handler.
func (app *App) prepareBeforeListen() http.Handler {
	app.registerRoutes()
	app.Mux.Handle(IfSlashPrefixString(app.Prefix), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, "API is running")
		if err != nil {
			log.Fatalf("error when running server %v", err)
		}
	}))

	var handler http.Handler
	handler = app.Mux
	if app.cors != nil {
		corsHandler := app.cors.Handler(handler)
		handler = corsHandler
	}

	if len(app.Middleware) > 0 {
		for _, m := range app.Middleware {
			handler = m(handler)
		}
	}

	return handler
}

// Listen starts the API server on the specified port.
//
// It first prepares the server's handler with the module's routes and
// middleware. Then it starts the server and logs a message to the console
// indicating that the server is running.
//
// The server is then shut down when the process receives a SIGINT or SIGTERM
// signal. It waits for 10 seconds for the server to shut down, and if it does
// not shut down within that time, it prints an error message to the console.
//
// Finally, it runs any hooks registered with the AFTER_SHUTDOWN run-at value.
func (app *App) Listen(port int) {
	handler := app.prepareBeforeListen()
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
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
