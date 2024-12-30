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

type Message struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func Test_RPC(t *testing.T) {
	app := appServer("localhost:8080")

	go func() {
		app.Listen()
	}()

	// Allow some time for the server to start
	time.Sleep(100 * time.Millisecond)

	clientApp := appClient("localhost:8080")
	testServer := httptest.NewServer(clientApp.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	// Test event based
	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_Event(t *testing.T) {
	app := appServer("localhost:4000")

	go func() {
		app.Listen()
	}()

	// Allow some time for the server to start
	time.Sleep(1000 * time.Millisecond)

	clientApp := appClient("localhost:4000")
	testServer := httptest.NewServer(clientApp.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	// Test event based
	resp, err := testClient.Get(testServer.URL + "/api/test/event")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func appServer(addr string) microservices.Service {
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
	app := tcp.New(appModule, microservices.Options{
		Addr: addr,
	})

	return app
}

func appClient(addr string) *core.App {
	clientController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("event", func(ctx core.Ctx) error {
			client := microservices.Inject(module)
			if client == nil {
				return ctx.Status(http.StatusInternalServerError).JSON(core.Map{"error": "client not found"})
			}
			// Example JSON messages to send
			messages := []Message{
				{"haha", 30},
			}

			for _, msg := range messages {
				client.Publish("user.updated", msg)
			}

			return ctx.JSON(core.Map{"data": "update"})
		})

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
			Imports: []core.Modules{microservices.RegisterClient(tcp.NewClient(tcp.Options{
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
