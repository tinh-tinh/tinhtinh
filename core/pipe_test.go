package core_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func TestDefaultDto(t *testing.T) {
	type Pagination struct {
		Page  int `validate:"isInt" default:"1"`
		Limit int `validate:"isInt" default:"10"`
	}

	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Pipe(core.QueryParser[Pagination]{}).Get("", func(ctx core.Ctx) error {
			pagin := ctx.Queries().(*Pagination)
			return ctx.JSON(core.Map{
				"data": pagin.Page,
			})
		})

		return ctrl
	}

	module := func() core.Module {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{appController},
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
	require.Equal(t, `{"data":1}`, string(data))
}

func Test_Pipe(t *testing.T) {
	type SignUpDto struct {
		Name     string `validate:"required"`
		Email    string `validate:"required,isEmail"`
		Password string `validate:"isStrongPassword"`
		Age      int    `validate:"isInt"`
	}

	type FilterDto struct {
		Name  string `validate:"required" query:"name"`
		Email string `validate:"required,isEmail" query:"email"`
		Age   int    `validate:"isInt" query:"age"`
	}

	type ParamDto struct {
		ID int `validate:"required,isInt" path:"id"`
	}

	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.
			Pipe(core.BodyParser[SignUpDto]{}).
			Post("", func(ctx core.Ctx) error {
				return ctx.JSON(core.Map{
					"data": ctx.Body().(*SignUpDto),
				})
			})

		ctrl.
			Pipe(core.QueryParser[FilterDto]{}).
			Get("", func(ctx core.Ctx) error {
				return ctx.JSON(core.Map{
					"data": "2",
				})
			})

		ctrl.
			Pipe(core.PathParser[ParamDto]{}).
			Get("{id}", func(ctx core.Ctx) error {
				return ctx.JSON(core.Map{
					"data": "2",
				})
			})

		return ctrl
	}

	module := func() core.Module {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{appController},
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

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":{"Name":"test","Email":"test@gmail.com","Password":"Test@1234546","Age":0}}`, string(data))

	resp, err = testClient.Post(testServer.URL+"/api/test", "application/json", strings.NewReader(`{"name":"test", "email":"test@gmail.com", "password":"Test@1234546", "age": "haha"}`))
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	resp, err = testClient.Post(testServer.URL+"/api/test", "application/json", strings.NewReader(`{"name":"test", "email":"test@gmail.com", "password":"Test@1234546", "age":333}`))
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":{"Name":"test","Email":"test@gmail.com","Password":"Test@1234546","Age":333}}`, string(data))

	resp, err = testClient.Get(testServer.URL + "/api/test?name=test&email=test")
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

	resp, err = testClient.Get(testServer.URL + "/api/test/abc")
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/123")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func BenchmarkPipe(b *testing.B) {
	type SignUpDto struct {
		Name     string `validate:"required"`
		Email    string `validate:"required,isEmail"`
		Password string `validate:"isStrongPassword"`
		Age      int    `validate:"isInt"`
	}
	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Pipe(core.BodyParser[SignUpDto]{}).Post("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": ctx.Body().(*SignUpDto),
			})
		})

		return ctrl
	}

	module := func() core.Module {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{appController},
		})

		return appModule
	}

	app := core.CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			testClient := testServer.Client()
			resp, err := testClient.Post(testServer.URL+"/api/test", "application/json", strings.NewReader(`{"name":"test","email":"test@mailinator.com","password":"12345678@Test","age":1}`))
			require.Nil(b, err)
			require.Equal(b, http.StatusOK, resp.StatusCode)
		}
	})
}
