package redis

import (
	"context"
	"reflect"
	"strings"
	"time"

	redis_store "github.com/redis/go-redis/v9"
	"github.com/tinh-tinh/tinhtinh/v2/common/era"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

type Options struct {
	microservices.Config
	*redis_store.Options
}

type Connect struct {
	Module  core.Module
	Context context.Context
	Conn    *redis_store.Client
	config  microservices.Config
	timeout time.Duration
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
		if opt.RetryOptions.Retry != 0 {
			time.Sleep(opt.RetryOptions.Delay)
			opt.RetryOptions.Retry--
			return NewClient(opt)
		}
		panic(err)
	}

	return connect
}

func (c *Connect) Config() microservices.Config {
	return c.config
}

func (c *Connect) Send(event string, data interface{}, headers ...microservices.Header) error {
	return microservices.DefaultSend(c)(event, data, headers...)
}

func (c *Connect) Publish(event string, data interface{}, headers ...microservices.Header) error {
	return microservices.DefaultPublish(c)(event, data, headers...)
}

func (c *Connect) Timeout(duration time.Duration) microservices.ClientProxy {
	c.timeout = duration
	return c
}

func (c *Connect) Emit(event string, message microservices.Message) error {
	payload, err := microservices.EncodeMessage(c, message)
	if err != nil {
		c.config.ErrorHandler(err)
		return err
	}
	if c.timeout > 0 {
		err = era.TimeoutFunc(c.timeout, func(ctx context.Context) error {
			err := c.Conn.Publish(ctx, event, payload).Err()
			return err
		})
	} else {
		err = c.Conn.Publish(c.Context, event, payload).Err()
	}
	if err != nil {
		c.config.ErrorHandler(err)
		return err
	}
	c.timeout = 0
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
