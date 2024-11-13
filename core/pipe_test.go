package core_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/core"
)

func Test_PipeMiddleware(t *testing.T) {
	type SignUpDto struct {
		Name     string `validate:"required"`
		Email    string `validate:"required,isEmail"`
		Password string `validate:"isStrongPassword"`
		Age      int    `validate:"isInt"`
	}
	appController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Pipe(core.Body(&SignUpDto{})).Post("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "2",
			})
		})

		return ctrl
	}

	module := func() *core.DynamicModule {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controller{appController},
		})

		return appModule
	}

	app := core.CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Post(testServer.URL+"/api/test", "application/json", strings.NewReader(`{"name":"test", "email":"test", "password":"test"}`))
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	resp, err = testClient.Post(testServer.URL+"/api/test", "application/json", strings.NewReader(`{"name":"test", "email":"test@gmail.com", "password":"Test@1234546"}`))
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = testClient.Post(testServer.URL+"/api/test", "application/json", strings.NewReader(`{"name":"test", "email":"test@gmail.com", "password":"Test@1234546", "age": "haha"}`))
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	resp, err = testClient.Post(testServer.URL+"/api/test", "application/json", strings.NewReader(`{"name":"test", "email":"test@gmail.com", "password":"Test@1234546", "age":333}`))
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_Query(t *testing.T) {
	type FilterDto struct {
		Name  string `validate:"required" query:"name"`
		Email string `validate:"required,isEmail" query:"email"`
		Age   int    `validate:"isInt" query:"age"`
	}
	appController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Pipe(core.Query(&FilterDto{})).Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "2",
			})
		})

		return ctrl
	}

	module := func() *core.DynamicModule {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controller{appController},
		})

		return appModule
	}

	app := core.CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test?name=test&email=test")
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test?name=test&email=test@gmail.com&age=g")
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test?name=test&email=test@gmail.com")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test?name=test&email=test@gmail.com&age=12")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_Param(t *testing.T) {
	type ParamDto struct {
		ID int `validate:"required,isInt" param:"id"`
	}
	appController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Pipe(core.Param(&ParamDto{})).Get("{id}", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "2",
			})
		})

		return ctrl
	}

	module := func() *core.DynamicModule {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controller{appController},
		})

		return appModule
	}

	app := core.CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test/abc")
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/123")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
