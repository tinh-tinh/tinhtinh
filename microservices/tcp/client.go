package tcp

import (
	"fmt"
	"net"
	"reflect"
	"time"

	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

type Options struct {
	microservices.Config
	Addr    string
	Timeout time.Duration
}

type Client struct {
	config microservices.Config
	Conn   net.Conn
}

func NewClient(opt Options) microservices.ClientProxy {
	var conn net.Conn
	var err error
	if opt.Timeout > 0 {
		conn, err = net.DialTimeout("tcp", opt.Addr, opt.Timeout)
	} else {
		conn, err = net.Dial("tcp", opt.Addr)
	}
	if err != nil {
		panic(err)
	}

	if reflect.ValueOf(opt.Config).IsZero() {
		opt.Config = microservices.DefaultConfig()
	} else {
		opt.Config = microservices.ParseConfig(opt.Config)
	}

	client := &Client{
		Conn:   conn,
		config: opt.Config,
	}

	return client
}

func (client *Client) Timeout(duration time.Duration) microservices.ClientProxy {
	err := client.Conn.SetDeadline(time.Now().Add(duration))
	if err != nil {
		client.config.ErrorHandler(err)
	}
	return client
}

func (client *Client) Send(event string, data interface{}, headers ...microservices.Header) error {
	payload, err := microservices.EncodeMessage(client, microservices.Message{
		Type:    microservices.RPC,
		Headers: microservices.AssignHeader(client.config.Header, headers...),
		Event:   event,
		Data:    data,
	})
	if err != nil {
		client.Serializer(err)
		return err
	}
	payload = append(payload, '\n')
	_, err = client.Conn.Write(payload)
	if err != nil {
		client.config.ErrorHandler(err)
		return err
	}

	fmt.Printf("Send message: %v for event %s\n", data, event)
	return nil
}

func (client *Client) Publish(event string, data interface{}, headers ...microservices.Header) error {
	payload, err := microservices.EncodeMessage(client, microservices.Message{
		Type:    microservices.PubSub,
		Headers: microservices.AssignHeader(client.config.Header, headers...),
		Event:   event,
		Data:    data,
	})
	if err != nil {
		client.Serializer(err)
		return err
	}
	payload = append(payload, '\n')
	_, err = client.Conn.Write(payload)
	if err != nil {
		client.config.ErrorHandler(err)
		return err
	}

	fmt.Printf("Publish message: %v for event %s\n", data, event)
	return nil
}

func (c *Client) Serializer(v interface{}) ([]byte, error) {
	return c.config.Serializer(v)
}

func (c *Client) Deserializer(data []byte, v interface{}) error {
	return c.config.Deserializer(data, v)
}
