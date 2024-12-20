package microservices

import (
	"encoding/json"
	"reflect"

	"github.com/tinh-tinh/tinhtinh/v2/core"
)

type Ctx interface {
	Payload(data ...interface{}) interface{}
	Header() map[string]string
}

type Factory func(ctx Ctx)

func ConvertFactory(fnc Factory) core.Factory {
	return func(params ...interface{}) interface{} {
		if len(params) > 0 {
			fnc(params[0].(Ctx))
		}
		return nil
	}
}

type DefaultCtx struct {
	payload interface{}
	service Service
}

func ParseCtx(data interface{}, service Service) Ctx {
	return &DefaultCtx{payload: data, service: service}
}

func (c *DefaultCtx) Payload(data ...interface{}) interface{} {
	if len(data) > 0 {
		schema := data[0]
		if reflect.TypeOf(c.payload).Kind() == reflect.String {
			_ = c.service.Deserializer([]byte(c.payload.(string)), schema)
			return schema
		}
		dataBytes, _ := json.Marshal(c.payload)
		_ = c.service.Deserializer(dataBytes, schema)
		return schema
	}
	return c.payload
}

func (c *DefaultCtx) Header() map[string]string {
	return nil
}
