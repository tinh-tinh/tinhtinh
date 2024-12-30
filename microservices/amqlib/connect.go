package amqlib

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

type Options struct {
	microservices.Config
	Addr string
}

type Connect struct {
	clientHeaders map[string]string
	Conn          *amqp091.Connection
	Module        core.Module
	Context       context.Context
	config        microservices.Config
}

// Client usage
func NewClient(opt Options) microservices.ClientProxy {
	conn, err := amqp091.Dial(opt.Addr)
	if err != nil {
		panic(err)
	}

	connect := &Connect{
		Context: context.Background(),
		Conn:    conn,
		config:  opt.Config,
	}

	if reflect.ValueOf(connect.config).IsZero() {
		connect.config = microservices.DefaultConfig()
	}

	return connect
}

func (c *Connect) Send(event string, data interface{}, headers ...microservices.Header) error {
	defer c.Conn.Close()

	ch, err := c.Conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_direct", // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message := microservices.Message{
		Type:    microservices.RPC,
		Event:   event,
		Data:    data,
		Headers: common.CloneMap(c.config.Header),
	}

	if len(headers) > 0 {
		for _, v := range headers {
			common.MergeMaps(message.Headers, v)
		}
	}

	body, err := c.Serializer(message)
	if err != nil {
		return err
	}
	err = ch.PublishWithContext(ctx,
		"logs_direct", // exchange
		event,         // routing key
		false,         // mandatory
		false,         // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return err
	}

	fmt.Printf("Send message: %v for event %s\n", data, event)
	return nil
}

func (c *Connect) Publish(event string, data interface{}, headers ...microservices.Header) error {
	defer c.Conn.Close()

	ch, err := c.Conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_direct", // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message := microservices.Message{
		Type:    microservices.PubSub,
		Event:   event,
		Data:    data,
		Headers: c.clientHeaders,
	}
	body, err := c.Serializer(message)
	if err != nil {
		return err
	}
	err = ch.PublishWithContext(ctx,
		"logs_direct", // exchange
		event,         // routing key
		false,         // mandatory
		false,         // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}

func (c *Connect) Serializer(v interface{}) ([]byte, error) {
	return c.config.Serializer(v)
}

func (c *Connect) Deserializer(data []byte, v interface{}) error {
	return c.config.Deserializer(data, v)
}

func (c *Connect) SetHeaders(key string, value string) microservices.ClientProxy {
	c.clientHeaders[key] = value
	return c
}

func (c *Connect) GetHeaders(key string) string {
	return c.clientHeaders[key]
}

// Server usage
func New(module core.ModuleParam, opts ...Options) core.Service {
	connect := &Connect{
		Module:  module(),
		config:  microservices.DefaultConfig(),
		Context: context.Background(),
	}

	if len(opts) > 0 {
		if opts[0].Addr != "" {
			conn, err := amqp091.Dial(opts[0].Addr)
			if err != nil {
				panic(err)
			}
			connect.Conn = conn
		}
		if !reflect.ValueOf(opts[0].Config).IsZero() {
			connect.config = microservices.ParseConfig(opts[0].Config)
		}
	}

	return connect
}

func Open(opts ...Options) core.Service {
	connect := &Connect{
		config:  microservices.DefaultConfig(),
		Context: context.Background(),
	}

	if len(opts) > 0 {
		if opts[0].Addr != "" {
			conn, err := amqp091.Dial(opts[0].Addr)
			if err != nil {
				panic(err)
			}
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
	fmt.Println("Listening to RabbitMQ")
	store := c.Module.Ref(microservices.STORE).(*microservices.Store)
	if store == nil {
		panic("store not found")
	}

	ch, err := c.Conn.Channel()
	if err != nil {
		panic(err)
	}

	err = ch.ExchangeDeclare(
		"logs_direct", // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		panic(err)
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		panic(err)
	}

	if store.Subscribers[string(microservices.RPC)] != nil {
		for _, sub := range store.Subscribers[string(microservices.RPC)] {
			err = ch.QueueBind(
				q.Name,        // queue name
				sub.Name,      // routing key
				"logs_direct", // exchange
				false,
				nil)
			if err != nil {
				panic(err)
			}
			go c.Handler(ch, q.Name, sub)
		}
	}

	if store.Subscribers[string(microservices.PubSub)] != nil {
		for _, sub := range store.Subscribers[string(microservices.PubSub)] {
			err = ch.QueueBind(
				q.Name,        // queue name
				sub.Name,      // routing key
				"logs_direct", // exchange
				false,
				nil)
			if err != nil {
				panic(err)
			}
			go c.Handler(ch, q.Name, sub)
		}
	}
}

func (c *Connect) Handler(ch *amqp091.Channel, q string, sub microservices.SubscribeHandler) {
	msgs, err := ch.Consume(
		q,     // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return
	}
	for d := range msgs {
		var message microservices.Message
		err := c.Deserializer(d.Body, &message)
		if err != nil {
			fmt.Println("Error deserializing message: ", err)
			continue
		}
		fmt.Println(message)
		if reflect.ValueOf(message).IsZero() {
			sub.Handle(c, microservices.Message{
				Data: d.Body,
			})
		} else {
			sub.Handle(c, message)
		}
	}
}

func (c *Connect) ErrorHandler(err error) {
	log.Printf("Error when running tcp: %v\n", err)
}
