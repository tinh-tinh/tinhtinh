package amqlib

import (
	"context"
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

type Client struct {
	Conn         *amqp091.Connection
	Context      context.Context
	Serializer   core.Encode
	Deserializer core.Decode
}

func NewClient(opt microservices.ConnectOptions) microservices.ClientProxy {
	conn, err := amqp091.Dial(opt.Addr)
	if err != nil {
		panic(err)
	}

	client := &Client{
		Conn:         conn,
		Serializer:   json.Marshal,
		Deserializer: json.Unmarshal,
		Context:      context.Background(),
	}
	if opt.Deserializer != nil {
		client.Deserializer = opt.Deserializer
	}
	if opt.Serializer != nil {
		client.Serializer = opt.Serializer
	}
	return client
}

func (client *Client) Close() {
	client.Conn.Close()
}

func (client *Client) Send(event string, data interface{}) error {
	ch, err := client.Conn.Channel()
	if err != nil {
		return err
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
		return err
	}

	payload, err := client.Serializer(data)
	if err != nil {
		return err
	}
	err = ch.PublishWithContext(client.Context,
		"logs_direct", // exchange
		event,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        payload,
		},
	)
	return err
}

func (client *Client) Broadcast(data interface{}) error {
	return nil
}
