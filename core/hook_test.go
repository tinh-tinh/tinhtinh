package core

import (
	"net/http/httptest"
	"testing"
)

func Test_OnInit(t *testing.T) {
	appModule := func() *DynamicModule {
		module := NewModule(NewModuleOptions{})

		module.OnInit(func(module *DynamicModule) {
			t.Log("OnInit")
		})
		return module
	}

	app := CreateFactory(appModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.prepareBeforeListen())
	defer testServer.Close()
}

func Test_BeforeShutdown(t *testing.T) {
	appModule := func() *DynamicModule {
		module := NewModule(NewModuleOptions{})

		return module
	}

	app := CreateFactory(appModule)
	app.SetGlobalPrefix("/api")
	app.BeforeShutdown(func() {
		t.Log("BeforeShutdown")
	})

	testServer := httptest.NewServer(app.prepareBeforeListen())
	defer testServer.Close()

	t.Log(testServer.URL)
}
