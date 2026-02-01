package microservices

import (
	"encoding/json"
	"time"

	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/common/compress"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/dto/validator"
)

type Header map[string]string

type Config struct {
	Serializer       core.Encode
	Deserializer     core.Decode
	Header           Header
	ErrorHandler     ErrorHandler
	CompressAlg      compress.Alg
	RetryOptions     RetryOptions
	CustomValidation core.PipeFnc
}

type RetryOptions struct {
	Retry int
	Delay time.Duration
}

// DefaultConfig returns a Config instance with default serialization, deserialization, error handling, header, and validation settings.
func DefaultConfig() Config {
	return Config{
		Serializer:       json.Marshal,
		Deserializer:     json.Unmarshal,
		Header:           make(Header),
		ErrorHandler:     DefaultErrorHandler(),
		CustomValidation: validator.Scanner,
	}
}

// ParseConfig merges one or more Config instances into a single Config, overriding default values with non-zero fields from the provided configs.
// If multiple configs are provided, later configs take precedence for overlapping fields. Header maps are merged by key. Returns the resulting merged Config.
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

		if mergeConfig.CompressAlg != "" {
			defaultConfig.CompressAlg = mergeConfig.CompressAlg
		}

		if mergeConfig.RetryOptions.Delay != 0 || mergeConfig.RetryOptions.Retry != 0 {
			defaultConfig.RetryOptions = mergeConfig.RetryOptions
		}

		if mergeConfig.CustomValidation != nil {
			defaultConfig.CustomValidation = mergeConfig.CustomValidation
		}
	}

	return defaultConfig
}

func NewConfig(config Config) Config {
	if config.IsZero() {
		return DefaultConfig()
	}
	return ParseConfig(config)
}

func (c Config) IsZero() bool {
	return c.Serializer == nil &&
		c.Deserializer == nil &&
		len(c.Header) == 0 &&
		c.ErrorHandler == nil &&
		c.CompressAlg == "" &&
		c.CustomValidation == nil &&
		c.RetryOptions.Delay == 0 &&
		c.RetryOptions.Retry == 0
}
