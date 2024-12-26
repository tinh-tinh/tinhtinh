package amqlib

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

type Connect struct {
	Conn         *amqp091.Connection
	Module       core.Module
	Context      context.Context
	serializer   core.Encode
	deserializer core.Decode
}

// Client usage
func NewClient(opt microservices.ConnectOptions) microservices.ClientProxy {
	conn, err := amqp091.Dial(opt.Addr)
	if err != nil {
		panic(err)
	}

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

func (c *Connect) Send(event string, data interface{}) error {
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

	message := microservices.Message{Type: microservices.RPC, Event: event, Data: data}
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

func (c *Connect) Publish(event string, data interface{}) error {
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

	message := microservices.Message{Type: microservices.RPC, Event: event, Data: data}
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
	return c.serializer(v)
}

func (c *Connect) Deserializer(data []byte, v interface{}) error {
	return c.deserializer(data, v)
}

// Server usage
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
			conn, err := amqp091.Dial(opts[0].Addr)
			if err != nil {
				panic(err)
			}
			connect.Conn = conn
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
		if !reflect.ValueOf(message).IsZero() {
			sub.Handle(c, message.Data)
			// data := microservices.ParseCtx(message.Data, c)
			// factory(data)
		} else {
			sub.Handle(c, message.Data)
		}
	}
}

func (c *Connect) ErrorHandler(err error) {
	log.Printf("Error when running tcp: %v\n", err)
}
