package microservices

import (
	"fmt"

	"github.com/tinh-tinh/tinhtinh/v2/common/exception"
)

type ErrorHandler func(err error)

// DefaultErrorHandler returns an ErrorHandler that prints formatted RPC error codes and messages to standard output.
func DefaultErrorHandler() ErrorHandler {
	return func(err error) {
		instance := exception.AdapterRpcError(err)

		content := instance.Code + ": " + instance.Msg
		fmt.Println(content)
	}
}
