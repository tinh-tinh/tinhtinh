package core_test

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func Transform(ctx core.Ctx) core.CallHandler {
	fmt.Println("Before ...")
	now := time.Now()
	return func(data any) any {
		res := make(core.Map)
		for key, val := range data.(core.Map) {
			if val != nil {
				res[key] = val
			}
		}
		fmt.Printf("After ...%vns\n", time.Since(now).Nanoseconds())
		return res
	}
}

func Test_Interceptor(t *testing.T) {
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Interceptor(Transform).Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data":    "ok",
				"total":   10,
				"message": nil,
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
	require.Equal(t, `{"data":"ok","total":10}`, string(data))
}

func Test_ParseInterceptorModule(t *testing.T) {
	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data":    "ok",
				"total":   10,
				"message": nil,
			})
		})

		return ctrl
	}

	appModule := func() core.Module {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{appController},
			Interceptor: Transform,
		})

		return appModule
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":"ok","total":10}`, string(data))
}

func Test_InterceptorMultiApi(t *testing.T) {
	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test").Interceptor(Transform).Registry()

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data":    "ok",
				"total":   10,
				"message": nil,
			})
		})

		ctrl.Post("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data":    "ok",
				"total":   10,
				"message": nil,
			})
		})

		return ctrl
	}

	appModule := func() core.Module {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{appController},
		})

		return appModule
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":"ok","total":10}`, string(data))

	resp, err = testClient.Post(testServer.URL+"/api/test", "application/json", nil)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":"ok","total":10}`, string(data))
}

func Test_InterceptorWilXML(t *testing.T) {
	type XMLData struct {
		Data    string `xml:"data"`
		Total   int    `xml:"total"`
		Message any    `xml:"message"`
	}

	TransformXML := func(ctx core.Ctx) core.CallHandler {
		fmt.Println("Before XML ...")
		now := time.Now()
		return func(data any) any {
			res := XMLData{}
			for key, val := range data.(core.Map) {
				if val != nil {
					switch key {
					case "data":
						res.Data = val.(string)
					case "total":
						res.Total = val.(int)
					case "message":
						res.Message = val
					}
				}
			}
			fmt.Printf("After XML ...%vns\n", time.Since(now).Nanoseconds())
			return res
		}
	}

	appController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Interceptor(TransformXML).Get("", func(ctx core.Ctx) error {
			return ctx.XML(core.Map{
				"data":    "ok",
				"total":   10,
				"message": nil,
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

	var xmlData XMLData
	err = xml.Unmarshal(data, &xmlData)
	require.Nil(t, err)
	require.Equal(t, "ok", xmlData.Data)
	require.Equal(t, 10, xmlData.Total)
	require.Nil(t, xmlData.Message)
}
