package router

import (
	"github.com/tinh-tinh/tinhtinh/core"
)

type RouteChildren struct {
	Path   string
	Module core.Module
}

type Options struct {
	Path     string
	Module   core.Module
	Children []*RouteChildren
}

func Register(opt ...Options) core.Module {
	routers := RegistryRoutes(opt)
	return func(module *core.DynamicModule) *core.DynamicModule {
		imports := []core.Module{}
		for _, router := range routers {
			imports = append(imports, func(module *core.DynamicModule) *core.DynamicModule {
				temp := module.New(core.NewModuleOptions{
					Imports: []core.Module{router.Module},
				})
				for _, subRouter := range temp.Routers {
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
