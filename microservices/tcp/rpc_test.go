package tcp_test

import (
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
type Calculator struct{}

// Add is the RPC method. It must follow the signature:
// func (t *T) MethodName(argType T1, replyType *T2) error
func (t *Calculator) Add(args *Args, reply *int) error {
	*reply = args.A + args.B
	log.Printf("Server handled Add: %d + %d = %d\n", args.A, args.B, *reply)
	return nil
}

func CalculateApp() *core.App {
	handlers := func(module core.Module) core.Provider {
		handler := microservices.NewHandler(module, microservices.TCP)

		handler.RegisterRPC(new(Calculator))

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

			args := &Args{A: 15, B: 27}
			var reply int

			err := client.Send("Calculator.Add", args, &reply)
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
		Addr: "localhost:4379",
	}))
	calculateApp.StartAllMicroservices()
	testServerCalculate := httptest.NewServer(calculateApp.PrepareBeforeListen())
	defer testServerCalculate.Close()

	time.Sleep(100 * time.Millisecond)

	callApp := CallApp("localhost:4379")
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
