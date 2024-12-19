package redis_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	redis_store "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/microservices/redis"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

type Message struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func Test_Client(t *testing.T) {
	module := core.NewModule(core.NewModuleOptions{
		Imports: []core.Modules{microservices.RegisterClient(redis.NewClient(microservices.ConnectOptions{
			Addr: "localhost:6379",
		}))},
	})

	require.NotNil(t, microservices.Inject(module))

	module2 := core.NewModule(core.NewModuleOptions{})
	require.Nil(t, microservices.Inject(module2))
}

func Test_HybridApp(t *testing.T) {
	appService := func(module core.Module) core.Provider {
		handler := microservices.NewHandler(module, core.ProviderOptions{})

		handler.OnResponse("user.created", func(param ...interface{}) interface{} {
			if len(param) == 0 {
				return nil
			}
			msg := param[0]
			var decodedData redis_store.Message
			if msg != nil {
				dataBytes, _ := json.Marshal(msg)
				_ = json.Unmarshal(dataBytes, &decodedData)
				fmt.Println("User Created Data:", decodedData)
			}

			return nil
		})

		handler.OnEvent("user.*", func(param ...interface{}) interface{} {
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

	app := core.CreateFactory(appModule)
	app.ConnectMicroservice(redis.Open(microservices.ConnectOptions{
		Addr: "localhost:6379",
	}))

	app.StartAllMicroservices()
	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	// Allow some time for the server to start
	time.Sleep(100 * time.Millisecond)

	clientApp := appClient()
	testServer2 := httptest.NewServer(clientApp.PrepareBeforeListen())
	defer testServer2.Close()
	testClient := testServer2.Client()

	// resp, err := testClient.Get(testServer.URL + "/api/test")
	// require.Nil(t, err)
	// require.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err := testClient.Get(testServer2.URL + "/api/test/update")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
