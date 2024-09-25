package root

import (
	"fmt"
	"net/http"

	"github.com/tinh-tinh/tinhtinh/core"
)

func NewProvider(module *core.DynamicModule) *core.DynamicProvider {
	provider := module.NewProvider(core.ProviderOptions{
		Name: "root",
		Factory: func(param ...interface{}) interface{} {
			req := param[0].(*http.Request)
			return fmt.Sprintf("%vRoot", req.Header.Get("x-api-key"))
		},
		Inject: []core.Provide{core.REQUEST},
	})

	return provider
}
