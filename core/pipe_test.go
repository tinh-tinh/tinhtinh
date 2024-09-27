package core

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_PipeMiddleware(t *testing.T) {
	type SignUpDto struct {
		Name     string `validate:"required"`
		Email    string `validate:"required,isEmail"`
		Password string `validate:"isStrongPassword"`
	}
	appController := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "1",
			})
		})

		ctrl.Pipe(Body(&SignUpDto{})).Post("", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "2",
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

	testServer := httptest.NewServer(app.prepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Post(testServer.URL+"/api/test", "application/json", strings.NewReader(`{"name":"test", "email":"test", "password":"test"}`))
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	resp, err = testClient.Post(testServer.URL+"/api/test", "application/json", strings.NewReader(`{"name":"test", "email":"test@gmail.com", "password":"Test@1234546"}`))
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_Query(t *testing.T) {
	type FilterDto struct {
		Name  string `validate:"required" query:"name"`
		Email string `validate:"required,isEmail" query:"email"`
	}
	appController := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Pipe(Query(&FilterDto{})).Get("", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "2",
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

	testServer := httptest.NewServer(app.prepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test?name=test&email=test")
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test?name=test&email=test@gmail.com")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_Param(t *testing.T) {
	type ParamDto struct {
		ID string `validate:"required,isInt" param:"id"`
	}
	appController := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Pipe(Param(&ParamDto{})).Get("{id}", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "2",
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

	testServer := httptest.NewServer(app.prepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test/abc")
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/123")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
