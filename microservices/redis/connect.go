package redis

import (
	"context"
	"reflect"
	"strings"

	redis_store "github.com/redis/go-redis/v9"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

type Options struct {
	microservices.Config
	*redis_store.Options
}

type Connect struct {
	Context context.Context
	Module  core.Module
	Conn    *redis_store.Client
	config  microservices.Config
}

// Client usage
func NewClient(opt Options) microservices.ClientProxy {
	conn := redis_store.NewClient(opt.Options)

	connect := &Connect{
		Context: context.Background(),
		Conn:    conn,
		config:  microservices.NewConfig(opt.Config),
	}

	if err := connect.Conn.Ping(connect.Context).Err(); err != nil {
		panic(err)
	}

	return connect
}

func (c *Connect) Headers() microservices.Header {
	return c.config.Header
}

func (c *Connect) Serializer(v interface{}) ([]byte, error) {
	return c.config.Serializer(v)
}

func (c *Connect) Deserializer(data []byte, v interface{}) error {
	return c.config.Deserializer(data, v)
}

func (c *Connect) ErrorHandler(err error) {
	c.config.ErrorHandler(err)
}

func (c *Connect) Send(event string, data interface{}, headers ...microservices.Header) error {
	return microservices.DefaultSend(c)(event, data, headers...)
}

func (c *Connect) Publish(event string, data interface{}, headers ...microservices.Header) error {
	return microservices.DefaultPublish(c)(event, data, headers...)
}

func (c *Connect) Emit(event string, message microservices.Message) error {
	payload, err := microservices.EncodeMessage(c, message)
	if err != nil {
		c.ErrorHandler(err)
		return err
	}
	err = c.Conn.Publish(c.Context, event, payload).Err()
	if err != nil {
		c.ErrorHandler(err)
		return err
	}
	return nil
}

// Server usage
func New(module core.ModuleParam, opts ...Options) microservices.Service {
	connect := &Connect{
		Context: context.Background(),
		Module:  module(),
		config:  microservices.DefaultConfig(),
	}

	if len(opts) > 0 {
		if opts[0].Options != nil {
			conn := redis_store.NewClient(opts[0].Options)
			connect.Conn = conn
		}
		if !reflect.ValueOf(opts[0].Config).IsZero() {
			connect.config = microservices.ParseConfig(opts[0].Config)
		}
	}

	if err := connect.Conn.Ping(connect.Context).Err(); err != nil {
		panic(err)
	}

	return connect
}

func Open(opts ...Options) core.Service {
	connect := &Connect{
		Context: context.Background(),
		config:  microservices.DefaultConfig(),
	}

	if len(opts) > 0 {
		if opts[0].Options != nil {
			conn := redis_store.NewClient(opts[0].Options)
			connect.Conn = conn
		}
		if !reflect.ValueOf(opts[0].Config).IsZero() {
			connect.config = microservices.ParseConfig(opts[0].Config)
		}
	}

	return connect
}

func (c *Connect) Create(module core.Module) {
	c.Module = module
}

func (c *Connect) Listen() {
	store, ok := c.Module.Ref(microservices.STORE).(*microservices.Store)
	if !ok {
		panic("store not found")
	}

	if store.GetRPC() != nil {
		for _, sub := range store.GetRPC() {
			subscriber := c.Conn.Subscribe(c.Context, sub.Name)
			go c.Handler(subscriber, sub)
		}
	}

	if store.GetPubSub() != nil {
		var subscriber *redis_store.PubSub
		for _, sub := range store.GetPubSub() {
			if strings.HasSuffix(sub.Name, "*") {
				subscriber = c.Conn.PSubscribe(c.Context, sub.Name)
			} else {
				subscriber = c.Conn.Subscribe(c.Context, sub.Name)
			}
			go c.Handler(subscriber, sub)
		}
	}
}

func (c *Connect) Handler(subscriber *redis_store.PubSub, sub *microservices.SubscribeHandler) {
	for {
		msg, err := subscriber.ReceiveMessage(c.Context)
		if err != nil {
			return
		}

		message := microservices.DecodeMessage(c, []byte(msg.Payload))
		sub.Handle(c, message)
	}
}
