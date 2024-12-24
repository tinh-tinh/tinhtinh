package microservices

import (
	"reflect"
)

type Ctx interface {
	Payload(data ...interface{}) interface{}
}

type Factory func(ctx Ctx)

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
		dataBytes, _ := c.service.Serializer(c.payload)
		_ = c.service.Deserializer(dataBytes, schema)
		return schema
	}
	return c.payload
}
