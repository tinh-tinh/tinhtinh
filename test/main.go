package main

import (
	"fmt"

	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/middleware"
)

func UserModule(module *core.DynamicModule) *core.DynamicModule {
	userModule := module.New(core.NewModuleOptions{})

	userModule.OnInit(func(module *core.DynamicModule) {
		fmt.Println("ON Init")
	})
	return userModule
}

func AppService(module *core.DynamicModule) *core.DynamicProvider {
	provider := module.NewProvider("abc", "test")
	return provider
}

func AppController(module *core.DynamicModule) *core.DynamicController {
	ctrl := module.NewController("test")

	ctrl.Get("/", func(ctx core.Ctx) {
		name := ctrl.Inject("test")
		ctx.JSON(core.Map{
			"data": name,
		})
	})

	return ctrl
}

func AppModule() *core.DynamicModule {
	appModule := core.NewModule(core.NewModuleOptions{
		Imports:     []core.Module{UserModule},
		Controllers: []core.Controller{AppController},
		Providers:   []core.Provider{AppService},
	})

	return appModule
}

func main() {
	app := core.CreateFactory(AppModule, "api").EnableCors(middleware.CorsOptions{})
	app.BeforeShutdown(func() {
		fmt.Println("Before shutdown")
	})
	app.AfterShutdown(func() {
		fmt.Print("After shutdown")
	})
	app.Listen(3000)
}
