package core

import "testing"

func AppReqService(module *DynamicModule) *DynamicProvider {
	provider := module.NewReqProvider("test", func(ctx Ctx) interface{} {
		return "abc" + ctx.Headers("x-api-name")
	})
	return provider
}

func AppService(module *DynamicModule) *DynamicProvider {
	provider := module.NewProvider("test")
	return provider
}

func AppController(module *DynamicModule) *DynamicController {
	ctrl := module.NewController("test")

	ctrl.Get("/", func(ctx Ctx) {
		name := ctrl.InjectFactory("abc", ctx)
		ctx.JSON(Map{
			"data": name,
		})
	})

	return ctrl
}

func AppModule() *DynamicModule {
	appModule := NewModule(NewModuleOptions{
		Controllers: []Controller{AppController},
		Providers:   []Provider{AppService},
	})

	return appModule
}

func Test_App(t *testing.T) {
	app := CreateFactory(AppModule, "/api")
	app.Listen(3000)
}
