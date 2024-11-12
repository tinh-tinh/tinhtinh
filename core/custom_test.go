package core

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CustomCtx(t *testing.T) {
	tenant := CreateWrapper(func(data interface{}, ctx Ctx) string {
		isMaster, ok := data.(bool)
		if ok && isMaster {
			return "master"
		}
		return ctx.Req().Header.Get("x-tenant-id")
	})
	appController := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", tenant.Handler(true, func(w WrappedCtx[string]) error {
			return w.Ctx.JSON(Map{
				"data": w.Data,
			})
		}))

		ctrl.Get("tenant", tenant.Handler(false, func(w WrappedCtx[string]) error {
			return w.Ctx.JSON(Map{
				"data": w.Data,
			})
		}))

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{appController},
		})

		return appModule
	}

	app := CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	req, err := http.NewRequest("GET", testServer.URL+"/api/test", nil)
	require.Nil(t, err)

	req.Header.Set("x-tenant-id", "babadook")
	resp, err := testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "master", res.Data)

	req, err = http.NewRequest("GET", testServer.URL+"/api/test/tenant", nil)
	require.Nil(t, err)

	req.Header.Set("x-tenant-id", "babadook")
	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res2 Response
	err = json.Unmarshal(data, &res2)
	require.Nil(t, err)
	require.Equal(t, "babadook", res2.Data)
}
