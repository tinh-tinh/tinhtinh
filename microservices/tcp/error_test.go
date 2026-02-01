package tcp_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/microservices"
	"github.com/tinh-tinh/tinhtinh/microservices/tcp"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func Test_Client_Error(t *testing.T) {
	appProvider := func(module core.Module) core.Provider {
		handler := microservices.NewHandler(module)

		handler.OnEvent("abc", func(ctx microservices.Ctx) error {
			return nil
		})

		return handler
	}
	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{
				microservices.RegisterClient(microservices.ClientOptions{
					Name: microservices.TCP,
					Transport: tcp.NewClient(tcp.Options{
						Addr: "localhost:9091",
					}),
				}),
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
			Imports:   []core.Modules{microservices.Register()},
			Providers: []core.Providers{appProvider},
		})
		return module
	}

	server := tcp.NewServer(tcp.Options{
		Addr: "localhost:9000",
	})
	server.Create(serverModule())
	go server.Listen()

	time.Sleep(100 * time.Millisecond)

	clientController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		client := microservices.InjectClient(module, microservices.TCP)
		ctrl.Get("", func(ctx core.Ctx) error {
			go client.Timeout(1*time.Microsecond).Publish("abc", 1000)
			return ctx.JSON(core.Map{"data": "ok"})
		})

		return ctrl
	}

	clientModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{
				microservices.RegisterClient(microservices.ClientOptions{
					Name: microservices.TCP,
					Transport: tcp.NewClient(tcp.Options{
						Addr: "localhost:9000",
						Config: microservices.Config{
							Serializer: func(v interface{}) ([]byte, error) {
								return nil, errors.New("error")
							},
						},
					}),
				}),
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
		server := tcp.NewServer(tcp.Options{
			Addr: "localhost",
		})
		server.Create(serverModule())
		server.Listen()
	})

	require.Panics(t, func() {
		serverModule := func() core.Module {
			module := core.NewModule(core.NewModuleOptions{})
			return module
		}
		server := tcp.NewServer(tcp.Options{
			Addr: "localhost:9090",
		})
		server.Create(serverModule())
		server.Listen()
	})
}
