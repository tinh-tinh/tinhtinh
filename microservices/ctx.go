package microservices

import (
	"bytes"
	"context"
	"encoding/gob"
	"io"
	"reflect"
)

type Ctx interface {
	Headers(key string) string
	Payload(data ...interface{}) interface{}
	ErrorHandler(err error)
	Set(key interface{}, value interface{})
	Get(key interface{}) interface{}
	Next() error
}

type DefaultCtx struct {
	message Message
	service Service
	context context.Context
}

func NewCtx(data Message, service Service) Ctx {
	return &DefaultCtx{
		message: data,
		service: service,
		context: context.Background(),
	}
}

func (c *DefaultCtx) marshallCompress(data interface{}) {
	buf := bytes.NewBuffer(c.message.Bytes)
	dec := gob.NewDecoder(buf)

	for {
		err := dec.Decode(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
	}
}

func (c *DefaultCtx) Payload(data ...interface{}) interface{} {
	payload := c.message.Data
	if len(data) > 0 {
		schema := data[0]
		if c.message.Data == nil && c.message.Bytes != nil {
			c.marshallCompress(schema)
			return schema
		}
		if reflect.TypeOf(payload).Kind() == reflect.String {
			_ = c.service.Config().Deserializer([]byte(payload.(string)), schema)
			return schema
		}
		dataBytes, _ := c.service.Config().Serializer(payload)
		_ = c.service.Config().Deserializer(dataBytes, schema)
		return schema
	}
	return payload
}

func (c *DefaultCtx) ErrorHandler(err error) {
	c.service.Config().ErrorHandler(err)
}

func (c *DefaultCtx) Next() error {
	return nil
}

func (c *DefaultCtx) Set(key interface{}, val interface{}) {
	ctx := context.WithValue(c.context, key, val)
	c.context = ctx
}

func (c *DefaultCtx) Get(key interface{}) interface{} {
	return c.context.Value(key)
}

func (c *DefaultCtx) Headers(key string) string {
	return c.message.Headers[key]
}
