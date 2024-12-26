package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	message := microservices.Message{Type: microservices.RPC, Event: event, Data: data}
	payload, err := c.Serializer(message)
	if err != nil {
		return err
	}
	err = c.Conn.Publish(c.Context, event, payload).Err()
	if err != nil {
		return err
	}
	fmt.Printf("Send mesage: %v for event: %s\n", data, event)
	return nil
}

func (c *Connect) Publish(event string, data interface{}) error {
	message := microservices.Message{Type: microservices.PubSub, Event: event, Data: data}
	payload, err := c.Serializer(message)
	if err != nil {
		return err
	}
	err = c.Conn.Publish(c.Context, event, payload).Err()
	if err != nil {
		return err
	}
	fmt.Printf("Publish mesage: %v for event %s\n", data, event)
	return nil
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
	store := c.Module.Ref(microservices.STORE).(*microservices.Store)
	if store == nil {
		panic("store not found")
	}

	if store.Subscribers[string(microservices.RPC)] != nil {
		for _, sub := range store.Subscribers[string(microservices.RPC)] {
			subscriber := c.Conn.Subscribe(c.Context, sub.Name)
			go c.Handler(subscriber, sub)
		}
	}

	if store.Subscribers[string(microservices.PubSub)] != nil {
		var subscriber *redis_store.PubSub
		for _, sub := range store.Subscribers[string(microservices.PubSub)] {
			if strings.HasSuffix(sub.Name, "*") {
				subscriber = c.Conn.PSubscribe(c.Context, sub.Name)
			} else {
				subscriber = c.Conn.Subscribe(c.Context, sub.Name)
			}
			go c.Handler(subscriber, sub)
		}
	}
}

func (c *Connect) Handler(subscriber *redis_store.PubSub, sub microservices.SubscribeHandler) {
	for {
		msg, err := subscriber.ReceiveMessage(c.Context)
		if err != nil {
			return
		}
		var message microservices.Message
		err = c.Deserializer([]byte(msg.Payload), &message)
		if err != nil {
			fmt.Println("Error deserializing message: ", err)
			return
		}

		sub.Handle(c, message.Data)
		// data := microservices.ParseCtx(message.Data, c)
		// factory(data)
	}
}

func (c *Connect) Serializer(v interface{}) ([]byte, error) {
	return c.serializer(v)
}

func (c *Connect) Deserializer(data []byte, v interface{}) error {
	return c.deserializer(data, v)
}

func (c *Connect) ErrorHandler(err error) {
	log.Printf("Error when running tcp: %v\n", err)
}
