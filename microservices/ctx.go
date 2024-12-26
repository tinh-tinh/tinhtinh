package microservices

import (
	"context"
	"reflect"
)

type Ctx interface {
	Payload(data ...interface{}) interface{}
	ErrorHandler(err error)
	Set(key interface{}, value interface{})
	Get(key interface{}) interface{}
	Next() error
	SetFactory(f Factory)
}

type DefaultCtx struct {
	payload interface{}
	service Service
	factory Factory
	context context.Context
}

func NewCtx(data interface{}, service Service) Ctx {
	return &DefaultCtx{
		payload: data,
		service: service,
		context: context.Background(),
	}
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

func (c *DefaultCtx) ErrorHandler(err error) {
	c.service.ErrorHandler(err)
}

func (c *DefaultCtx) Next() error {
	return c.factory.Handle(c)
}

func (c *DefaultCtx) Set(key interface{}, val interface{}) {
	ctx := context.WithValue(c.context, key, val)
	c.context = ctx
}

func (c *DefaultCtx) Get(key interface{}) interface{} {
	return c.context.Value(key)
}

func (c *DefaultCtx) SetFactory(factory Factory) {
	c.factory = factory
}
