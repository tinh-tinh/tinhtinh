package microservices

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/common/compress"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/dto/validator"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/logger"
)

type Header map[string]string

type Config struct {
	Serializer       core.Encode
	Deserializer     core.Decode
	Header           Header
	ErrorHandler     ErrorHandler
	Logger           *logger.Logger
	CompressAlg      compress.Alg
	RetryOptions     RetryOptions
	CustomValidation core.PipeFnc
}

type RetryOptions struct {
	Retry int
	Delay time.Duration
}

func DefaultConfig() Config {
	logger := logger.Create(logger.Options{})
	return Config{
		Serializer:       json.Marshal,
		Deserializer:     json.Unmarshal,
		Header:           make(Header),
		ErrorHandler:     DefaultErrorHandler(logger),
		Logger:           logger,
		CustomValidation: validator.Scanner,
	}
}

func ParseConfig(cfg ...Config) Config {
	defaultConfig := DefaultConfig()
	if len(cfg) > 0 {
		mergeConfig := common.MergeStruct(cfg...)
		if mergeConfig.Serializer != nil {
			defaultConfig.Serializer = mergeConfig.Serializer
		}

		if mergeConfig.Deserializer != nil {
			defaultConfig.Deserializer = mergeConfig.Deserializer
		}

		if len(cfg[0].Header) > 0 {
			for k, v := range mergeConfig.Header {
				defaultConfig.Header[k] = v
			}
		}

		if mergeConfig.ErrorHandler != nil {
			defaultConfig.ErrorHandler = mergeConfig.ErrorHandler
		}

		if mergeConfig.Logger != nil {
			defaultConfig.Logger = mergeConfig.Logger
			if mergeConfig.ErrorHandler == nil {
				defaultConfig.ErrorHandler = DefaultErrorHandler(cfg[0].Logger)
			}
		}

		if mergeConfig.CompressAlg != "" {
			defaultConfig.CompressAlg = mergeConfig.CompressAlg
		}

		if !reflect.ValueOf(mergeConfig.RetryOptions).IsZero() {
			defaultConfig.RetryOptions = mergeConfig.RetryOptions
		}

		if mergeConfig.CustomValidation != nil {
			defaultConfig.CustomValidation = mergeConfig.CustomValidation
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
