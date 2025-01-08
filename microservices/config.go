package microservices

import (
	"encoding/json"
	"reflect"

	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/logger"
)

type Header map[string]string

type Config struct {
	Serializer   core.Encode
	Deserializer core.Decode
	Header       Header
	ErrorHandler ErrorHandler
	Logger       *logger.Logger
}

func DefaultConfig() Config {
	logger := logger.Create(logger.Options{})
	return Config{
		Serializer:   json.Marshal,
		Deserializer: json.Unmarshal,
		Header:       make(Header),
		ErrorHandler: DefaultErrorHandler(logger),
		Logger:       logger,
	}
}

func ParseConfig(cfg ...Config) Config {
	defaultConfig := DefaultConfig()
	if len(cfg) > 0 {
		if cfg[0].Serializer != nil {
			defaultConfig.Serializer = cfg[0].Serializer
		}

		if cfg[0].Deserializer != nil {
			defaultConfig.Deserializer = cfg[0].Deserializer
		}

		if len(cfg[0].Header) > 0 {
			for k, v := range cfg[0].Header {
				defaultConfig.Header[k] = v
			}
		}

		if cfg[0].ErrorHandler != nil {
			defaultConfig.ErrorHandler = cfg[0].ErrorHandler
		}

		if cfg[0].Logger != nil {
			defaultConfig.Logger = cfg[0].Logger
			if cfg[0].ErrorHandler == nil {
				defaultConfig.ErrorHandler = DefaultErrorHandler(cfg[0].Logger)
			}
		}
	}

	return defaultConfig
}

func NewConfig(config Config) Config {
	if reflect.ValueOf(config).IsZero() {
		config = DefaultConfig()
	} else {
		config = ParseConfig(config)
	}
	return config
}
