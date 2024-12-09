package core_test

import (
	"net/http/httptest"
	"testing"

	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func Test_OnInit(t *testing.T) {
	appModule := func() *core.DynamicModule {
		module := core.NewModule(core.NewModuleOptions{})

		module.OnInit(func(module *core.DynamicModule) {
			t.Log("OnInit")
		})
		return module
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
}

func Test_BeforeShutdown(t *testing.T) {
	appModule := func() *core.DynamicModule {
		module := core.NewModule(core.NewModuleOptions{})

		return module
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")
	app.BeforeShutdown(func() {
		t.Log("BeforeShutdown")
	})

	app.AfterShutdown(func() {
		t.Log("AfterShutdown")
	})

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	t.Log(testServer.URL)
}
