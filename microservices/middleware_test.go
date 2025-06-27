package microservices_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/microservices"
	"github.com/tinh-tinh/tinhtinh/microservices/tcp"
	"github.com/tinh-tinh/tinhtinh/v2/common/exception"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func appMiddleware(addr string) microservices.Service {
	middleware := func(ctx microservices.Ctx) error {
		if ctx.Headers("key") != "value" {
			return exception.ThrowRpc("error")
		}
		ctx.Set("key", "value")
		return ctx.Next()
	}

	appService := func(module core.Module) core.Provider {
		handler := microservices.NewHandler(module, core.ProviderOptions{})

		handler.Use(middleware).OnResponse("middleware", func(ctx microservices.Ctx) error {
			fmt.Printf("Receive data %v with key is %v\n", ctx.Payload(), ctx.Get("key"))
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
	app := tcp.New(appModule, tcp.Options{
		Addr: addr,
	})

	return app
}

func clientMiddleware(addr string, event string) *core.App {
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			client := microservices.InjectClient(module, TCP_SERVICE)
			if client == nil {
				return ctx.Status(http.StatusInternalServerError).JSON(core.Map{"error": "client not found"})
			}
			// Example JSON messages to send
			messages := []Message{
				{"haha", 30},
				{"haha", 30},
				{"haha", 30},
				{"haha", 30},
				{"haha", 30},
			}

			for i, msg := range messages {
				if i%2 == 0 {
					go client.Send(event, msg, microservices.Header{"key": "value"})
				} else {
					client.Send(event, msg)
				}
			}

			return ctx.JSON(core.Map{"data": "update"})
		})

		return ctrl
	}

	module := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{
				microservices.RegisterClient(microservices.ClientOptions{
					Name: TCP_SERVICE,
					Transport: tcp.NewClient(tcp.Options{
						Addr: addr,
					}),
				}),
			},
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

func Test_Middleware(t *testing.T) {
	app := appMiddleware("localhost:8081")

	go func() {
		app.Listen()
	}()

	// Allow some time for the server to start
	time.Sleep(100 * time.Millisecond)

	clientApp := clientMiddleware("localhost:8081", "middleware")
	testServer := httptest.NewServer(clientApp.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	// Test event based
	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
