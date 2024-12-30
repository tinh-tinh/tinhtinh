package tcp

import (
	"bufio"
	"fmt"
	"net"
	"reflect"

	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

type Client struct {
	config microservices.Config
	Conn   net.Conn
}

func NewClient(opt microservices.TcpOptions) microservices.ClientProxy {
	conn, err := net.Dial("tcp", opt.Addr)
	if err != nil {
		panic(err)
	}
	client := &Client{
		Conn:   conn,
		config: opt.Config,
	}

	if reflect.ValueOf(client.config).IsZero() {
		client.config = microservices.ParseConfig(opt.Config)
	}

	return client
}

func (client *Client) Send(event string, data interface{}, headers ...microservices.Header) error {
	writer := bufio.NewWriter(client.Conn)

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
		return err
	}

	jsonData = append(jsonData, '\n')
	_, err = writer.Write(jsonData)
	if err != nil {
		return err
	}

	writer.Flush()
	fmt.Printf("Send message: %v for event %s\n", data, event)
	return nil
}

func (client *Client) Publish(event string, data interface{}, headers ...microservices.Header) error {
	writer := bufio.NewWriter(client.Conn)

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
		return err
	}

	jsonData = append(jsonData, '\n')
	_, err = writer.Write(jsonData)
	if err != nil {
		return err
	}

	writer.Flush()
	fmt.Printf("Publish message: %v for event %s\n", data, event)
	return nil
}
