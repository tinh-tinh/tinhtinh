package tcp_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
	"github.com/tinh-tinh/tinhtinh/v2/microservices/tcp"
)

func Test_Client_Error(t *testing.T) {
	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{
				microservices.RegisterClient(tcp.NewClient(tcp.Options{
					Addr: "localhost:9091",
				})),
			},
		})

		return module
	}
	require.Panics(t, func() {
		app := core.CreateFactory(appModule)
		app.SetGlobalPrefix("api")
	})

	serverModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{microservices.Register()},
		})
		return module
	}

	server := tcp.New(serverModule, tcp.Options{
		Addr: "localhost:9000",
	})
	go server.Listen()

	time.Sleep(100 * time.Millisecond)

	clientController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		client := microservices.Inject(module)
		ctrl.Get("", func(ctx core.Ctx) error {
			go client.Send("abc", 1000)
			return ctx.JSON(core.Map{"data": "ok"})
		})

		return ctrl
	}

	clientModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{
				microservices.RegisterClient(tcp.NewClient(tcp.Options{
					Addr: "localhost:9000",
					Config: microservices.Config{
						Serializer: func(v interface{}) ([]byte, error) {
							return nil, errors.New("error")
						},
					},
				})),
			},
			Controllers: []core.Controllers{clientController},
		})

		return module
	}

	app := core.CreateFactory(clientModule)
	app.SetGlobalPrefix("api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()
	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_Server_Error(t *testing.T) {
	require.Panics(t, func() {
		serverModule := func() core.Module {
			module := core.NewModule(core.NewModuleOptions{
				Imports: []core.Modules{microservices.Register()},
			})
			return module
		}
		server := tcp.New(serverModule, tcp.Options{
			Addr: "localhost",
		})
		server.Listen()
	})

	require.Panics(t, func() {
		serverModule := func() core.Module {
			module := core.NewModule(core.NewModuleOptions{})
			return module
		}
		server := tcp.New(serverModule, tcp.Options{
			Addr: "localhost:9090",
		})
		server.Listen()
	})
}
