package static

import (
	"net/http"
	"strings"

	"github.com/tinh-tinh/tinhtinh/core"
)

func ForRoot(path string) core.Module {
	path = strings.ReplaceAll(path, "/", "")
	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("")

		routePath := "/" + path + "/"
		ctrl.Handler(path, http.StripPrefix(routePath, http.FileServer(http.Dir(path))))

		return ctrl
	}

	return func(module *core.DynamicModule) *core.DynamicModule {
		staticModule := module.New(core.NewModuleOptions{
			Controllers: []core.Controller{controller},
		})

		return staticModule
	}
}
