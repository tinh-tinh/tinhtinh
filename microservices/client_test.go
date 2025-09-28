package microservices_test

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/microservices"
	"github.com/tinh-tinh/tinhtinh/microservices/tcp"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func Test_Client(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:8000")
	require.Nil(t, err)

	go http.Serve(listener, nil)
	module := core.NewModule(core.NewModuleOptions{
		Imports: []core.Modules{
			microservices.RegisterClient(microservices.ClientOptions{
				Name: TCP_SERVICE,
				Transport: tcp.NewClient(tcp.Options{
					Addr: "localhost:8000",
				}),
			}),
		},
	})

	require.NotNil(t, microservices.InjectClient(module, TCP_SERVICE))

	module2 := core.NewModule(core.NewModuleOptions{})
	require.Nil(t, microservices.InjectClient(module2, TCP_SERVICE))
}

func Test_Client_Factory(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:8084")
	require.Nil(t, err)

	go http.Serve(listener, nil)
	module := core.NewModule(core.NewModuleOptions{
		Imports: []core.Modules{microservices.RegisterClientFactory(
			func(ref core.RefProvider) []microservices.ClientOptions {
				return []microservices.ClientOptions{
					{
						Name: TCP_SERVICE,
						Transport: tcp.NewClient(tcp.Options{
							Addr: "localhost:8084",
						}),
					},
				}
			},
		)},
	})

	require.NotNil(t, microservices.InjectClient(module, TCP_SERVICE))

	module2 := core.NewModule(core.NewModuleOptions{})
	require.Nil(t, microservices.InjectClient(module2, TCP_SERVICE))
}

func Test_HybridApp(t *testing.T) {
	appService := func(module core.Module) core.Provider {
		handler := microservices.NewHandler(module)

		handler.OnEvent("user.created", func(ctx microservices.Ctx) error {
			fmt.Println("User Created Data:", ctx.Payload())
			return nil
		})

		handler.OnEvent("user.updated", func(ctx microservices.Ctx) error {
			var message Message
			err := ctx.PayloadParser(&message)
			if err != nil {
				return err
			}
			fmt.Println("User Updated Data:", message)
			return nil
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

	app := core.CreateFactory(appModule)
	app.ConnectMicroservice(tcp.Open(tcp.Options{
		Addr: "localhost:3005",
	}))

	app.StartAllMicroservices()
	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	// Allow some time for the server to start
	time.Sleep(100 * time.Millisecond)

	clientApp := appClient("localhost:3005")
	testServer2 := httptest.NewServer(clientApp.PrepareBeforeListen())
	defer testServer2.Close()
	testClient := testServer2.Client()

	resp, err := testClient.Get(testServer2.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
