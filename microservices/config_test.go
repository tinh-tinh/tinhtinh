package microservices_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/logger"
)

func Test_Default(t *testing.T) {
	config := microservices.DefaultConfig()
	require.NotNil(t, config.Serializer)
	require.NotNil(t, config.Deserializer)
	require.NotNil(t, config.Logger)
	require.NotNil(t, config.ErrorHandler)
	require.NotNil(t, config.Header)

	cfg := microservices.ParseConfig(microservices.Config{
		Serializer: func(v interface{}) ([]byte, error) {
			return nil, nil
		},
		Deserializer: func(data []byte, v interface{}) error {
			return nil
		},
		ErrorHandler: func(err error) {},
		Logger:       logger.Create(logger.Options{}),
		Header:       map[string]string{"abc": "123"},
	})
	require.NotEqual(t, config, cfg)
}
