package core

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/middleware/cors"
)

func Test_EnableCors(t *testing.T) {
	appController := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "1",
			})
		})

		ctrl.Post("", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "2",
			})
		})

		ctrl.Patch("{id}", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "3",
			})
		})

		ctrl.Put("{id}", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "4",
			})
		})

		ctrl.Delete("{id}", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "5",
			})
		})
		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{appController},
		})

		return appModule
	}

	app := CreateFactory(module, "api")
	app.EnableCors(cors.CorsOptions{
		AllowedMethods: []string{"POST"},
	})

	testServer := httptest.NewServer(app.prepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}
