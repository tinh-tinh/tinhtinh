package core_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/common/exception"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

const (
	Context core.CtxKey = "context"
)

func Test_Execution(t *testing.T) {
	type Abc struct {
		Name string
	}
	middleware := func(ctx core.Ctx) error {
		abc := &Abc{
			Name: ctx.Query("name"),
		}
		ctx.Set(Context, abc)
		return ctx.Next()
	}

	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Use(middleware).Get("", func(ctx core.Ctx) error {
			data := core.Execution[Abc](Context, ctx)
			return ctx.JSON(core.Map{
				"data": data.Name,
			})
		})

		ctrl.Get("none", func(ctx core.Ctx) error {
			data := core.Execution[Abc](Context, ctx)
			if data == nil {
				return exception.NotFound("not found ctx")
			}
			return ctx.JSON(core.Map{
				"data": data.Name,
			})
		})

		return ctrl
	}

	module := func() core.Module {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{controller},
		})

		return appModule
	}

	app := core.CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test?name=haha")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	type Response struct {
		Data string `json:"data"`
	}

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "haha", res.Data)

	resp, err = testClient.Get(testServer.URL + "/api/test/none")
	require.Nil(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
