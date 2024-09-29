package core

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Tag(t *testing.T) {
	ctrl := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")
		ctrl.Tag("test").Get("/t1", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "ok",
			})
		})

		ctrl.Tag("test2").Get("/t2", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "ok",
			})
		})
		return ctrl
	}

	appModule := NewModule(NewModuleOptions{
		Controllers: []Controller{ctrl},
	})

	findT1 := slices.IndexFunc(appModule.Routers, func(r *Router) bool {
		return r.Path == "/t1"
	})
	require.NotEqual(t, -1, findT1)
	require.Equal(t, "test", appModule.Routers[findT1].Tag)

	findT2 := slices.IndexFunc(appModule.Routers, func(r *Router) bool {
		return r.Path == "/t2"
	})
	require.NotEqual(t, -1, findT2)
	require.Equal(t, "test2", appModule.Routers[findT2].Tag)
}

func Test_Version(t *testing.T) {
	ctrl := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")
		ctrl.Version("1").Get("/t1", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "ok",
			})
		})

		ctrl.Version("2").Get("/t2", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "ok",
			})
		})
		return ctrl
	}

	appModule := NewModule(NewModuleOptions{
		Controllers: []Controller{ctrl},
	})

	findT1 := slices.IndexFunc(appModule.Routers, func(r *Router) bool {
		return r.Path == "/t1"
	})
	require.NotEqual(t, -1, findT1)
	require.Equal(t, "1", appModule.Routers[findT1].Version)

	findT2 := slices.IndexFunc(appModule.Routers, func(r *Router) bool {
		return r.Path == "/t2"
	})
	require.NotEqual(t, -1, findT2)
	require.Equal(t, "2", appModule.Routers[findT2].Version)
}

func Test_AddSecurity(t *testing.T) {
	ctrl := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")
		ctrl.AddSecurity("auth").Get("/t1", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "ok",
			})
		})

		ctrl.Get("/t2", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "ok",
			})
		})
		return ctrl
	}

	appModule := NewModule(NewModuleOptions{
		Controllers: []Controller{ctrl},
	})

	findT1 := slices.IndexFunc(appModule.Routers, func(r *Router) bool {
		return r.Path == "/t1"
	})
	require.NotEqual(t, -1, findT1)
	require.Equal(t, "auth", appModule.Routers[findT1].Security[0])

	findT2 := slices.IndexFunc(appModule.Routers, func(r *Router) bool {
		return r.Path == "/t2"
	})
	require.NotEqual(t, -1, findT2)
	require.Equal(t, 0, len(appModule.Routers[findT2].Security))
}

func Test_free(t *testing.T) {
	module := NewModule(NewModuleOptions{})

	middleware := func(ctx Ctx) error {
		return ctx.Next()
	}

	controller := module.NewController("test").Use(middleware)

	require.NotEmpty(t, controller.middlewares)

	controller.Get("", func(ctx Ctx) {
		ctx.JSON(Map{
			"data": "ok",
		})
	})

	require.Empty(t, controller.middlewares)
}

func Test_Registry(t *testing.T) {
	module := NewModule(NewModuleOptions{})

	middleware := func(ctx Ctx) error {
		return ctx.Next()
	}

	controller := module.NewController("test").Use(middleware).Registry()

	controller.Get("", func(ctx Ctx) {
		ctx.JSON(Map{
			"data": 1,
		})
	})

	require.Len(t, controller.globalMiddlewares, 1)
	require.Len(t, controller.middlewares, 0)
}
