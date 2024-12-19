package redis_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/microservices/redis"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

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
