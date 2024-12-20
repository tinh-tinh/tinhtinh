package redis

import (
	"context"
	"encoding/json"
	"strings"

	redis_store "github.com/redis/go-redis/v9"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

type Connect struct {
	Context      context.Context
	Module       core.Module
	serializer   core.Encode
	deserializer core.Decode
	Conn         *redis_store.Client
}

// Client usage
func NewClient(opt microservices.ConnectOptions) microservices.ClientProxy {
	conn := redis_store.NewClient(&redis_store.Options{
		Addr:     opt.Addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	connect := &Connect{
		Context:      context.Background(),
		Conn:         conn,
		serializer:   json.Marshal,
		deserializer: json.Unmarshal,
	}
	if opt.Deserializer != nil {
		connect.deserializer = opt.Deserializer
	}
	if opt.Serializer != nil {
		connect.serializer = opt.Serializer
	}

	return connect
}

func (c *Connect) Close() {
	c.Conn.Close()
}

func (c *Connect) Send(event string, data interface{}) error {
	payload, err := c.Serializer(data)
	if err != nil {
		return err
	}
	err = c.Conn.Publish(c.Context, event, payload).Err()
	if err != nil {
		return err
	}
	return nil
}

func (client *Connect) Broadcast(data interface{}) error {
	return client.Send("*", data)
}

// Server usage
func New(module core.ModuleParam, opts ...microservices.ConnectOptions) microservices.Service {
	connect := &Connect{
		Module:       module(),
		serializer:   json.Marshal,
		deserializer: json.Unmarshal,
	}

	if len(opts) > 0 {
		if opts[0].Serializer != nil {
			connect.serializer = opts[0].Serializer
		}

		if opts[0].Deserializer != nil {
			connect.deserializer = opts[0].Deserializer
		}
		if opts[0].Addr != "" {
			conn := redis_store.NewClient(&redis_store.Options{
				Addr:     opts[0].Addr,
				Password: "", // no password set
				DB:       0,  // use default DB
			})
			connect.Conn = conn
		}
	}

	return connect
}

func Open(opts ...microservices.ConnectOptions) core.Service {
	connect := &Connect{
		serializer:   json.Marshal,
		deserializer: json.Unmarshal,
		Context:      context.Background(),
	}

	if len(opts) > 0 {
		if opts[0].Serializer != nil {
			connect.serializer = opts[0].Serializer
		}

		if opts[0].Deserializer != nil {
			connect.deserializer = opts[0].Deserializer
		}

		if opts[0].Addr != "" {
			conn := redis_store.NewClient(&redis_store.Options{
				Addr:     opts[0].Addr,
				Password: "", // no password set
				DB:       0,  // use default DB
			})
			connect.Conn = conn
		}
	}

	return connect
}

func (c *Connect) Create(module core.Module) {
	c.Module = module
}

func (c *Connect) Listen() {
	for _, prd := range c.Module.GetDataProviders() {
		var subscriber *redis_store.PubSub
		if strings.HasSuffix(string(prd.GetName()), "*") {
			subscriber = c.Conn.PSubscribe(c.Context, string(prd.GetName()))
		} else {
			subscriber = c.Conn.Subscribe(c.Context, string(prd.GetName()))
		}
		go c.Handler(subscriber, prd.GetFactory())
	}
}

func (c *Connect) Handler(params ...interface{}) {
	subscriber := params[0].(*redis_store.PubSub)
	factory := params[1].(core.Factory)
	for {
		msg, err := subscriber.ReceiveMessage(c.Context)
		if err != nil {
			return
		}

		data := microservices.ParseCtx(msg.Payload, c)
		factory(data)
	}
}

func (c *Connect) Serializer(v interface{}) ([]byte, error) {
	return c.serializer(v)
}

func (c *Connect) Deserializer(data []byte, v interface{}) error {
	return c.deserializer(data, v)
}
