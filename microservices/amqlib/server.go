package amqlib

import (
	"context"
	"fmt"
	"reflect"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tinh-tinh/tinhtinh/microservices"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

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
			fmt.Printf("Connected to RabbitMQ at %s\n", opts[0].Addr)
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
		"logs",        // name
		"fanout",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		panic(err)
	}

	if store.Subscribers != nil {
		for _, sub := range store.Subscribers {
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

			err = ch.QueueBind(
				q.Name,        // queue name
				"",            // routing key
				"logs",        // exchange
				false,
				nil)
			if err != nil {
				panic(err)
			}
			go c.HandlerPubSub(ch, q.Name, sub)
		}
	}

	if len(store.RpcHandlers) > 0 {
		q, err := ch.QueueDeclare(
			"rpc_queue",   // name
			false,         // durable
			false,         // delete when unused
			false,         // exclusive
			false,         // no-wait
			nil,           // arguments
		)
		if err != nil {
			panic(err)
		}
		go c.HandlerRPC(ch, q.Name, store.RpcHandlers)
	}
}

func (c *Connect) HandlerPubSub(ch *amqp091.Channel, q string, sub *microservices.SubscribeHandler) {
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
		fmt.Printf("Error consuming messages: %v\n", err)
		return
	}
	for d := range msgs {
		message := microservices.DecodeMessage(c, d.Body)
		if message.Event != sub.Name {
			continue
		}
		fmt.Printf("Received message: %v for event %s\n", message.Data, message.Event)
		if reflect.ValueOf(message).IsZero() {
			sub.Handle(c, microservices.Message{
				Bytes: d.Body,
			})
		} else {
			sub.Handle(c, message)
		}
	}
}

func (c *Connect) HandlerRPC(ch *amqp091.Channel, q string, subs microservices.RpcHandlers) {
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
		message := microservices.DecodeMessage(c, d.Body)
		var ctx microservices.Ctx
		if reflect.ValueOf(message).IsZero() {
			ctx = microservices.NewCtx(microservices.Message{
				Data: d.Body,
			}, c)
		} else {
			ctx = microservices.NewCtx(message, c)
		}

		var sub *microservices.RpcHandler
		for _, s := range subs {
			if s.Name == message.Event {
				sub = s
				break
			}
		}

		if sub == nil {
			continue
		}

		reply, err := sub.Factory(ctx)
		if err != nil {
			ctx.ErrorHandler(err)
			continue
		}

		if d.ReplyTo != "" {
			err = ch.PublishWithContext(c.Context,
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp091.Publishing{
					ContentType:   "application/json",
					CorrelationId: d.CorrelationId,
					Body:          reply,
				})
			if err != nil {
				ctx.ErrorHandler(err)
			}
		}
	}
}
