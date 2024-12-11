package tcp_test

import (
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
	"github.com/tinh-tinh/tinhtinh/v2/microservices/tcp"
)

func Test_Client(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:8000")
	require.Nil(t, err)

	go http.Serve(listener, nil)
	module := core.NewModule(core.NewModuleOptions{
		Imports: []core.Modules{microservices.RegisterClient(tcp.NewClient(microservices.ConnectOptions{
			Addr: "localhost:8000",
		}))},
	})

	require.NotNil(t, microservices.Inject(module))

	module2 := core.NewModule(core.NewModuleOptions{})
	require.Nil(t, microservices.Inject(module2))
}
