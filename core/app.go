package core

import (
	"context"
	"encoding/json"
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

	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/common/color"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/cors"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/session"
)

type App struct {
	pool sync.Pool
	// prefix is the URL prefix of the API.
	Prefix string
	// Mux is the http.ServeMux that the App uses to serve requests.
	// The App uses this Mux to serve requests.
	Mux *http.ServeMux
	// Module is the module that the App uses to initialize itself.
	// The App uses this Module to initialize itself.
	Module Module
	// cors is the CORS middleware.
	cors *cors.Cors
	// version is the type version of the API.
	version *Version
	// hooks are the hooks that the App uses to initialize itself.
	// Two hooks can registered is: BeforeShutdown and AfterShutdown
	hooks []*Hook
	// middleware are the middleware that the App uses to initialize itself.
	Middleware []middlewareRaw
	// encoder is the encoder that the App uses to initialize itself.
	encoder Encode
	// decoder is the decoder that the App uses to initialize itself.
	decoder Decode
	// session is the session that the App uses to initialize itself.
	session *session.Config
	// errorHandler is the error handler that the App uses to initialize itself.
	errorHandler ErrorHandler
	timeout      time.Duration
	Services     []Service
}

type ModuleParam func() Module
type AppOptions struct {
	// Encoder is the encoder that the App uses to initialize itself.
	Encoder Encode
	// Decoder is the decoder that the App uses to initialize itself.
	Decoder Decode
	// Session is the session that the App uses to initialize itself.
	Session *session.Config
	// ErrorHandler is the error handler that the App uses to initialize itself.
	ErrorHandler ErrorHandler
	Timeout      time.Duration
}

type Service interface {
	Create(module Module)
	Listen()
}

// CreateFactory creates a new App with the given module and options.
//
// The module is the module that the App uses to initialize itself.
// The options are the options that the App uses to initialize itself.
// The options can be used to override the default encoder, decoder, session, and
// error handler.
//
// The App is created with a default encoder and decoder of json.Marshal and
// json.Unmarshal respectively.
// The App is created with a default session of nil.
// The App is created with a default error handler of ErrorHandlerDefault.
//
// The App is initialized by calling the init method of the module.
// The App is then returned.
func CreateFactory(module ModuleParam, opt ...AppOptions) *App {
	app := &App{
		Module:       module(),
		Mux:          http.NewServeMux(),
		encoder:      json.Marshal,
		decoder:      json.Unmarshal,
		errorHandler: ErrorHandlerDefault,
	}

	app.pool = sync.Pool{
		New: func() any {
			return NewCtx(app)
		},
	}

	if len(opt) > 0 {
		if opt[0].Encoder != nil {
			app.encoder = opt[0].Encoder
		}
		if opt[0].Decoder != nil {
			app.decoder = opt[0].Decoder
		}
		if opt[0].Session != nil {
			app.session = opt[0].Session
		}
		if opt[0].ErrorHandler != nil {
			app.errorHandler = opt[0].ErrorHandler
		}
		if opt[0].Timeout != 0 {
			app.timeout = opt[0].Timeout
		}
	}

	fmt.Printf("%s %s %s %s\n",
		color.Green("[TT]"),
		color.White(time.Now().Format("2006-01-02 15:04:05")),
		color.Yellow("[Module Initializer]"),
		color.Green(common.GetFunctionName(module)),
	)
	app.Module.init()
	return app
}

// SetGlobalPrefix sets the global prefix of the API. The global prefix is
// prepended to the URL paths of the API.
func (app *App) SetGlobalPrefix(prefix string) *App {
	app.Prefix = IfSlashPrefixString(prefix)
	return app
}

// EnableCors enables CORS on the API server. The passed in options are used
// to configure the CORS middleware.
func (app *App) EnableCors(opt cors.Options) *App {
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
func (app *App) Use(middleware ...middlewareRaw) *App {
	app.Middleware = append(app.Middleware, middleware...)
	return app
}

// PrepareBeforeListen is a helper function that prepares the App instance's
// HTTP handler before listening. It registers the routes from the App
// instance's Module, and adds a handler that writes "API is running" to the
// request writer. It also adds the App instance's CORS middleware if it is
// not nil. Finally, it chains the App instance's middleware handlers and
// returns the final handler.
func (app *App) PrepareBeforeListen() http.Handler {
	app.registerRoutes()
	prefix := "/"
	if app.Prefix != "" {
		prefix = IfSlashPrefixString(app.Prefix)
	}
	app.Mux.Handle(prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	if app.timeout != 0 {
		handler = http.TimeoutHandler(handler, app.timeout, "timeout")
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
	handler := app.PrepareBeforeListen()
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

func (app *App) ConnectMicroservice(svc Service) {
	svc.Create(app.Module)
	app.Services = append(app.Services, svc)
}

func (app *App) StartAllMicroservices() {
	for _, svc := range app.Services {
		go svc.Listen()
	}
}
