package router

import (
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

type RouteChildren struct {
	Path   string
	Module core.Modules
}

type Options struct {
	Path     string
	Module   core.Modules
	Children []*RouteChildren
}

func Register(opt ...Options) core.Modules {
	routers := RegistryRoutes(opt)
	return func(module core.Module) core.Module {
		imports := []core.Modules{}
		for _, router := range routers {
			imports = append(imports, func(module core.Module) core.Module {
				temp := module.New(core.NewModuleOptions{
					Imports: []core.Modules{router.Module},
				})
				for _, subRouter := range temp.GetRouters() {
					subRouter.Name = router.Path + core.IfSlashPrefixString(subRouter.Name)
				}

				return temp
			})
		}

		routerModule := module.New(core.NewModuleOptions{
			Imports: imports,
		})

		return routerModule
	}
}

func RegistryRoutes(options []Options) []*RouteChildren {
	routers := []*RouteChildren{}
	for _, option := range options {
		path := core.IfSlashPrefixString(option.Path)
		if option.Module != nil {
			routers = append(routers, &RouteChildren{
				Path:   path,
				Module: option.Module,
			})
		}
		if option.Children != nil {
			for _, child := range option.Children {
				childPath := path + core.IfSlashPrefixString(child.Path)
				routers = append(routers, &RouteChildren{
					Path:   childPath,
					Module: child.Module,
				})
			}
		}
	}

	return routers
}
