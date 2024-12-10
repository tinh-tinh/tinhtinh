package compression_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/compression"
)

func Test_Compress(t *testing.T) {
	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "ok",
			})
		})

		return ctrl
	}

	appModule := func() core.Module {
		return core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{appController},
		})
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")

	app.Use(compression.Handler())

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	req, err := http.NewRequest("GET", testServer.URL+"/api/test", nil)
	require.Nil(t, err)
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "gzip", resp.Header.Get("Content-Encoding"))
}
