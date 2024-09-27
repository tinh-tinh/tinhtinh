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

	app := CreateFactory(appModule, "api")

	testServer := httptest.NewServer(app.prepareBeforeListen())
	defer testServer.Close()
}
