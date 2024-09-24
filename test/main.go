package main

import (
	"fmt"
	"net/http"

	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/middleware/cors"
	"github.com/tinh-tinh/tinhtinh/middleware/logger"
)

func UserProvider(module *core.DynamicModule) *core.DynamicProvider {
	provider := module.NewProvider(core.ProviderOptions{
		Name: "user",
		Factory: func(param ...interface{}) interface{} {
			fmt.Println("param", param[1])
			return fmt.Sprintf("Root%vUser", param[1])
		},
		Inject: []core.Provide{core.REQUEST, "root", "jajaj"},
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
		Scope:       core.Request,
		Controllers: []core.Controller{UserController},
		Providers:   []core.Provider{UserProvider},
	})

	return userModule
}

func AppService(module *core.DynamicModule) *core.DynamicProvider {
	provider := module.NewProvider(core.ProviderOptions{
		Name:  "app",
		Value: "test",
	})
	return provider
}

func AppController(module *core.DynamicModule) *core.DynamicController {
	ctrl := module.NewController("test")

	ctrl.Get("/", func(ctx core.Ctx) {
		name := ctrl.Inject("sjkfkjgjkebjgrkebjkgb")
		ctx.JSON(core.Map{
			"data": name,
		})
	})

	return ctrl
}

func RootProvider(module *core.DynamicModule) *core.DynamicProvider {
	provider := module.NewProvider(core.ProviderOptions{
		Name: "root",
		Factory: func(param ...interface{}) interface{} {
			req := param[0].(*http.Request)
			return fmt.Sprintf("%vRoot", req.Header.Get("x-api-key"))
		},
		Inject: []core.Provide{core.REQUEST},
	})

	return provider
}

func RootModule(module *core.DynamicModule) *core.DynamicModule {
	rootModule := module.New(core.NewModuleOptions{
		Scope:     core.Request,
		Providers: []core.Provider{RootProvider},
		Exports:   []core.Provider{RootProvider},
	})
	return rootModule
}

func AbcModule(module *core.DynamicModule) *core.DynamicModule {
	abcModule := module.New(core.NewModuleOptions{
		Scope: core.Request,
	})

	abcModule.NewProvider(core.ProviderOptions{
		Name: "jajaj",
		Factory: func(param ...interface{}) interface{} {
			return fmt.Sprintf("%vAbc", param[0])
		},
		Inject: []core.Provide{core.REQUEST},
	})
	abcModule.Export("jajaj")

	return abcModule
}

func AppModule() *core.DynamicModule {
	appModule := core.NewModule(core.NewModuleOptions{
		Imports: []core.Module{
			RootModule,
			AbcModule,
			UserModule,
		},
		// Controllers: []core.Controller{AppController},
		// Providers:   []core.Provider{AppService},
	})

	return appModule
}

func main() {
	app := core.CreateFactory(AppModule, "api").EnableCors(cors.CorsOptions{
		AllowedMethods: []string{"POST", "GET"},
		AllowedHeaders: []string{"*"},
	})

	app.Use(logger.Middleware(logger.MiddlewareOptions{
		Rotate: true,
		Format: "${method} ${path} ${status} ${latency}",
	}))
	app.BeforeShutdown(func() {
		fmt.Println("Before shutdown")
	})
	app.AfterShutdown(func() {
		fmt.Print("After shutdown")
	})
	app.Listen(3000)
}
