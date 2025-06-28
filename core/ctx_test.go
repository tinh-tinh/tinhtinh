package core_test

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/cookie"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/session"
)

func Test_Ctx_Req(t *testing.T) {
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": ctx.Req().Host,
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

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, strings.Replace(testServer.URL, "http://", "", 1), res.Data)
}

func Test_Ctx_Res(t *testing.T) {
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			ctx.Res().Header().Set("key", "value")
			return ctx.JSON(core.Map{
				"data": "ok",
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

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	require.Equal(t, "value", resp.Header.Get("key"))
}

func Test_Ctx_Headers(t *testing.T) {
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": ctx.Headers("x-key"),
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

	req, err := http.NewRequest("GET", testServer.URL+"/api/test", nil)
	require.Nil(t, err)

	req.Header.Set("x-key", "value")

	resp, err := testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "value", res.Data)
}

func Test_Ctx_BodyParser(t *testing.T) {
	type BodyData struct {
		Name string `json:"name"`
	}
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Post("", func(ctx core.Ctx) error {
			var bodyData BodyData
			err := ctx.BodyParser(&bodyData)
			if err != nil {
				return common.InternalServerException(ctx.Res(), err.Error())
			}
			return ctx.JSON(core.Map{
				"data": bodyData.Name,
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

	resp, err := testClient.Post(testServer.URL+"/api/test", "application/json", strings.NewReader(`{"name":"test"}`))

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "test", res.Data)
}

func Test_Ctx_Body(t *testing.T) {
	type BodyData struct {
		Name string `json:"name"`
	}
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Pipe(core.Body(BodyData{})).Post("", func(ctx core.Ctx) error {
			data := ctx.Body().(*BodyData)
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

	resp, err := testClient.Post(testServer.URL+"/api/test", "application/json", strings.NewReader(`{"name":"test"}`))

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "test", res.Data)
}

func Test_Ctx_Params(t *testing.T) {
	type ID struct {
		ID string `path:"id"`
	}
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Pipe(core.Param(ID{})).Get("/{id}", func(ctx core.Ctx) error {
			data := ctx.Paths().(*ID)
			return ctx.JSON(core.Map{
				"data": data.ID,
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

	resp, err := testClient.Get(testServer.URL + "/api/test/123")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "123", res.Data)
}

func Test_Ctx_Queries(t *testing.T) {
	type QueryData struct {
		Name string `query:"name"`
		Age  uint   `query:"age"`
	}
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Pipe(core.Query(QueryData{})).Get("", func(ctx core.Ctx) error {
			data := ctx.Queries().(*QueryData)
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

	resp, err := testClient.Get(testServer.URL + "/api/test?name=test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "test", res.Data)
}

func Test_Ctx_Param(t *testing.T) {
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("/{id}", func(ctx core.Ctx) error {
			data := ctx.Path("id")
			return ctx.JSON(core.Map{
				"data": data,
			})
		})

		ctrl.Get("key/{key}", func(ctx core.Ctx) error {
			data := ctx.PathInt("key", 1)
			return ctx.JSON(core.Map{
				"data": data,
			})
		})

		ctrl.Get("key2/{key}", func(ctx core.Ctx) error {
			data := ctx.PathFloat("key", 1)
			return ctx.JSON(core.Map{
				"data": data,
			})
		})

		ctrl.Get("key3/{key}", func(ctx core.Ctx) error {
			data := ctx.PathBool("key", true)
			return ctx.JSON(core.Map{
				"data": data,
			})
		})

		ctrl.Get("key4/{key}", func(ctx core.Ctx) error {
			data := ctx.PathInt("key")
			return ctx.JSON(core.Map{
				"data": data,
			})
		})

		ctrl.Get("key5/{key}", func(ctx core.Ctx) error {
			data := ctx.PathFloat("key")
			return ctx.JSON(core.Map{
				"data": data,
			})
		})

		ctrl.Get("key6/{key}", func(ctx core.Ctx) error {
			data := ctx.PathBool("key")
			return ctx.JSON(core.Map{
				"data": data,
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

	resp, err := testClient.Get(testServer.URL + "/api/test/123")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "123", res.Data)

	// Case ParamInt
	resp2, err := testClient.Get(testServer.URL + "/api/test/key/456")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	data2, err := io.ReadAll(resp2.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data2, &res)
	require.Nil(t, err)
	require.Equal(t, float64(456), res.Data)

	resp2, err = testClient.Get(testServer.URL + "/api/test/key/abc")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	// Case ParamFloat
	resp3, err := testClient.Get(testServer.URL + "/api/test/key2/10.84573984573984")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp3.StatusCode)

	data3, err := io.ReadAll(resp3.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data3, &res)
	require.Nil(t, err)
	require.Equal(t, 10.84573984573984, res.Data)

	resp3, err = testClient.Get(testServer.URL + "/api/test/key2/abc")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp3.StatusCode)

	// Case ParamBool
	resp4, err := testClient.Get(testServer.URL + "/api/test/key3/true")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp4.StatusCode)

	data4, err := io.ReadAll(resp4.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data4, &res)
	require.Nil(t, err)
	require.Equal(t, true, res.Data)

	resp4, err = testClient.Get(testServer.URL + "/api/test/key3/abc")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp4.StatusCode)

	// Case Param with default value
	resp5, err := testClient.Get(testServer.URL + "/api/test/key4/abc")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp5.StatusCode)

	data5, err := io.ReadAll(resp5.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data5, &res)
	require.Nil(t, err)
	require.Equal(t, 0.0, res.Data)

	// Case ParamFloat without default value
	resp6, err := testClient.Get(testServer.URL + "/api/test/key5/abc")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp6.StatusCode)

	data6, err := io.ReadAll(resp6.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data6, &res)
	require.Nil(t, err)
	require.Equal(t, 0.0, res.Data)

	// Case ParamBool
	resp7, err := testClient.Get(testServer.URL + "/api/test/key6/78")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp7.StatusCode)

	data7, err := io.ReadAll(resp7.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data7, &res)
	require.Nil(t, err)
	require.Equal(t, false, res.Data)
}

func Test_Ctx_Query(t *testing.T) {
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			data := ctx.Query("name")
			return ctx.JSON(core.Map{
				"data": data,
			})
		})

		ctrl.Get("key", func(ctx core.Ctx) error {
			data := ctx.QueryFloat("key", 10)
			return ctx.JSON(core.Map{
				"data": data,
			})
		})

		ctrl.Get("key2", func(ctx core.Ctx) error {
			data := ctx.QueryFloat("key")
			return ctx.JSON(core.Map{
				"data": data,
			})
		})

		ctrl.Get("key3", func(ctx core.Ctx) error {
			data := ctx.QueryBool("key")
			return ctx.JSON(core.Map{
				"data": data,
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

	resp, err := testClient.Get(testServer.URL + "/api/test?name=test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "test", res.Data)

	// Case QueryInt
	resp2, err := testClient.Get(testServer.URL + "/api/test/key?key=123")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	data2, err := io.ReadAll(resp2.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data2, &res)
	require.Nil(t, err)
	require.Equal(t, float64(123), res.Data)

	resp2, err = testClient.Get(testServer.URL + "/api/test/key?key=abc")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	// Case QueryFloat
	resp3, err := testClient.Get(testServer.URL + "/api/test/key2?key=10.84573984573984")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp3.StatusCode)

	data3, err := io.ReadAll(resp3.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data3, &res)
	require.Nil(t, err)
	require.Equal(t, 10.84573984573984, res.Data)

	// Case QueryInt with invalid integer
	resp5, err := testClient.Get(testServer.URL + "/api/test/key?key=invalid")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp5.StatusCode)

	data5, err := io.ReadAll(resp5.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data5, &res)
	require.Nil(t, err)
	require.Equal(t, float64(10), res.Data)

	// Case QueryFloat with invalid float
	resp6, err := testClient.Get(testServer.URL + "/api/test/key2?key=invalid")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp6.StatusCode)

	data6, err := io.ReadAll(resp6.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data6, &res)
	require.Nil(t, err)
	require.Equal(t, 0.0, res.Data)

	// Case QueryBool with invalid boolean
	resp7, err := testClient.Get(testServer.URL + "/api/test/key3?key=invalid")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp7.StatusCode)

	data7, err := io.ReadAll(resp7.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data7, &res)
	require.Nil(t, err)
	require.Equal(t, false, res.Data)

	// Case QueryBool
	resp4, err := testClient.Get(testServer.URL + "/api/test/key3?key=true")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp4.StatusCode)

	data4, err := io.ReadAll(resp4.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data4, &res)
	require.Nil(t, err)
	require.Equal(t, true, res.Data)

	resp4, err = testClient.Get(testServer.URL + "/api/test/key3?key=abc")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp4.StatusCode)
}

func Test_Ctx_QueryInt(t *testing.T) {
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			data := ctx.QueryInt("name")
			return ctx.JSON(core.Map{
				"data": data,
			})
		})

		ctrl.Get("d", func(ctx core.Ctx) error {
			data := ctx.QueryInt("default", 10)
			return ctx.JSON(core.Map{
				"data": data,
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

	resp, err := testClient.Get(testServer.URL + "/api/test?name=10")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	type IntResponse struct {
		Data int `json:"data"`
	}

	var res IntResponse
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, 10, res.Data)

	resp, err = testClient.Get(testServer.URL + "/api/test?name=abc")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, 0, res.Data)

	resp, err = testClient.Get(testServer.URL + "/api/test/d?default=5")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, 5, res.Data)

	resp, err = testClient.Get(testServer.URL + "/api/test/d")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, 10, res.Data)
}

func Test_Ctx_QueryBool(t *testing.T) {
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			data := ctx.QueryBool("name")
			return ctx.JSON(core.Map{
				"data": data,
			})
		})

		ctrl.Get("/d", func(ctx core.Ctx) error {
			data := ctx.QueryBool("default", true)
			return ctx.JSON(core.Map{
				"data": data,
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

	resp, err := testClient.Get(testServer.URL + "/api/test?name=true")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	type BoolResponse struct {
		Data bool `json:"data"`
	}

	var res BoolResponse
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, true, res.Data)

	resp, err = testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, false, res.Data)

	resp, err = testClient.Get(testServer.URL + "/api/test/d")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, true, res.Data)

	resp, err = testClient.Get(testServer.URL + "/api/test/d?default=false")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, false, res.Data)
}

