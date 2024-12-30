package redis_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	redis_store "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/microservices/redis"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
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

	guard := func(ref core.RefProvider, ctx microservices.Ctx) bool {
		return ctx.Headers("tenant") != "1"
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

			fmt.Printf("Order created: %v\n", orderService.orders)
			return nil
		})

		handler.Guard(guard).OnEvent("order.*", func(ctx microservices.Ctx) error {
			fmt.Printf("Order Updated: %v\n", ctx.Payload())
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

		client := microservices.Inject(module)
		ctrl.Post("", func(ctx core.Ctx) error {
			go client.Send("order.created", &Order{
				ID:   "order1",
				Name: "order1",
			})
			return ctx.JSON(core.Map{
				"data": []string{"product1", "product2"},
			})
		})

		ctrl.Post("multiple", func(ctx core.Ctx) error {
			go client.Publish("order.updated", &Order{
				ID:   "order1",
				Name: "order1",
			}, microservices.Header{"tenant": "1"})
			return ctx.JSON(core.Map{
				"data": []string{"product1", "product2"},
			})
		})

		return ctrl
	}

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{
				microservices.RegisterClient(redis.NewClient(redis.Options{
					Options: &redis_store.Options{
						Addr: addr,
					},
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

func DeliveryApp() microservices.Service {
	service := func(module core.Module) core.Provider {
		handler := microservices.NewHandler(module, core.ProviderOptions{})

		handler.OnEvent("order.*", func(ctx microservices.Ctx) error {
			fmt.Println("Delivery when have order:", ctx.Payload(&Order{}))
			return nil
		})

		return handler
	}

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{microservices.Register()},
			Providers: []core.Providers{
				service,
			},
		})
		return module
	}

	app := redis.New(appModule, redis.Options{
		Options: &redis_store.Options{
			Addr: "localhost:6379",
		},
	})

	return app
}

func Test_Practice(t *testing.T) {
	deliveryApp := DeliveryApp()
	go deliveryApp.Listen()

	orderApp := OrderApp()
	orderApp.ConnectMicroservice(redis.Open(redis.Options{
		Options: &redis_store.Options{
			Addr: "localhost:6379",
		},
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

	productApp := ProductApp("localhost:6379")
	testProductServer := httptest.NewServer(productApp.PrepareBeforeListen())
	defer testProductServer.Close()

	testClientProduct := testProductServer.Client()

	resp, err = testClientProduct.Post(testProductServer.URL+"/product-api/products", "application/json", nil)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	time.Sleep(1000 * time.Millisecond)

	resp, err = testClientOrder.Get(testOrderServer.URL + "/order-api/orders")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, `{"data":{"order1":true}}`, string(data))

	resp, err = testClientProduct.Post(testProductServer.URL+"/product-api/products/multiple", "application/json", nil)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	time.Sleep(1000 * time.Millisecond)
}

func Benchmark_Practice(b *testing.B) {
	orderApp := OrderApp()
	orderApp.ConnectMicroservice(redis.Open(redis.Options{
		Options: &redis_store.Options{
			Addr: "localhost:6379",
		},
	}))

	orderApp.StartAllMicroservices()
	testOrderServer := httptest.NewServer(orderApp.PrepareBeforeListen())
	defer testOrderServer.Close()

	time.Sleep(100 * time.Millisecond)

	productApp := ProductApp("localhost:6379")
	testProductServer := httptest.NewServer(productApp.PrepareBeforeListen())
	defer testProductServer.Close()

	testClientProduct := testProductServer.Client()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			resp, err := testClientProduct.Post(testProductServer.URL+"/product-api/products", "application/json", nil)
			require.Nil(b, err)
			require.Equal(b, http.StatusOK, resp.StatusCode)
		}
	})
}
