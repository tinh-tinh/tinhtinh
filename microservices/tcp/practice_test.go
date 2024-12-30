package tcp_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
	"github.com/tinh-tinh/tinhtinh/v2/microservices/tcp"
)

type Order struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func OrderApp() *core.App {
	type OrderService struct {
		orders map[string]interface{}
		mutex  sync.RWMutex
	}

	const ORDER core.Provide = "orders"
	service := func(module core.Module) core.Provider {
		prd := module.NewProvider(core.ProviderOptions{
			Name: ORDER,
			Value: &OrderService{
				orders: make(map[string]interface{}),
				mutex:  sync.RWMutex{},
			},
		})

		return prd
	}

	handlerService := func(module core.Module) core.Provider {
		handler := microservices.NewHandler(module, core.ProviderOptions{})

		orderService := module.Ref(ORDER).(*OrderService)
		handler.OnResponse("order.created", func(ctx microservices.Ctx) error {
			data := ctx.Payload(&Order{}).(*Order)

			orderService.mutex.Lock()
			if orderService.orders[data.ID] == nil {
				orderService.orders[data.ID] = true
			}

			orderService.mutex.Unlock()
			return nil
		})

		return handler
	}

	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("orders")

		ctrl.Get("", func(ctx core.Ctx) error {
			orderService := module.Ref(ORDER).(*OrderService)
			return ctx.JSON(core.Map{
				"data": orderService.orders,
			})
		})

		return ctrl
	}

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports:     []core.Modules{microservices.Register()},
			Controllers: []core.Controllers{controller},
			Providers: []core.Providers{
				service,
				handlerService,
			},
		})
		return module
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("order-api")

	return app
}

func ProductApp(addr string) *core.App {
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("products")

		ctrl.Post("", func(ctx core.Ctx) error {
			client := microservices.Inject(module)

			client.Send("order.created", &Order{
				ID:   "order1",
				Name: "order1",
			})
			return ctx.JSON(core.Map{
				"data": []string{"product1", "product2"},
			})
		})

		return ctrl
	}

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{
				microservices.RegisterClient(tcp.NewClient(tcp.Options{
					Addr: addr,
				})),
			},
			Controllers: []core.Controllers{controller},
		})
		return module
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("product-api")

	return app
}

func Test_Practice(t *testing.T) {
	orderApp := OrderApp()
	orderApp.ConnectMicroservice(tcp.Open(microservices.Options{
		Addr: "localhost:3006",
	}))

	orderApp.StartAllMicroservices()
	testOrderServer := httptest.NewServer(orderApp.PrepareBeforeListen())
	defer testOrderServer.Close()

	testClientOrder := testOrderServer.Client()

	resp, err := testClientOrder.Get(testOrderServer.URL + "/order-api/orders")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":{}}`, string(data))

	productApp := ProductApp("localhost:3006")
	testProductServer := httptest.NewServer(productApp.PrepareBeforeListen())
	defer testProductServer.Close()

	testClientProduct := testProductServer.Client()

	resp, err = testClientProduct.Post(testProductServer.URL+"/product-api/products", "application/json", nil)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	time.Sleep(100 * time.Millisecond)

	resp, err = testClientOrder.Get(testOrderServer.URL + "/order-api/orders")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":{"order1":true}}`, string(data))
}
