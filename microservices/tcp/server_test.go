package tcp_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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

func Test_EventBase(t *testing.T) {
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

	// Read output
	reader, writer, err := os.Pipe()
	require.Nil(t, err)

	originalStdout := os.Stdout
	defer func() {
		os.Stdout = originalStdout // Restore original stdout after test
		reader.Close()
		writer.Close()
	}()
	os.Stdout = writer

	// Test event based
	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	writer.Close() // Close writer to signal we're done writing
	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	require.Nil(t, err)

	// require.Equal(t, "Send message: {Alice 30} for event user.*\nReceived message: map[age:30 name:Alice] from event user.*\nuser.\nUser Created Data: {Alice 30}\nUser Updated Data: {Alice 30}\n", buf.String())
}

func Test_Response(t *testing.T) {
	app := appServer("localhost:4000")

	go func() {
		app.Listen()
	}()

	// Allow some time for the server to start
	time.Sleep(100 * time.Millisecond)

	clientApp := appClient("localhost:4000")
	testServer := httptest.NewServer(clientApp.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	// Read output
	reader, writer, err := os.Pipe()
	require.Nil(t, err)

	originalStdout := os.Stdout
	defer func() {
		os.Stdout = originalStdout // Restore original stdout after test
		reader.Close()
		writer.Close()
	}()
	os.Stdout = writer

	// Test event based
	resp, err := testClient.Get(testServer.URL + "/api/test/update")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	writer.Close() // Close writer to signal we're done writing
	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	require.Nil(t, err)

	// require.Equal(t, "Send message: {haha 30} for event user.updated\nReceived message: map[age:30 name:haha] from event user.updated\nUser Updated Data: {haha 30}\n", buf.String())
}

func Test_Broadcast(t *testing.T) {
	app := appServer("localhost:5000")

	go func() {
		app.Listen()
	}()

	// Allow some time for the server to start
	time.Sleep(100 * time.Millisecond)

	clientApp := appClient("localhost:5000")
	testServer := httptest.NewServer(clientApp.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	// Read output
	reader, writer, err := os.Pipe()
	require.Nil(t, err)

	originalStdout := os.Stdout
	defer func() {
		os.Stdout = originalStdout // Restore original stdout after test
		reader.Close()
		writer.Close()
	}()
	os.Stdout = writer

	// Test event based
	resp, err := testClient.Get(testServer.URL + "/api/test/broadcast")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	writer.Close() // Close writer to signal we're done writing
	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	require.Nil(t, err)

	// require.Equal(t, "Send message: {Broadcast 1 30} for event *\nReceived message: map[age:30 name:Broadcast 1] from event *\nUser Created Data: {Broadcast 1 30}\nUser Updated Data: {Broadcast 1 30}\n", buf.String())
}

func appServer(addr string) microservices.Service {
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
				fmt.Println("User Created Data:", decodedData)
			}

			return nil
		})

		handler.OnEvent("user.updated", func(param ...interface{}) interface{} {
			if len(param) == 0 {
				return nil
			}
			msg := param[0]
			var decodedData Message
			if msg != nil {
				dataBytes, _ := json.Marshal(msg)
				_ = json.Unmarshal(dataBytes, &decodedData)
				fmt.Println("User Updated Data:", decodedData)
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
		Addr: addr,
	})

	return app
}

func appClient(addr string) *core.App {
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
			}

			for _, msg := range messages {
				client.Send("user.*", msg)
			}

			return ctx.JSON(core.Map{"data": "ok"})
		})

		ctrl.Get("broadcast", func(ctx core.Ctx) error {
			client := microservices.Inject(module)
			if client == nil {
				return ctx.Status(http.StatusInternalServerError).JSON(core.Map{"error": "client not found"})
			}
			// Example JSON messages to send
			messages := []Message{
				{"Broadcast 1", 30},
			}

			for _, msg := range messages {
				client.Broadcast(msg)
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
