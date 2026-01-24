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
	Payload() any
	PayloadParser(schema any) error
	ErrorHandler(err error)
	Set(key any, value any)
	Get(key any) any
	Next() error
	Scan(val any) error
	Path() string
	Reply(data any) ([]byte, error)
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

func (c *DefaultCtx) unmarshallCompress(data interface{}) error {
	buf := bytes.NewBuffer(c.message.Bytes)
	dec := gob.NewDecoder(buf)

	for {
		err := dec.Decode(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *DefaultCtx) Payload() interface{} {
	return c.message.Data
}

func (c *DefaultCtx) PayloadParser(schema interface{}) error {
	payload := c.message.Data

	if c.message.Data == nil && c.message.Bytes != nil {
		return c.unmarshallCompress(schema)
	}
	if reflect.TypeOf(payload).Kind() == reflect.String {
		return c.service.Config().Deserializer([]byte(payload.(string)), schema)
	}
	dataBytes, _ := c.service.Config().Serializer(payload)
	return c.service.Config().Deserializer(dataBytes, schema)
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

func (c *DefaultCtx) Scan(val any) error {
	return c.service.Config().CustomValidation(val)
}

func (c *DefaultCtx) Path() string {
	return c.message.Event
}

func (c *DefaultCtx) Reply(data any) ([]byte, error) {
	return c.service.Config().Serializer(data)
}
