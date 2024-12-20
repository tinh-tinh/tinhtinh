package core_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/core"
)

func Test_Compose(t *testing.T) {
	// Guard
	guard := func(ctrl core.RefProvider, ctx *core.Ctx) bool {
		return ctx.Query("key") == "value"
	}

	// Pipe
	type FilterDto struct {
		Name string `validate:"required" query:"name"`
	}

	// Metadata
	const role_key = "roles"
	roleFnc := func(roles ...string) *core.Metadata {
		return core.SetMetadata(role_key, roles)
	}

	composite := core.Composition().Guard(guard).Pipe(core.Query(FilterDto{})).Metadata(roleFnc("admin"))
	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Composition(composite).Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": ctx.GetMetadata(role_key),
			})
		})

		return ctrl
	}

	module := func() *core.DynamicModule {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controller{controller},
		})

		return appModule
	}

	app := core.CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()
	resp, err := testClient.Get(testServer.URL + "/api/test?name=abc&key=value")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, []interface{}{"admin"}, res.Data)

	resp, err = testClient.Get(testServer.URL + "/api/test?name=abc")
	require.Nil(t, err)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test?key=value")
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
