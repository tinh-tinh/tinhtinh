package tcp

import (
	"fmt"
	"net"
	"reflect"
	"time"

	"github.com/tinh-tinh/tinhtinh/v2/common"
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
	conn, err := net.Dial("tcp", opt.Addr)
	if err != nil {
		panic(err)
	}

	if reflect.ValueOf(opt.Config).IsZero() {
		opt.Config = microservices.DefaultConfig()
	}

	if opt.Timeout > 0 {
		conn.SetDeadline(time.Now().Add(opt.Timeout))
	}

	client := &Client{
		Conn:   conn,
		config: opt.Config,
	}

	return client
}

func (client *Client) Send(event string, data interface{}, headers ...microservices.Header) error {
	message := microservices.Message{
		Type:    microservices.RPC,
		Headers: common.CloneMap(client.config.Header),
		Event:   event,
		Data:    data,
	}
	if len(headers) > 0 {
		for _, v := range headers {
			common.MergeMaps(message.Headers, v)
		}
	}

	jsonData, err := client.config.Serializer(message)
	if err != nil {
		client.config.ErrorHandler(err)
		return err
	}

	jsonData = append(jsonData, '\n')
	_, err = client.Conn.Write(jsonData)
	if err != nil {
		client.config.ErrorHandler(err)
		return err
	}

	fmt.Printf("Send message: %v for event %s\n", data, event)
	return nil
}

func (client *Client) Publish(event string, data interface{}, headers ...microservices.Header) error {
	message := microservices.Message{
		Type:    microservices.PubSub,
		Headers: common.CloneMap(client.config.Header),
		Event:   event,
		Data:    data,
	}
	if len(headers) > 0 {
		for _, v := range headers {
			common.MergeMaps(message.Headers, v)
		}
	}

	jsonData, err := client.config.Serializer(message)
	if err != nil {
		return err
	}

	jsonData = append(jsonData, '\n')

	_, err = client.Conn.Write(jsonData)
	if err != nil {
		client.config.ErrorHandler(err)
		return err
	}

	fmt.Printf("Publish message: %v for event %s\n", data, event)
	return nil
}
