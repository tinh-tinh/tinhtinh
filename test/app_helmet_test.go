package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/middleware/helmet"
)

func Test_HelmetFullSetting(t *testing.T) {
	appController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) {
			ctx.JSON(core.Map{
				"data": "ok",
			})
		})

		return ctrl
	}

	appModule := func() *core.DynamicModule {
		return core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controller{appController},
		})
	}

	app := core.CreateFactory(appModule)

	app.SetGlobalPrefix("/api")
	app.Use(helmet.Handler(helmet.Options{
		XPoweredBy:                helmet.XPoweredBy{Enabled: true, Value: "OK"},
		CrossOriginEmbedderPolicy: helmet.CrossOriginEmbedderPolicy{Enabled: true, Policy: "require-corp"},
		CrossOriginOpenerPolicy:   helmet.CrossOriginOpenerPolicy{Enabled: true, Policy: "same-origin-allow-popups"},
		CrossOriginResourcePolicy: helmet.CrossOriginResourcePolicy{Enabled: true, Policy: "same-site"},
		XXSSProtection:            helmet.XXSSProtection{Enabled: true, Value: "1; mode=block"},
		XContentTypeOptions:       helmet.XContentTypeOptions{Enabled: true, Value: "nosniff"},
		XDownloadOptions:          helmet.XDownloadOptions{Enabled: true, Value: "noopen"},
		XFrameOptions:             helmet.XFrameOptions{Enabled: true, Action: "deny"},
		XPermittedCrossDomainPolicies: helmet.XPermittedCrossDomainPolicies{
			Enabled:           true,
			PermittedPolicies: "none",
		},
		ReferrerPolicy: helmet.ReferrerPolicy{Enabled: true, Policy: "no-referrer"},
		StrictTransportSecurity: helmet.StrictTransportSecurity{
			Enabled:           true,
			MaxAge:            3600,
			IncludeSubDomains: true,
			Preload:           true,
		},
		ContentSecurityPolicy: helmet.ContentSecurityPolicy{
			OptionSecurityPolicy: &helmet.OptionSecurityPolicy{
				DefaultSrc:              []string{"'self'"},
				ScriptSrc:               []string{"'self'"},
				StyleSrc:                []string{"'self'"},
				ImgSrc:                  []string{"'self'"},
				FontSrc:                 []string{"'self'"},
				FrameSrc:                []string{"'self'"},
				ManifestSrc:             []string{"'self'"},
				ObjectSrc:               []string{"'self'"},
				UpgradeInsecureRequests: true,
				ReportOnly:              true,
			},
		},
		XDnsPrefetchControl: helmet.XDnsPrefetchControl{Enabled: true, Value: "1"},
	}))

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "OK", resp.Header.Get("X-Powered-By"))
	require.Equal(t, "require-corp", resp.Header.Get("Cross-Origin-Embedder-Policy"))
	require.NotEqual(t, "credentialless", resp.Header.Get("Cross-Origin-Embedder-Policy"))
	require.Equal(t, "same-origin-allow-popups", resp.Header.Get("Cross-Origin-Opener-Policy"))
	require.Equal(t, "same-site", resp.Header.Get("Cross-Origin-Resource-Policy"))
	require.Equal(t, "1; mode=block", resp.Header.Get("X-XSS-Protection"))
	require.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
	require.Equal(t, "noopen", resp.Header.Get("X-Download-Options"))
	require.Equal(t, "deny", resp.Header.Get("X-Frame-Options"))
	require.Equal(t, "none", resp.Header.Get("X-Permitted-Cross-Domain-Policies"))
	require.Equal(t, "no-referrer", resp.Header.Get("Referrer-Policy"))
	require.Equal(t, "max-age=3600; includeSubDomains; preload", resp.Header.Get("Strict-Transport-Security"))
	require.Equal(t, "1", resp.Header.Get("X-DNS-Prefetch-Control"))
}

func Test_DefaultHelmet(t *testing.T) {
	appController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) {
			ctx.JSON(core.Map{
				"data": "ok",
			})
		})

		return ctrl
	}

	appModule := func() *core.DynamicModule {
		return core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controller{appController},
		})
	}

	app := core.CreateFactory(appModule)

	app.SetGlobalPrefix("/api")
	app.Use(helmet.Handler(helmet.Options{}))

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	require.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
	require.Equal(t, "off", resp.Header.Get("X-DNS-Prefetch-Control"))
	require.Equal(t, "SAMEORIGIN", resp.Header.Get("X-Frame-Options"))
	require.Equal(t, "noopen", resp.Header.Get("X-Download-Options"))
	require.Equal(t, "0", resp.Header.Get("X-XSS-Protection"))
}

func Test_HelmetSomeSetting(t *testing.T) {
	appController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) {
			ctx.JSON(core.Map{
				"data": "ok",
			})
		})

		return ctrl
	}

	appModule := func() *core.DynamicModule {
		return core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controller{appController},
		})
	}

	app := core.CreateFactory(appModule)

	app.SetGlobalPrefix("/api")

	app.Use(helmet.Handler(helmet.Options{
		CrossOriginEmbedderPolicy: helmet.CrossOriginEmbedderPolicy{Policy: "credentialless", Enabled: true},
		ReferrerPolicy:            helmet.ReferrerPolicy{Policy: []string{"no-referrer-when-downgrade", "origin"}, Enabled: true},
		XDnsPrefetchControl:       helmet.XDnsPrefetchControl{Enabled: false},
	}))

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	require.Equal(t, "credentialless", resp.Header.Get("Cross-Origin-Embedder-Policy"))
	require.Equal(t, "no-referrer-when-downgrade,origin", resp.Header.Get("Referrer-Policy"))
	require.Equal(t, "off", resp.Header.Get("X-DNS-Prefetch-Control"))
}
