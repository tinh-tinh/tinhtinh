package logger_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/middleware/logger"
)

func Test_Module(t *testing.T) {
	module := core.NewModule(core.NewModuleOptions{
		Imports: []core.Module{
			logger.Module(logger.Options{
				Max:    1,
				Rotate: true,
			}),
		},
	})

	require.NotNil(t, logger.InjectLog(module))
}
