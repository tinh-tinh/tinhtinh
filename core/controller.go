package core

import (
	"net/http"
	"runtime"
	"strings"
)

type Handler func(ctx Ctx) error

type Controller interface {
	Version(version string) Controller
	Metadata(metadata ...*Metadata) Controller
	Pipe(dtos ...PipeDto) Controller
	Guard(guards ...Guard) Controller
	Interceptor(interceptor Interceptor) Controller
	Use(middleware ...Middleware) Controller
	Composition(ctrl Controller) Controller
	Registry() Controller
	Get(path string, handler Handler)
	Post(path string, handler Handler)
	Put(path string, handler Handler)
	Patch(path string, handler Handler)
	Delete(path string, handler Handler)
	Handler(path string, handler http.Handler)
	Ref(name Provide, ctx ...Ctx) interface{}
	getMiddlewares() []Middleware
	getMetadata() []*Metadata
	GetDtos() []PipeDto
	Sse(path string, sseFnc SseFnc)
}

type DynamicController struct {
	// name of the controller. This is prefix when registry routes in controller
	name string
	// Mark version
	version string
	// Use for apply metadata for all routes in controller
	globalMetadata []*Metadata
	// Use for apply metadata for each route
	metadata []*Metadata
	// Data validator of each routes
	Dtos []PipeDto
	// Use for apply middlewares for each route
	middlewares []Middleware
	// Use for apply middlewares for all routes
	globalMiddlewares []Middleware
	// Parent module for this controller
	module            *DynamicModule
	globalInterceptor Interceptor
	interceptor       Interceptor
}

// NewController creates a new controller with the given name.
// The controller will be initialized with the default middlewares
// and a blank list of dtos and security roles. The module
// parameter is used to store the controller in the module's
// registry.
func (module *DynamicModule) NewController(name string) Controller {
	return &DynamicController{
		name:              strings.ToLower(name),
		globalMiddlewares: module.Middlewares,
		globalInterceptor: module.interceptor,
		Dtos:              []PipeDto{},
		module:            module,
		version:           "",
	}
}

// Version sets the version for the controller. The version is used to
// generate the route path for the controller. The version is used in
// combination with the module's prefix and the controller's name and
// tag to generate the route path. The version is required for the
// controller to be registered with the module. The version is also
// used as the default value for the controller's version if the
// version is not set.
func (c *DynamicController) Version(version string) Controller {
	c.version = version
	return c
}

// Pipe registers the given dtos with the controller. The dtos are run in the
// order they are added to the controller. The dtos are run before the
// controller's handlers. The dtos are run after the module's middleware
// handlers. The module middleware handlers are run after the module's
// parent middleware handlers. The module middleware handlers are run
// before the module's controllers. The controller's dtos are run before
// the controller's handlers. If any of the dtos return an error, the
// request will be rejected with the error.
func (c *DynamicController) Pipe(dtos ...PipeDto) Controller {
	c.Dtos = append(c.Dtos, dtos...)
	middleware := PipeMiddleware(dtos...)
	c.middlewares = append(c.middlewares, middleware)
	return c
}

// Registry saves the current middleware handlers to the controller's global
// middleware handlers and clears the current middleware handlers. This is
// useful when the controller is used as a sub-controller in another
// controller. The global middleware handlers of the sub-controller are
// appended to the middleware handlers of the parent controller, and the
// global middleware handlers of the sub-controller are cleared. This
// ensures that the middleware handlers of the sub-controller are not
// executed twice.
func (c *DynamicController) Registry() Controller {
	c.globalMiddlewares = append(c.globalMiddlewares, c.middlewares...)
	c.middlewares = []Middleware{}
	c.globalMetadata = append(c.globalMetadata, c.metadata...)
	c.globalInterceptor = c.interceptor
	c.metadata = []*Metadata{}

	return c
}

// Get registers a new GET route with the given path and handler.
func (c *DynamicController) Get(path string, handler Handler) {
	c.registry("GET", path, handler)
}

// Post registers a new POST route with the given path and handler.
func (c *DynamicController) Post(path string, handler Handler) {
	c.registry("POST", path, handler)
}

// Patch registers a new PATCH route with the given path and handler.
func (c *DynamicController) Patch(path string, handler Handler) {
	c.registry("PATCH", path, handler)
}

// Put registers a new PUT route with the given path and handler.
func (c *DynamicController) Put(path string, handler Handler) {
	c.registry("PUT", path, handler)
}

// Delete registers a new DELETE route with the given path and handler.
func (c *DynamicController) Delete(path string, handler Handler) {
	c.registry("DELETE", path, handler)
}

// Handler registers a new route with the given path and handler.
//
// The route's middlewares are the combination of the controller's global
// middlewares and the controller's current middlewares. The route's metadata
// are the combination of the controller's global metadata and the
// controller's current metadata. The route's dtos are the controller's dtos.
// The route's version is the controller's version.
//
// The route is registered with the controller's module and the controller's
// middlewares and metadata are cleared.
func (c *DynamicController) Handler(path string, handler http.Handler) {
	router := &Router{
		Name:        c.name,
		Path:        path,
		Middlewares: append(c.globalMiddlewares, c.middlewares...),
		Metadata:    append(c.globalMetadata, c.metadata...),
		Dtos:        c.Dtos,
		Version:     c.version,
		interceptor: c.interceptor,
		httpHandler: handler,
	}
	c.module.Routers = append(c.module.Routers, router)
	c.free()
}

// registry registers a new route with the given method, path and handler.
//
// The route's middlewares are the combination of the controller's global
// middlewares and the controller's current middlewares. The route's metadata
// are the combination of the controller's global metadata and the
// controller's current metadata. The route's dtos are the controller's dtos.
// The route's version is the controller's version.
//
// The route is registered with the controller's module and the controller's
// middlewares and metadata are cleared.
func (c *DynamicController) registry(method string, path string, handler Handler) {
	if path == "/" {
		path = ""
	}
	router := &Router{
		Name:        c.name,
		Method:      method,
		Path:        path,
		Middlewares: append(c.globalMiddlewares, c.middlewares...),
		Metadata:    append(c.globalMetadata, c.metadata...),
		Handler:     handler,
		Dtos:        c.Dtos,
		Version:     c.version,
	}
	if c.interceptor != nil {
		router.interceptor = c.interceptor
	} else {
		router.interceptor = c.globalInterceptor
	}
	c.module.Routers = append(c.module.Routers, router)
	c.free()
}

// free clears the controller's middlewares, dtos, security and metadata.
// It is called after a route is registered with the controller's module.
// It is useful when the controller is used as a sub-controller in another
// controller. The global middlewares of the sub-controller are appended to
// the middleware handlers of the parent controller, and the global middlewares
// of the sub-controller are cleared. This ensures that the middleware handlers
// of the sub-controller are not executed twice.
func (c *DynamicController) free() {
	c.middlewares = []Middleware{}
	c.Dtos = nil
	c.interceptor = nil
	c.metadata = []*Metadata{}
	runtime.GC()
}

func (c *DynamicController) Ref(name Provide, ctx ...Ctx) interface{} {
	return c.module.Ref(name, ctx...)
}

func (c *DynamicController) getMiddlewares() []Middleware {
	return c.middlewares
}

func (c *DynamicController) getGlobalMiddlewares() []Middleware {
	return c.globalMiddlewares
}

func (c *DynamicController) getMetadata() []*Metadata {
	return c.metadata
}

func (c *DynamicController) GetDtos() []PipeDto {
	return c.Dtos
}
