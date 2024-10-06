package core

import (
	"net/http"
	"runtime"
	"strings"
)

type Handler func(ctx Ctx)

type DynamicController struct {
	name              string
	tag               string
	version           string
	metadata          []*Metadata
	Dtos              []Pipe
	Security          []string
	middlewares       []Middleware
	globalMiddlewares []Middleware
	module            *DynamicModule
}

// NewController creates a new controller with the given name.
// The controller will be initialized with the default middlewares
// and a blank list of dtos and security roles. The module
// parameter is used to store the controller in the module's
// registry.
func (module *DynamicModule) NewController(name string) *DynamicController {
	return &DynamicController{
		name:              strings.ToLower(name),
		tag:               name,
		globalMiddlewares: module.Middlewares,
		Dtos:              []Pipe{},
		Security:          []string{},
		module:            module,
		version:           "",
	}
}

// Use appends the given middleware functions to the module's list of
// middleware handlers. The middleware handlers are run in the order they
// are added to the module. The middleware handlers are run before the
// controller's middleware handlers. The module's middleware handlers
// are run after the module's parent middleware handlers. The module
// middleware handlers are run before the module's controllers.
func (module *DynamicModule) Use(middleware ...Middleware) *DynamicModule {
	module.Middlewares = append(module.Middlewares, middleware...)

	return module
}

// Tag sets the tag for the controller. The tag is used to generate the
// route path for the controller. The tag is used in combination with the
// module's prefix and the controller's name to generate the route path.
// The tag is required for the controller to be registered with the
// module. The tag is also used as the default value for the
// controller's tag if the tag is not set.
func (c *DynamicController) Tag(tag string) *DynamicController {
	c.tag = tag
	return c
}

// Version sets the version for the controller. The version is used to
// generate the route path for the controller. The version is used in
// combination with the module's prefix and the controller's name and
// tag to generate the route path. The version is required for the
// controller to be registered with the module. The version is also
// used as the default value for the controller's version if the
// version is not set.
func (c *DynamicController) Version(version string) *DynamicController {
	c.version = version
	return c
}

// Use appends the given middleware functions to the controller's list of
// middleware handlers. The middleware handlers are run in the order they
// are added to the controller. The middleware handlers are run before the
// controller's handlers. The controller's middleware handlers are run
// after the module's middleware handlers. The module middleware handlers
// are run after the module's parent middleware handlers. The module
// middleware handlers are run before the module's controllers. The
// controller middleware handlers are run before the controller's
// handlers.
func (c *DynamicController) Use(middleware ...Middleware) *DynamicController {
	c.middlewares = append(c.middlewares, middleware...)
	return c
}

// Guard registers the given guard functions with the controller. The guard
// functions are run in the order they are added to the controller. The guard
// functions are run before the controller's middleware handlers. The guard
// functions are run after the module's middleware handlers. The guard functions
// are run before the controller's handlers. If any of the guard functions
// return false, the request will be rejected with a 403 status code.
func (c *DynamicController) Guard(guards ...Guard) *DynamicController {
	for _, v := range guards {
		mid := c.ParseGuard(v)
		c.middlewares = append(c.middlewares, mid)
	}
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
func (c *DynamicController) Pipe(dtos ...Pipe) *DynamicController {
	c.Dtos = append(c.Dtos, dtos...)
	middleware := PipeMiddleware(dtos...)
	c.middlewares = append(c.middlewares, middleware)
	return c
}

// AddSecurity adds the given security roles to the controller. The security
// roles are checked in the order they are added to the controller. If any of
// the security roles return false, the request will be rejected with a 403
// status code.
func (c *DynamicController) AddSecurity(security ...string) *DynamicController {
	c.Security = append(c.Security, security...)
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
func (c *DynamicController) Registry() *DynamicController {
	c.globalMiddlewares = append(c.globalMiddlewares, c.middlewares...)
	c.middlewares = []Middleware{}

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

func (c *DynamicController) Handler(path string, handler http.Handler) {
	router := &Router{
		Name:        c.name,
		Tag:         c.tag,
		Path:        path,
		Middlewares: append(c.globalMiddlewares, c.middlewares...),
		Dtos:        c.Dtos,
		Security:    c.Security,
		Version:     c.version,
		httpHandler: handler,
	}
	c.module.Routers = append(c.module.Routers, router)
	c.free()
}

func (c *DynamicController) registry(method string, path string, handler Handler) {
	router := &Router{
		Name:        c.name,
		Method:      method,
		Tag:         c.tag,
		Path:        path,
		Middlewares: append(c.globalMiddlewares, c.middlewares...),
		Handler:     handler,
		Dtos:        c.Dtos,
		Security:    c.Security,
		Version:     c.version,
	}
	c.module.Routers = append(c.module.Routers, router)
	c.free()
}

func (c *DynamicController) free() {
	c.middlewares = []Middleware{}
	c.Dtos = nil
	c.Security = nil
	runtime.GC()
}

func (c *DynamicController) Inject(name Provide) interface{} {
	return c.module.Ref(name)
}
