package microservices_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
	"github.com/tinh-tinh/tinhtinh/v2/microservices/tcp"
)

type User struct {
	Name string `json:"name" validate:"isName"`
	Age  int    `json:"age" validate:"isInt"`
}

func Test_Pipe(t *testing.T) {
	app := appPipe("localhost:8085")

	go func() {
		app.Listen()
	}()

	// Allow some time for the server to start
	time.Sleep(100 * time.Millisecond)

	clientApp := clientPipe("localhost:8085")
	testServer := httptest.NewServer(clientApp.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	// Test event based
	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func appPipe(addr string) microservices.Service {
	appService := func(module core.Module) core.Provider {
		handler := microservices.NewHandler(module, core.ProviderOptions{})

		handler.Pipe(&User{}).OnResponse("user.created", func(ctx microservices.Ctx) error {
			fmt.Println("User Created Data:", ctx.Get(microservices.PIPE))
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
	app := tcp.New(appModule, microservices.ConnectOptions{
		Addr: addr,
	})

	return app
}

func clientPipe(addr string) *core.App {
	clientController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			client := microservices.Inject(module)
			if client == nil {
				return ctx.Status(http.StatusInternalServerError).JSON(core.Map{"error": "client not found"})
			}
			// Example JSON messages to send
			messages := []Message{
				{"Alice", 30},
			}

			for _, msg := range messages {
				client.Send("user.created", msg)
			}

			return ctx.JSON(core.Map{"data": "ok"})
		})

		return ctrl
	}

	clientModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{microservices.RegisterClient(tcp.NewClient(microservices.ConnectOptions{
				Addr: addr,
			}))},
			Controllers: []core.Controllers{
				clientController,
			},
		})
		return module
	}
	clientApp := core.CreateFactory(clientModule)
	clientApp.SetGlobalPrefix("api")

	return clientApp
}
