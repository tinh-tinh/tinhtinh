package microservices

import (
	"github.com/tinh-tinh/tinhtinh/v2/common/exception"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/logger"
)

type ErrorHandler func(err error)

func DefaultErrorHandler(logger *logger.Logger) ErrorHandler {
	return func(err error) {
		instance := exception.AdapterRpcError(err)

		content := instance.Code + ": " + instance.Msg
		logger.Error(content)
	}
}
