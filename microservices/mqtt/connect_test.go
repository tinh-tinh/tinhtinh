package mqtt_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	mqtt_store "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/microservices/mqtt"
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
			var data *Order
			err := ctx.PayloadParser(&data)
			if err != nil {
				return err
			}

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

		client := microservices.InjectClient(module, microservices.MQTT)
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
		opts := mqtt_store.NewClientOptions().AddBroker(addr).SetClientID("product-app")
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{
				microservices.RegisterClient(microservices.ClientOptions{
					Name: microservices.MQTT,
					Transport: mqtt.NewClient(mqtt.Options{
						ClientOptions: opts,
					}),
				}),
			},
			Controllers: []core.Controllers{controller},
		})
		return module
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("product-api")

	return app
}

func DeliveryApp(addr string) microservices.Service {
	service := func(module core.Module) core.Provider {
		handler := microservices.NewHandler(module, core.ProviderOptions{})

		handler.OnEvent("order.*", func(ctx microservices.Ctx) error {
			var data *Order
			err := ctx.PayloadParser(&data)
			if err != nil {
				return err
			}

			fmt.Println("Delivery when have order:", data)
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

	opts := mqtt_store.NewClientOptions().AddBroker(addr).SetClientID("delivery-app")
	app := mqtt.New(appModule, mqtt.Options{
		ClientOptions: opts,
	})

	return app
}

func Test_Practice(t *testing.T) {
	deliveryApp := DeliveryApp("mqtt://localhost:1883")
	go deliveryApp.Listen()

	orderApp := OrderApp()
	opts := mqtt_store.NewClientOptions().AddBroker("mqtt://localhost:1883").SetClientID("order-app")
	orderApp.ConnectMicroservice(mqtt.Open(mqtt.Options{
		ClientOptions: opts,
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

	productApp := ProductApp("mqtt://localhost:1883")
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
