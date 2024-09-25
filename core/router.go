package core

import (
	"net/http"
	"time"

	"github.com/tinh-tinh/tinhtinh/utils"
)

type Router struct {
	Name     string
	Method   string
	Tag      string
	Path     string
	Handler  http.Handler
	Dtos     []Pipe
	Security []string
	Version  string
}

func (app *App) registerRoutes() {
	routes := make(map[string][]*Router)

	for _, r := range app.Module.Routers {
		route := ParseRoute(r.Method + " " + r.Path)
		if app.version != nil && app.version.Type == URIVersion && r.Version != "" {
			route.SetPrefix("v" + r.Version)
		}
		route.SetPrefix(app.Prefix + "/" + r.Name)
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
	}
}
