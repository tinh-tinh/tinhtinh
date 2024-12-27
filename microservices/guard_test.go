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

func appGuard(addr string) microservices.Service {
	guard := func(ref core.RefProvider, ctx microservices.Ctx) bool {
		return ctx.Headers("key") == "value"
	}

	appService := func(module core.Module) core.Provider {
		handler := microservices.NewHandler(module, core.ProviderOptions{})

		handler.Guard(guard).OnResponse("guard", func(ctx microservices.Ctx) error {
			fmt.Printf("Receive data %v\n", ctx.Payload())
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

func clientGuard(addr string, event string) *core.App {
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			client := microservices.Inject(module)
			if client == nil {
				return ctx.Status(http.StatusInternalServerError).JSON(core.Map{"error": "client not found"})
			}
			// Example JSON messages to send
			messages := []Message{
				{"haha", 30},
			}

			for _, msg := range messages {
				client.SetHeaders("key", "value").Send(event, msg)
			}

			return ctx.JSON(core.Map{"data": "update"})
		})

		return ctrl
	}

	module := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{microservices.RegisterClient(tcp.NewClient(microservices.ConnectOptions{
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

func Test_Guard(t *testing.T) {
	app := appGuard("localhost:8082")

	go func() {
		app.Listen()
	}()

	// Allow some time for the server to start
	time.Sleep(100 * time.Millisecond)

	clientApp := clientGuard("localhost:8082", "guard")
	testServer := httptest.NewServer(clientApp.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	// Test event based
	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
