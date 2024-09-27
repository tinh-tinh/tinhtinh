package core

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Ctx_Req(t *testing.T) {
	controller := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": ctx.Req().Host,
			})
		})

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{controller},
		})

		return appModule
	}

	app := CreateFactory(module, "api")

	testServer := httptest.NewServer(app.prepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	json.Unmarshal(data, &res)
	require.Equal(t, strings.Replace(testServer.URL, "http://", "", 1), res.Data)
}

func Test_Ctx_Res(t *testing.T) {
	controller := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx Ctx) {
			ctx.Res().Header().Set("key", "value")
			ctx.JSON(Map{
				"data": "ok",
			})
		})

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{controller},
		})

		return appModule
	}

	app := CreateFactory(module, "api")

	testServer := httptest.NewServer(app.prepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	require.Equal(t, "value", resp.Header.Get("key"))
}

func Test_Ctx_Headers(t *testing.T) {
	controller := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": ctx.Headers("x-key"),
			})
		})

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{controller},
		})

		return appModule
	}

	app := CreateFactory(module, "api")

	testServer := httptest.NewServer(app.prepareBeforeListen())
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
	json.Unmarshal(data, &res)
	require.Equal(t, "value", res.Data)
}

func Test_Ctx_BodyParser(t *testing.T) {
	type BodyData struct {
		Name string `json:"name"`
	}
	controller := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Post("", func(ctx Ctx) {
			var bodyData *BodyData
			ctx.BodyParser(&bodyData)
			ctx.JSON(Map{
				"data": bodyData.Name,
			})
		})

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{controller},
		})

		return appModule
	}

	app := CreateFactory(module, "api")

	testServer := httptest.NewServer(app.prepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Post(testServer.URL+"/api/test", "application/json", strings.NewReader(`{"name":"test"}`))

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	json.Unmarshal(data, &res)
	require.Equal(t, "test", res.Data)
}

func Test_Ctx_Body(t *testing.T) {
	type BodyData struct {
		Name string `json:"name"`
	}
	controller := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Pipe(Body(&BodyData{})).Post("", func(ctx Ctx) {
			data := ctx.Body().(*BodyData)
			ctx.JSON(Map{
				"data": data.Name,
			})
		})

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{controller},
		})

		return appModule
	}

	app := CreateFactory(module, "api")

	testServer := httptest.NewServer(app.prepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Post(testServer.URL+"/api/test", "application/json", strings.NewReader(`{"name":"test"}`))

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	json.Unmarshal(data, &res)
	require.Equal(t, "test", res.Data)
}

func Test_Ctx_Params(t *testing.T) {
	type ID struct {
		ID string `param:"id"`
	}
	controller := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Pipe(Param(&ID{})).Get("/{id}", func(ctx Ctx) {
			data := ctx.Params().(*ID)
			ctx.JSON(Map{
				"data": data.ID,
			})
		})

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{controller},
		})

		return appModule
	}

	app := CreateFactory(module, "api")

	testServer := httptest.NewServer(app.prepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test/123")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	json.Unmarshal(data, &res)
	require.Equal(t, "123", res.Data)
}

func Test_Ctx_Queries(t *testing.T) {
	type QueryData struct {
		Name string `query:"name"`
	}
	controller := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Pipe(Query(&QueryData{})).Get("", func(ctx Ctx) {
			data := ctx.Queries().(*QueryData)
			ctx.JSON(Map{
				"data": data.Name,
			})
		})

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{controller},
		})

		return appModule
	}

	app := CreateFactory(module, "api")

	testServer := httptest.NewServer(app.prepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test?name=test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	json.Unmarshal(data, &res)
	require.Equal(t, "test", res.Data)
}

func Test_Ctx_Param(t *testing.T) {
	controller := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("/{id}", func(ctx Ctx) {
			data := ctx.Param("id")
			ctx.JSON(Map{
				"data": data,
			})
		})

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{controller},
		})

		return appModule
	}

	app := CreateFactory(module, "api")

	testServer := httptest.NewServer(app.prepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test/123")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	json.Unmarshal(data, &res)
	require.Equal(t, "123", res.Data)
}

func Test_Ctx_Query(t *testing.T) {
	controller := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx Ctx) {
			data := ctx.Query("name")
			ctx.JSON(Map{
				"data": data,
			})
		})

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{controller},
		})

		return appModule
	}

	app := CreateFactory(module, "api")

	testServer := httptest.NewServer(app.prepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test?name=test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	json.Unmarshal(data, &res)
	require.Equal(t, "test", res.Data)
}

func Test_Ctx_QueryInt(t *testing.T) {
	controller := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx Ctx) {
			data := ctx.QueryInt("name")
			ctx.JSON(Map{
				"data": data,
			})
		})

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{controller},
		})

		return appModule
	}

	app := CreateFactory(module, "api")

	testServer := httptest.NewServer(app.prepareBeforeListen())
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
	json.Unmarshal(data, &res)
	require.Equal(t, 10, res.Data)
}

func Test_Ctx_QueryBool(t *testing.T) {
	controller := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx Ctx) {
			data := ctx.QueryBool("name")
			ctx.JSON(Map{
				"data": data,
			})
		})

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{controller},
		})

		return appModule
	}

	app := CreateFactory(module, "api")

	testServer := httptest.NewServer(app.prepareBeforeListen())
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
	json.Unmarshal(data, &res)
	require.Equal(t, true, res.Data)
}

func Test_Ctx_Status(t *testing.T) {
	controller := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx Ctx) {
			ctx.Status(http.StatusNotFound)
		})

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{controller},
		})

		return appModule
	}

	app := CreateFactory(module, "api")

	testServer := httptest.NewServer(app.prepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func Test_CtxContext(t *testing.T) {
	const key CtxKey = "key"
	middleware := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), key, "value"))

			h.ServeHTTP(w, r)
		})
	}
	controller := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Use(middleware).Get("", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": ctx.Get(key),
			})
		})

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{controller},
		})

		return appModule
	}

	app := CreateFactory(module, "api")

	testServer := httptest.NewServer(app.prepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var res Response
	json.Unmarshal(data, &res)
	require.Equal(t, "value", res.Data)
}