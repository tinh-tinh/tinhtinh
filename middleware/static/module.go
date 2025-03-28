package static

import (
	"net/http"
	"strings"

	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func ForRoot(path string) core.Modules {
	path = strings.ReplaceAll(path, "/", "")
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("")

		routePath := "/" + path + "/"
		ctrl.Handler(path, http.StripPrefix(routePath, http.FileServer(http.Dir(path))))

		return ctrl
	}

	return func(module core.Module) core.Module {
		staticModule := module.New(core.NewModuleOptions{
			Controllers: []core.Controllers{controller},
		})

		return staticModule
	}
}
