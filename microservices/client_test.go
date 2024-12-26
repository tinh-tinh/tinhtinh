package microservices_test

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
	"github.com/tinh-tinh/tinhtinh/v2/microservices/tcp"
)

func Test_Client(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:8000")
	require.Nil(t, err)

	go http.Serve(listener, nil)
	module := core.NewModule(core.NewModuleOptions{
		Imports: []core.Modules{microservices.RegisterClient(tcp.NewClient(microservices.ConnectOptions{
			Addr: "localhost:8000",
		}))},
	})

	require.NotNil(t, microservices.Inject(module))

	module2 := core.NewModule(core.NewModuleOptions{})
	require.Nil(t, microservices.Inject(module2))
}

func Test_HybridApp(t *testing.T) {
	appService := func(module core.Module) core.Provider {
		handler := microservices.NewHandler(module, core.ProviderOptions{})

		handler.OnResponse("user.created", func(ctx microservices.Ctx) error {
			fmt.Println("User Created Data:", ctx.Payload(&Message{}))
			return nil
		})

		handler.OnEvent("user.updated", func(ctx microservices.Ctx) error {
			fmt.Println("User Updated Data:", ctx.Payload(&Message{}))
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
	app.ConnectMicroservice(tcp.Open(microservices.ConnectOptions{
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