func Test_Ctx_Status(t *testing.T) {
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.Status(http.StatusNotFound).JSON(core.Map{
				"data": "ok",
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

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func Test_QueryParser(t *testing.T) {
	type QueryData struct {
		Age    uint    `query:"age"`
		Score  float32 `query:"score"`
		Format bool    `query:"format"`
	}

	type UnsupportStruct struct {
		hidden string         `query:"hidden"`
		Untype map[string]any `query:"untype"`
	}
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			var queryData QueryData
			err := ctx.QueryParser(&queryData)
			if err != nil {
				return common.InternalServerException(ctx.Res(), err.Error())
			}
			return ctx.JSON(core.Map{
				"data": queryData,
			})
		})

		ctrl.Get("unsupport", func(ctx core.Ctx) error {
			var queryData UnsupportStruct
			err := ctx.QueryParser(&queryData)
			if err != nil {
				return common.InternalServerException(ctx.Res(), err.Error())
			}
			return ctx.JSON(core.Map{
				"data": queryData.hidden,
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

	resp, err := testClient.Get(testServer.URL + "/api/test?age=12&format=true&score=4.5")

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, map[string]interface{}(map[string]interface{}{"Age": float64(12), "Format": true, "Score": 4.5}), res.Data)

	resp, err = testClient.Get(testServer.URL + "/api/test?age=true")
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test?format=35")
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test?score=string")
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/unsupport?untype=string")
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/unsupport?hidden=string")
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_ParamParser(t *testing.T) {
	type ParamData struct {
		ID     int  `path:"id"`
		Export bool `path:"export"`
	}
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("{id}/{export}", func(ctx core.Ctx) error {
			var queryData ParamData
			err := ctx.PathParser(&queryData)
			if err != nil {
				return common.InternalServerException(ctx.Res(), err.Error())
			}
			return ctx.JSON(core.Map{
				"data": queryData.ID,
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

	resp, err := testClient.Get(testServer.URL + "/api/test/345/true")

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, float64(345), res.Data)
}

func Test_Ctx_Session(t *testing.T) {
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Post("", func(ctx core.Ctx) error {
			ctx.Session("key", "val")

			return ctx.JSON(core.Map{
				"data": "ok",
			})
		})

		ctrl.Get("", func(ctx core.Ctx) error {
			data := ctx.Session("key")
			return ctx.JSON(core.Map{
				"data": data,
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

	session := session.New(session.Options{
		Secret: "secret",
	})
	app := core.CreateFactory(module, core.AppOptions{
		Session: session,
	})
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Post(testServer.URL+"/api/test", "application/json", nil)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	req, err := http.NewRequest("GET", testServer.URL+"/api/test", nil)
	require.Nil(t, err)

	req.AddCookie(resp.Cookies()[0])
	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	type StringResponse struct {
		Data string `json:"data"`
	}

	var res StringResponse
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "val", res.Data)
}

func Test_Cookie(t *testing.T) {
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Post("", func(ctx core.Ctx) error {
			ctx.SetCookie("key", "val", 3600)

			return ctx.JSON(core.Map{
				"data": "ok",
			})
		})

		ctrl.Get("", func(ctx core.Ctx) error {
			data := ctx.Cookies("key").Value
			return ctx.JSON(core.Map{
				"data": data,
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

	resp, err := testClient.Post(testServer.URL+"/api/test", "application/json", nil)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	req, err := http.NewRequest("GET", testServer.URL+"/api/test", nil)
	require.Nil(t, err)
	req.AddCookie(resp.Cookies()[0])

	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "val", res.Data)
}

func Test_SignedCookie(t *testing.T) {
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Post("", func(ctx core.Ctx) error {
			_, err := ctx.SignedCookie("key", "val")
			if err != nil {
				return common.InternalServerException(ctx.Res(), err.Error())
			}

			return ctx.JSON(core.Map{
				"data": "ok",
			})
		})

		ctrl.Get("", func(ctx core.Ctx) error {
			data, err := ctx.SignedCookie("key")
			if err != nil {
				return common.InternalServerException(ctx.Res(), err.Error())
			}
			return ctx.JSON(core.Map{
				"data": data,
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
	app.Use(cookie.Handler(cookie.Options{
		Key: "abc&1*~#^2^#s0^=)^^7%b34",
	}))

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	resp, err = testClient.Post(testServer.URL+"/api/test", "application/json", nil)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	req, err := http.NewRequest("GET", testServer.URL+"/api/test", nil)
	require.Nil(t, err)
	req.AddCookie(resp.Cookies()[0])

	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "val", res.Data)
}

func Test_Redirect(t *testing.T) {
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("/redirect", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "ok",
			})
		})

		ctrl.Get("/fail", func(ctx core.Ctx) error {
			return ctx.Redirect("!@#$%^&**&^%$#redirect")
		})

		ctrl.Get("/out", func(ctx core.Ctx) error {
			return ctx.Redirect("https://www.google.com")
		})

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.Redirect("/redirect")
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

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
	require.Equal(t, "ok", res.Data)

	resp, err = testClient.Get(testServer.URL + "/api/test/fail")
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/out")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_ArrayQuery(t *testing.T) {
	type QueryData struct {
		IdsInt     []int     `query:"idsInt" json:"idsInt"`
		IdsInt8    []int8    `query:"idsInt8" json:"idsInt8"`
		IdsInt16   []int16   `query:"idsInt16" json:"idsInt16"`
		IdsInt32   []int32   `query:"idsInt32" json:"idsInt32"`
		IdsInt64   []int64   `query:"idsInt64" json:"idsInt64"`
		IdsUInt    []uint    `query:"idsUInt" json:"idsUInt"`
		IdsUInt8   []uint8   `query:"idsUInt8" json:"idsUInt8"`
		IdsUInt16  []uint16  `query:"idsUInt16" json:"idsUInt16"`
		IdsUInt32  []uint32  `query:"idsUInt32" json:"idsUInt32"`
		IdsUInt64  []uint64  `query:"idsUInt64" json:"idsUInt64"`
		IdsFloat32 []float32 `query:"idsFloat32" json:"idsFloat32"`
		IdsFloat64 []float64 `query:"idsFloat64" json:"idsFloat64"`
		IdsStr     []string  `query:"idsStr" json:"idsStr"`
		IdsBool    []bool    `query:"idsBool" json:"idsBool"`
	}

	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			var queryData QueryData
			err := ctx.QueryParser(&queryData)
			if err != nil {
				return common.InternalServerException(ctx.Res(), err.Error())
			}

			return ctx.JSON(core.Map{
				"data": queryData,
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
	tests := []struct {
		field  string
		values []string
		expect any
		isErr  bool
	}{
		{"idsInt", []string{"1", "2"}, []any{float64(1), float64(2)}, false}, // expect error
		{"idsInt8", []string{"1", "2"}, []any{float64(1), float64(2)}, false},
		{"idsInt16", []string{"1", "2"}, []any{float64(1), float64(2)}, false},
		{"idsInt32", []string{"1", "2"}, []any{float64(1), float64(2)}, false},
		{"idsInt64", []string{"1", "2"}, []any{float64(1), float64(2)}, false},
		{"idsUInt", []string{"3", "4"}, []any{float64(3), float64(4)}, false},
		{"idsUInt16", []string{"3", "4"}, []any{float64(3), float64(4)}, false},
		{"idsUInt32", []string{"3", "4"}, []any{float64(3), float64(4)}, false},
		{"idsUInt64", []string{"3", "4"}, []any{float64(3), float64(4)}, false},
		{"idsFloat32", []string{"1.5", "2.5"}, []any{1.5, 2.5}, false},
		{"idsFloat64", []string{"1.5", "2.5"}, []any{1.5, 2.5}, false},
		{"idsStr", []string{"abc", "def"}, []any{"abc", "def"}, false},
		{"idsBool", []string{"true", "false"}, []any{true, false}, false},
		{"idsInt", []string{"1", "true"}, []any{float64(1), float64(2)}, true}, // expect error
		{"idsInt8", []string{"1", "2true"}, []any{float64(1), float64(2)}, true},
		{"idsInt16", []string{"1", "2true"}, []any{float64(1), float64(2)}, true},
		{"idsInt32", []string{"1", "2true"}, []any{float64(1), float64(2)}, true},
		{"idsInt64", []string{"1", "2true"}, []any{float64(1), float64(2)}, true},
		{"idsUInt", []string{"3", "4true"}, []any{float64(3), float64(4)}, true},
		{"idsUInt16", []string{"3", "4true"}, []any{float64(3), float64(4)}, true},
		{"idsUInt32", []string{"3", "4true"}, []any{float64(3), float64(4)}, true},
		{"idsUInt64", []string{"3", "4true"}, []any{float64(3), float64(4)}, true},
		{"idsFloat32", []string{"1.5", "2.5true"}, []any{1.5, 2.5}, true},
		{"idsFloat64", []string{"1.5", "2.5true"}, []any{1.5, 2.5}, true},
		{"idsBool", []string{"true", "25"}, []any{true, false}, true},
	}

	for _, test := range tests {
		url := testServer.URL + "/api/test?"
		for _, v := range test.values {
			url += fmt.Sprintf("%s=%s&", test.field, v)
		}
		resp, err := testClient.Get(url)
		require.Nil(t, err)
		if test.isErr {
			require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			continue
		}
		require.Equal(t, http.StatusOK, resp.StatusCode)

		data, err := io.ReadAll(resp.Body)
		require.Nil(t, err)

		var res Response
		err = json.Unmarshal(data, &res)
		require.Nil(t, err)

		mapper := res.Data.(map[string]any)
		require.NotNil(t, mapper[test.field])
		require.Equal(t, test.expect, mapper[test.field], fmt.Sprintf("Failed by case %s", test.field))

		// resp, err = testClient.Get(testServer.URL + "/api/test?idsInt=true&idsInt=2")
		// require.Nil(t, err)
		// require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	}
}

func Test_XML(t *testing.T) {
	type User struct {
		XMLName xml.Name `xml:"user"`
		Name    string   `xml:"name"`
		Email   string   `xml:"email"`
	}

	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			user := User{Name: "Alice", Email: "alice@example.com"}

			return ctx.Status(http.StatusOK).XML(user)
		})

		ctrl.Get("error", func(ctx core.Ctx) error {
			type Bad struct {
				C chan int // channels cannot be marshaled
			}
			bad := Bad{C: make(chan int)}
			return ctx.Status(500).XML(bad)
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
	resp, err := testClient.Get(testServer.URL + "/api/test/error")
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var user User
	err = xml.Unmarshal(data, &user)
	require.Nil(t, err)
	require.Equal(t, "Alice", user.Name)
	require.Equal(t, "alice@example.com", user.Email)
}

func createHTMLFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func Test_Render(t *testing.T) {
	// Prepare a simple HTML content
	layoutContent := `
<!DOCTYPE html>
<html>
<head>
    <title>{{block "title" .}}Default Title{{end}}</title>
</head>
<body>
    {{block "content" .}}Default Content{{end}}
</body>
</html>`
	err := createHTMLFile("layout.html", layoutContent)
	require.Nil(t, err)

	homeContent := `
{{define "title"}}Home Page{{end}}
{{define "content"}}
<h1>Welcome, {{.Name}}!</h1>
{{end}}`

	err = createHTMLFile("home.html", homeContent)
	require.Nil(t, err)

	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.Render("layout.html", core.Map{
				"Name": ctx.Query("name"),
			}, "layout.html", "home.html")
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
	resp, err := testClient.Get(testServer.URL + "/api/test?name=John")
	require.Nil(t, err)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	res := string(data)
	require.NotEmpty(t, res)
	require.Contains(t, res, "<h1>Welcome, John!</h1>")
	require.Contains(t, res, "<title>Home Page</title>")

	// Test with missing template files
	os.Remove("layout.html")
	os.Remove("home.html")

	resp, err = testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	// require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)

	var resMap core.Map
	err = json.Unmarshal(data, &resMap)
	require.Nil(t, err)
	require.Equal(t, float64(500), resMap["statusCode"])
}

func Benchmark_CtxJson(b *testing.B) {
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.Status(http.StatusOK).JSON(core.Map{
				"data": "ok",
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
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			resp, err := testClient.Get(testServer.URL + "/api/test")
			require.Nil(b, err)
			require.Equal(b, http.StatusOK, resp.StatusCode)
		}
	})
}
