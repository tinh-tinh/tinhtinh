package static

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/core"
)

func Test_Static(t *testing.T) {
	appController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "ok",
			})
		})

		return ctrl
	}

	appModule := func() *core.DynamicModule {
		return core.NewModule(core.NewModuleOptions{
			Imports:     []core.Module{ForRoot("upload")},
			Controllers: []core.Controller{appController},
		})
	}

	app := core.CreateFactory(appModule)

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/upload/")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
