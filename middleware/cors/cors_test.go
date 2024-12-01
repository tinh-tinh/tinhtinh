package cors_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/middleware/cors"
)

func appModule() *core.DynamicModule {
	appController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(core.Map{
				"data": "ok",
			})
		})

		return ctrl
	}

	return core.NewModule(core.NewModuleOptions{
		Controllers: []core.Controller{appController},
	})
}

func Test_DefaultCors(t *testing.T) {
	app := core.CreateFactory(appModule)

	app.SetGlobalPrefix("/api")
	app.EnableCors(cors.Options{})

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	// Preflight
	req, err := http.NewRequest("OPTIONS", testServer.URL+"/api/test", nil)
	req.Header.Set("Access-Control-Request-Method", "GET")
	require.Nil(t, err)

	resp, err := testClient.Do(req)
	require.Nil(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))

	// Actual req
	resp, err = testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
}

func Test_Cors(t *testing.T) {
	app := core.CreateFactory(appModule)

	app.SetGlobalPrefix("/api")
	app.EnableCors(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders: []string{"*"},
		Credentials:    true,
	})

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	req, err := http.NewRequest("OPTIONS", testServer.URL+"/api/test", nil)
	req.Header.Set("Access-Control-Request-Method", "GET")
	require.Nil(t, err)

	// Preflight
	resp, err := testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// ACtual req
	resp, err = testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
	require.Equal(t, "true", resp.Header.Get("Access-Control-Allow-Credentials"))
}

func Test_FailedCors(t *testing.T) {
	app := core.CreateFactory(appModule)

	app.SetGlobalPrefix("/api")
	app.EnableCors(cors.Options{
		AllowedOrigins: []string{"localhost"},
		AllowedMethods: []string{"PUT", "PATCH", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "X-Custom-Header"},
	})

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	// Actual req
	req, err := http.NewRequest("GET", testServer.URL+"/api/test", nil)
	require.Nil(t, err)

	req.Header.Set("Origin", "localhost")

	resp, err := testClient.Do(req)
	require.Nil(t, err)

	require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Preflight
	req, err = http.NewRequest("OPTIONS", testServer.URL+"/api/test", nil)

	// Case 1: Wrong origin
	req.Header.Set("Access-Control-Request-Method", "GET")
	require.Nil(t, err)

	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Case 2: Wrong method
	req, err = http.NewRequest("OPTIONS", testServer.URL+"/api/test", nil)
	require.Nil(t, err)

	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Origin", "localhost")

	resp, err = testClient.Do(req)
	require.Nil(t, err)

	require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	// Case 3: Wrong headers
	req, err = http.NewRequest("OPTIONS", testServer.URL+"/api/test", nil)
	require.Nil(t, err)

	req.Header.Set("Access-Control-Request-Method", "PUT")
	req.Header.Set("Origin", "localhost")
	req.Header.Set("Access-Control-Request-Headers", "x-api-key")

	resp, err = testClient.Do(req)
	require.Nil(t, err)

	require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	// Case 4: Success
	req, err = http.NewRequest("OPTIONS", testServer.URL+"/api/test", nil)
	require.Nil(t, err)

	req.Header.Set("Access-Control-Request-Method", "PUT")
	req.Header.Set("Origin", "localhost")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")

	resp, err = testClient.Do(req)
	require.Nil(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_OriginFnc(t *testing.T) {
	app := core.CreateFactory(appModule)

	app.SetGlobalPrefix("/api")

	app.EnableCors(cors.Options{
		AllowedOriginCtx: func(r *http.Request) bool {
			return r.Referer() == "localhost"
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders: []string{"*"},
		PassThrough:    true,
	})

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	// Actual req
	req, err := http.NewRequest("GET", testServer.URL+"/api/test", nil)
	require.Nil(t, err)

	req.Header.Set("Referer", "localhost")

	resp, err := testClient.Do(req)
	require.Nil(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Bad referer
	req, err = http.NewRequest("GET", testServer.URL+"/api/test", nil)
	require.Nil(t, err)

	req.Header.Set("Referer", "google.com")

	resp, err = testClient.Do(req)
	require.Nil(t, err)

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

}
