package amqlib

import (
	"context"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tinh-tinh/tinhtinh/microservices"
)

// Client usage
func NewClient(opt Options) microservices.ClientProxy {
	conn, err := amqp091.Dial(opt.Addr)
	if err != nil {
		fmt.Printf("Failed to connect to RabbitMQ: %v\n", err)
		panic(err)
	}

	fmt.Printf("Connected to RabbitMQ at %s\n", opt.Addr)
	connect := &Connect{
		Context: context.Background(),
		Conn:    conn,
		config:  microservices.NewConfig(opt.Config),
	}

	return connect
}

func (c *Connect) Config() microservices.Config {
	return c.config
}

func (c *Connect) Emit(event string, message microservices.Message) error {
	return nil
}

func (c *Connect) Timeout(duration time.Duration) microservices.ClientProxy {
	c.timeout = duration
	return c
}

func (c *Connect) Send(event string, data any, res any, headers ...microservices.Header) error {
	defer c.Conn.Close()

	ch, err := c.Conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// No ExchangeDeclare needed for RPC, using default exchange

	if c.timeout == 0 {
		c.timeout = 5 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	body, err := microservices.EncodeMessage(c, microservices.Message{
		Event:   event,
		Data:    data,
		Headers: microservices.AssignHeader(c.config.Header, headers...),
	})
	if err != nil {
		return err
	}
	err = ch.PublishWithContext(ctx,
		"",           // exchange
		"rpc_queue",  // routing key
		false,         // mandatory
		false,         // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return err
	}

	c.timeout = 0
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
		"logs",        // name
		"fanout",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return err
	}

	if c.timeout == 0 {
		c.timeout = 5 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	body, err := microservices.EncodeMessage(c, microservices.Message{
		Event:   event,
		Data:    data,
		Headers: microservices.AssignHeader(c.config.Header, headers...),
	})
	if err != nil {
		return err
	}
	err = ch.PublishWithContext(ctx,
		"logs",        // exchange
		"",            // routing key
		false,         // mandatory
		false,         // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return err
	}

	c.timeout = 0
	return nil
}
