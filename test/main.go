package main

import (
	"fmt"

	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/middleware"
)

func UserProvider(module *core.DynamicModule) *core.DynamicProvider {
	provider := module.NewFactoryProvider(core.FactoryOptions{
		Name: "user",
		Factory: func(param ...interface{}) interface{} {
			return fmt.Sprintf("%vUser", param[0])
		},
		Inject: []core.Provide{"root"},
	})

	return provider
}

func UserController(module *core.DynamicModule) *core.DynamicController {
	ctrl := module.NewController("user")

	ctrl.Get("/", func(ctx core.Ctx) {
		name := ctrl.Inject("user")
		ctx.JSON(core.Map{
			"data": name,
		})
	})

	return ctrl
}

func UserModule(module *core.DynamicModule) *core.DynamicModule {
	userModule := module.New(core.NewModuleOptions{
		Controllers: []core.Controller{UserController},
		Providers:   []core.Provider{UserProvider},
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

func RootModule(module *core.DynamicModule) *core.DynamicModule {
	rootModule := module.New(core.NewModuleOptions{})

	rootModule.NewProvider("root", "root")
	rootModule.Export("root")

	return rootModule
}

func AppModule() *core.DynamicModule {
	appModule := core.NewModule(core.NewModuleOptions{
		Imports: []core.Module{
			RootModule,
			UserModule,
		},
		Controllers: []core.Controller{AppController},
		Providers:   []core.Provider{AppService},
	})

	return appModule
}

func main() {
	app := core.CreateFactory(AppModule, "api").EnableCors(middleware.CorsOptions{
		AllowedMethods: []string{"POST", "GET"},
	})
	app.BeforeShutdown(func() {
		fmt.Println("Before shutdown")
	})
	app.AfterShutdown(func() {
		fmt.Print("After shutdown")
	})
	app.Listen(3000)
}
