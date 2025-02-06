package exception_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/common/exception"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
	"github.com/tinh-tinh/tinhtinh/v2/microservices/tcp"
)

func appServer(add string) microservices.Service {
	appService := func(module core.Module) core.Provider {
		handler := microservices.NewHandler(module, core.ProviderOptions{})

		handler.OnResponse("exception", func(ctx microservices.Ctx) error {
			panic(exception.ThrowRpc("error"))
		})

		return handler
	}

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{microservices.Register()},
			Providers: []core.Providers{
				appService,
			},
		})

		return module
	}

	app := tcp.New(appModule, tcp.Options{
		Addr: add,
	})

	return app
}

func appClient(addr string) *core.App {
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			client := microservices.Inject(module)
			go client.Send("exception", map[string]interface{}{"data": "ok"})
			return ctx.JSON(core.Map{})
		})

		return ctrl
	}

	module := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{microservices.RegisterClient(tcp.NewClient(tcp.Options{
				Addr: addr,
			}))},
			Controllers: []core.Controllers{
				controller,
			},
		})
		return module
	}
	app := core.CreateFactory(module)
	app.SetGlobalPrefix("api")

	return app
}

func Test_RCP_Exception(t *testing.T) {
	app := appServer("localhost:8083")

	go func() {
		app.Listen()
	}()

	// Allow some time for the server to start
	time.Sleep(100 * time.Millisecond)

	clientApp := appClient("localhost:8083")
	testServer := httptest.NewServer(clientApp.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	// Test event based
	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	time.Sleep(100 * time.Millisecond)
}
