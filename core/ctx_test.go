package core_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/common"
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/middleware/cookie"
	"github.com/tinh-tinh/tinhtinh/middleware/session"
)

func Test_Ctx_Req(t *testing.T) {
	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": ctx.Req().Host,
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
	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			ctx.Res().Header().Set("key", "value")
			return ctx.JSON(core.Map{
				"data": "ok",
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

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	require.Equal(t, "value", resp.Header.Get("key"))
}

func Test_Ctx_Headers(t *testing.T) {
	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": ctx.Headers("x-key"),
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
	controller := func(module *core.DynamicModule) *core.DynamicController {
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
	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Pipe(core.Body(BodyData{})).Post("", func(ctx core.Ctx) error {
			data := ctx.Body().(*BodyData)
			return ctx.JSON(core.Map{
				"data": data.Name,
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
		ID string `param:"id"`
	}
	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Pipe(core.Param(ID{})).Get("/{id}", func(ctx core.Ctx) error {
			data := ctx.Params().(*ID)
			return ctx.JSON(core.Map{
				"data": data.ID,
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
	}
	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Pipe(core.Query(QueryData{})).Get("", func(ctx core.Ctx) error {
			data := ctx.Queries().(*QueryData)
			return ctx.JSON(core.Map{
				"data": data.Name,
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
	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("/{id}", func(ctx core.Ctx) error {
			data := ctx.Param("id")
			return ctx.JSON(core.Map{
				"data": data,
			})
		})

		ctrl.Get("key/{key}", func(ctx core.Ctx) error {
			data := ctx.ParamInt("key")
			return ctx.JSON(core.Map{
				"data": data,
			})
		})

		ctrl.Get("key2/{key}", func(ctx core.Ctx) error {
			data := ctx.ParamFloat("key")
			return ctx.JSON(core.Map{
				"data": data,
			})
		})

		ctrl.Get("key3/{key}", func(ctx core.Ctx) error {
			data := ctx.ParamBool("key")
			return ctx.JSON(core.Map{
				"data": data,
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
	require.Equal(t, http.StatusInternalServerError, resp2.StatusCode)

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
	require.Equal(t, http.StatusInternalServerError, resp3.StatusCode)

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
	require.Equal(t, http.StatusInternalServerError, resp4.StatusCode)
}

func Test_Ctx_Query(t *testing.T) {
	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			data := ctx.Query("name")
			return ctx.JSON(core.Map{
				"data": data,
			})
		})

		ctrl.Get("key", func(ctx core.Ctx) error {
			data := ctx.QueryInt("key")
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
	require.Equal(t, http.StatusInternalServerError, resp2.StatusCode)

	// Case QueryFloat
	resp3, err := testClient.Get(testServer.URL + "/api/test/key2?key=10.84573984573984")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp3.StatusCode)

	data3, err := io.ReadAll(resp3.Body)
	require.Nil(t, err)

	err = json.Unmarshal(data3, &res)
	require.Nil(t, err)
	require.Equal(t, 10.84573984573984, res.Data)

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
	require.Equal(t, http.StatusInternalServerError, resp4.StatusCode)
}

func Test_Ctx_QueryInt(t *testing.T) {
	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			data := ctx.QueryInt("name")
			return ctx.JSON(core.Map{
				"data": data,
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
}

func Test_Ctx_QueryBool(t *testing.T) {
	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			data := ctx.QueryBool("name")
			return ctx.JSON(core.Map{
				"data": data,
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
}

func Test_Ctx_Status(t *testing.T) {
	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.Status(http.StatusNotFound).JSON(core.Map{
				"data": "ok",
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

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func Test_QueryParser(t *testing.T) {
	type QueryData struct {
		Age    int  `query:"age"`
		Format bool `query:"format"`
	}
	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			var queryData QueryData
			err := ctx.QueryParse(&queryData)
			if err != nil {
				return common.InternalServerException(ctx.Res(), err.Error())
			}
			return ctx.JSON(core.Map{
				"data": queryData,
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

	resp, err := testClient.Get(testServer.URL + "/api/test?age=12&format=true")

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	err = json.Unmarshal(data, &res)
	require.Nil(t, err)
}

func Test_ParamParser(t *testing.T) {
	type ParamData struct {
		ID     int  `param:"id"`
		Export bool `param:"export"`
	}
	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("{id}/{export}", func(ctx core.Ctx) error {
			var queryData ParamData
			err := ctx.ParamParse(&queryData)
			if err != nil {
				return common.InternalServerException(ctx.Res(), err.Error())
			}
			return ctx.JSON(core.Map{
				"data": queryData.ID,
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
	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Post("", func(ctx core.Ctx) error {
			ctx.Session("key", "val")

			return ctx.JSON(core.Map{
				"data": "ok",
			})
		})

		ctrl.Get("", func(ctx core.Ctx) error {
			data := ctx.Session("key")
			fmt.Print(data)
			return ctx.JSON(core.Map{
				"data": data,
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
	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Post("", func(ctx core.Ctx) error {
			ctx.SetCookie("key", "val", 3600)

			return ctx.JSON(core.Map{
				"data": "ok",
			})
		})

		ctrl.Get("", func(ctx core.Ctx) error {
			data := ctx.Cookies("key").Value
			fmt.Print(data)
			return ctx.JSON(core.Map{
				"data": data,
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
	controller := func(module *core.DynamicModule) *core.DynamicController {
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

	module := func() *core.DynamicModule {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controller{controller},
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
	controller := func(module *core.DynamicModule) *core.DynamicController {
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
