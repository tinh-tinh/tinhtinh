package tcp_test

import (
	"encoding/json"
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

func Test_Server(t *testing.T) {
	type Message struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	appService := func(module core.Module) core.Provider {
		handler := microservices.NewHandler(module, core.ProviderOptions{})

		handler.OnResponse("user.created", func(param ...interface{}) interface{} {
			if len(param) == 0 {
				return nil
			}
			msg := param[0]
			var decodedData Message
			if msg != nil {
				dataBytes, _ := json.Marshal(msg)
				_ = json.Unmarshal(dataBytes, &decodedData)
				fmt.Println("Decoded Data:", decodedData)
			}

			return nil
		})

		handler.OnResponse("user.updated", func(param ...interface{}) interface{} {
			if len(param) == 0 {
				return nil
			}
			msg := param[0]
			var decodedData Message
			if msg != nil {
				dataBytes, _ := json.Marshal(msg)
				_ = json.Unmarshal(dataBytes, &decodedData)
				fmt.Println("Decoded Data:", decodedData)
			}

			return nil
		})

		return handler
	}

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Providers: []core.Providers{
				appService,
			},
		})
		return module
	}
	app := tcp.New(appModule, microservices.ConnectOptions{
		Addr: "localhost:8080",
	})
	go func() {
		app.Listen()
	}()

	// Allow some time for the server to start
	time.Sleep(100 * time.Millisecond)

	clientController := func(module core.Module) core.Controller {
		ctrl := module.NewController("test")

		ctrl.Get("update", func(ctx core.Ctx) error {
			client := microservices.Inject(module)
			if client == nil {
				return ctx.Status(http.StatusInternalServerError).JSON(core.Map{"error": "client not found"})
			}
			// Example JSON messages to send
			messages := []Message{
				{"haha", 30},
				{"hihi", 25},
				{"huhu", 35},
			}

			for _, msg := range messages {
				client.Send("user.updated", msg)
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
				{"Bob", 25},
				{"Charlie", 35},
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
				Addr: "localhost:8080",
			}))},
			Controllers: []core.Controllers{
				clientController,
			},
		})
		return module
	}
	clientApp := core.CreateFactory(clientModule)
	clientApp.SetGlobalPrefix("api")

	testServer := httptest.NewServer(clientApp.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/update")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
