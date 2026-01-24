package tcp_test

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/microservices"
	"github.com/tinh-tinh/tinhtinh/microservices/tcp"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

// Args holds the arguments for the remote Add call.
type Args struct {
	A, B int
}

// Calculator defines the service object that contains the RPC methods.
// The type name ("Calculator") will be part of the method call signature (e.g., "Calculator.Add").
func CalculateApp() *core.App {
	handlers := func(module core.Module) core.Provider {
		handler := microservices.NewHandler(module, microservices.TCP)

		handler.OnReply("add", func(ctx microservices.Ctx) (reply []byte, err error) {
			args := Args{}
			err = ctx.PayloadParser(&args)
			if err != nil {
				return nil, errors.New("invalid payload")
			}
			results := args.A + args.B
			log.Printf("Server handled Add: %d + %d = %d\n", args.A, args.B, results)
			return ctx.Reply(results)
		})

		return handler
	}

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports:   []core.Modules{microservices.Register(microservices.TCP)},
			Providers: []core.Providers{handlers},
		})
		return module
	}

	app := core.CreateFactory(appModule)

	return app
}

func CallApp(addr string) *core.App {
	controller := func(module core.Module) core.Controller {
		ctrl := module.NewController("call")
		client := microservices.InjectClient(module, microservices.TCP)

		ctrl.Post("register", func(ctx core.Ctx) error {
			args := Args{A: 15, B: 27}
			var reply int

			err := client.Send("add", args, &reply)
			if err != nil {
				return err
			}
			return ctx.JSON(reply)
		})

		return ctrl
	}

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{
				microservices.RegisterClient(microservices.ClientOptions{
					Name: microservices.TCP,
					Transport: tcp.NewClient(tcp.Options{
						Addr: addr,
					}),
				}),
			},
			Controllers: []core.Controllers{controller},
		})

		return module
	}

	app := core.CreateFactory(appModule)

	return app
}

func TestRpc(t *testing.T) {
	calculateApp := CalculateApp()
	calculateApp.ConnectMicroservice(tcp.Open(tcp.Options{
		Addr: "localhost:5155",
	}))
	calculateApp.StartAllMicroservices()
	testServerCalculate := httptest.NewServer(calculateApp.PrepareBeforeListen())
	defer testServerCalculate.Close()

	time.Sleep(100 * time.Millisecond)

	callApp := CallApp("localhost:5155")
	testServerCall := httptest.NewServer(callApp.PrepareBeforeListen())
	defer testServerCall.Close()

	testClientCall := testServerCall.Client()

	resp, err := testClientCall.Post(testServerCall.URL+"/call/register", "application/json", nil)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, "42", string(data))
}

func MiddlewareApp() *core.App {
	handlers := func(module core.Module) core.Provider {
		handler := microservices.NewHandler(module, microservices.TCP)

		handler.Use(func(ctx microservices.Ctx) error {
			if ctx.Headers("auth") != "secret" {
				return errors.New("unauthorized")
			}
			return nil
		})

		handler.OnReply("check", func(ctx microservices.Ctx) (reply []byte, err error) {
			return ctx.Reply("authorized")
		})

		return handler
	}

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports:   []core.Modules{microservices.Register(microservices.TCP)},
			Providers: []core.Providers{handlers},
		})
		return module
	}

	return core.CreateFactory(appModule)
}

func PanicApp() *core.App {
	handlers := func(module core.Module) core.Provider {
		handler := microservices.NewHandler(module, microservices.TCP)

		handler.OnReply("panic", func(ctx microservices.Ctx) (reply []byte, err error) {
			panic("something went wrong")
		})

		return handler
	}

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports:   []core.Modules{microservices.Register(microservices.TCP)},
			Providers: []core.Providers{handlers},
		})
		return module
	}

	return core.CreateFactory(appModule)
}

func TestRpcMiddleware(t *testing.T) {
	app := MiddlewareApp()
	app.ConnectMicroservice(tcp.Open(tcp.Options{
		Addr: "localhost:5156",
	}))
	app.StartAllMicroservices()
	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	time.Sleep(100 * time.Millisecond)

	client := tcp.NewClient(tcp.Options{
		Addr: "localhost:5156",
	})

	var reply string
	// Test unauthorized
	err := client.Send("check", nil, &reply)
	require.Error(t, err)
	require.Equal(t, "unauthorized", err.Error())

	// Test authorized
	err = client.Send("check", nil, &reply, microservices.Header{"auth": "secret"})
	require.Nil(t, err)
	require.Equal(t, "authorized", reply)
}

func TestRpcPanic(t *testing.T) {
	app := PanicApp()
	app.ConnectMicroservice(tcp.Open(tcp.Options{
		Addr: "localhost:5157",
	}))
	app.StartAllMicroservices()
	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	time.Sleep(100 * time.Millisecond)

	client := tcp.NewClient(tcp.Options{
		Addr: "localhost:5157",
	})

	var reply string
	err := client.Send("panic", nil, &reply)
	require.Error(t, err)
	require.Equal(t, "something went wrong", err.Error())
}
