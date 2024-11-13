package core_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/core"
)

func Test_Version(t *testing.T) {
	ctrl := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")
		ctrl.Version("1").Get("/t1", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "ok",
			})
		})

		ctrl.Version("2").Get("/t2", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "ok",
			})
		})
		return ctrl
	}

	appModule := core.NewModule(core.NewModuleOptions{
		Controllers: []core.Controller{ctrl},
	})

	findT1 := slices.IndexFunc(appModule.Routers, func(r *core.Router) bool {
		return r.Path == "/t1"
	})
	require.NotEqual(t, -1, findT1)
	require.Equal(t, "1", appModule.Routers[findT1].Version)

	findT2 := slices.IndexFunc(appModule.Routers, func(r *core.Router) bool {
		return r.Path == "/t2"
	})
	require.NotEqual(t, -1, findT2)
	require.Equal(t, "2", appModule.Routers[findT2].Version)
}

func Test_free(t *testing.T) {
	module := core.NewModule(core.NewModuleOptions{})

	middleware := func(ctx core.Ctx) error {
		return ctx.Next()
	}

	controller := module.NewController("test").Use(middleware)

	controller.Get("", func(ctx core.Ctx) error {
		return ctx.JSON(core.Map{
			"data": "ok",
		})
	})
}

func Test_Registry(t *testing.T) {
	module := core.NewModule(core.NewModuleOptions{})

	middleware := func(ctx core.Ctx) error {
		return ctx.Next()
	}

	controller := module.NewController("test").Use(middleware).Registry()

	controller.Get("", func(ctx core.Ctx) error {
		return ctx.JSON(core.Map{
			"data": 1,
		})
	})
}
