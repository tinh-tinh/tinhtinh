package core

import (
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/tinh-tinh/tinhtinh/utils"
)

type Router struct {
	// Name of controller own the route
	Name string
	// Method of route
	Method string
	// DEPRECATED: Use metadata in swagger package
	Tag string
	// Path detail of route
	Path string
	// Metadata of route
	Metadata []*Metadata
	// Handler function for router handler. Will receive Ctx instance
	Handler Handler
	// Middlewares of route
	Middlewares []Middleware
	// Pipes of route
	Dtos []Pipe
	// DEPRECATED: Use metadata in swagger package
	Security []string
	// Version of route
	Version string
	// Raw http handler
	httpHandler http.Handler
}

// getHandler returns a new http.Handler that combines the raw httpHandler
// of the route with all the middlewares of the route.
//
// If the raw httpHandler is not nil, it will be used. Otherwise, a new
// http.Handler will be created with the Handler of the route and the Metadata
// of the route.
//
// The middlewares of the route will be applied to the returned http.Handler.
// The order of the middlewares will be the order of the Middlewares field of
// the Router.
func (r *Router) getHandler(app *App) http.Handler {
	var mergeHandler http.Handler
	if r.httpHandler != nil {
		mergeHandler = r.httpHandler
	} else {
		mergeHandler = ParseCtx(app, r.Handler, r.Metadata...)
	}
	for _, v := range r.Middlewares {
		mid := ParseCtxMiddleware(app, v)
		mergeHandler = mid(mergeHandler)
	}

	return mergeHandler
}

// free clears the registered routes of the app and runs the garbage collector.
// It should be called after all the routes have been registered with the app.
func (app *App) free() {
	app.Module.Routers = nil
	runtime.GC()
}

// registerRoutes registers all the routes of the app with the Mux.
//
// It creates a map of routes with their paths as keys and their values as
// slices of Routers. It then iterates over the map and registers each route
// with the Mux. If a route has a version, it will be registered with a prefix
// of "v" followed by the version number. If the app has a prefix, it will be
// registered with the prefix.
func (app *App) registerRoutes() {
	routes := make(map[string][]*Router)

	for _, r := range app.Module.Routers {
		route := ParseRoute(r.Method + " " + r.Path)
		if app.version != nil && app.version.Type == URIVersion && r.Version != "" {
			route.SetPrefix("v" + r.Version)
		}
		route.SetPrefix(r.Name)
		if app.Prefix != "" {
			route.SetPrefix(app.Prefix)
		}
		utils.Log(
			utils.Green("[TT] "),
			utils.White(time.Now().Format("2006-01-02 15:04:05")),
			utils.Yellow(" [RoutesResolver] "),
			utils.Green(route.GetPath()+"\n"),
		)
		if routes[route.GetPath()] != nil {
			routes[route.GetPath()] = append(routes[route.GetPath()], r)
		} else {
			routes[route.GetPath()] = []*Router{r}
		}
	}

	for k, v := range routes {
		app.Mux.Handle(k, app.versionMiddleware(v))
		delete(routes, k)
	}

	app.free()
}

type Route struct {
	Method string
	Path   string
}

// ParseRoute takes a URL string and returns a Route object.
//
// The returned Route object has a Method and a Path. The Method is the first
// element of the split URL, and the Path is the concatenation of the remaining
// elements of the split URL. The Path is prefixed with a slash if it is not
// empty.
func ParseRoute(url string) Route {
	route := strings.Split(url, " ")

	var path string
	for i := 1; i < len(route); i++ {
		path += IfSlashPrefixString(route[i])
	}

	return Route{
		Method: route[0],
		Path:   path,
	}
}

// SetPrefix sets the prefix of the route.
//
// The prefix is prepended to the current path of the route. The prefix is
// prefixed with a slash if it is not empty.
func (r *Route) SetPrefix(prefix string) {
	r.Path = IfSlashPrefixString(prefix) + r.Path
}

// GetPath returns the path of the route.
//
// If the route's method is empty, the path is returned with a trailing slash.
// Otherwise, the path is returned with the method prepended to it.
func (r *Route) GetPath() string {
	if r.Method == "" {
		return r.Path + "/"
	}
	return r.Method + " " + r.Path
}

// IfSlashPrefixString takes a string and returns it with a slash prefix if it is not empty
// and does not already have a slash prefix. The returned string is also formatted to
// have lowercase letters and no spaces.
func IfSlashPrefixString(s string) string {
	if s == "" {
		return s
	}
	s = strings.TrimSuffix(s, "/")
	if strings.HasPrefix(s, "/") {
		return ToFormat(s)
	}
	return "/" + ToFormat(s)
}

// ToFormat takes a string and returns a formatted string. The string is
// converted to lowercase and spaces are removed.
func ToFormat(s string) string {
	result := strings.ToLower(s)
	return strings.ReplaceAll(result, " ", "")
}
