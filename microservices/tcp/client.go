package tcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"

	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

type Client struct {
	Conn         net.Conn
	Serializer   core.Encode
	Deserializer core.Decode
	Wildcard     bool
	Delimiter    string
}

func NewClient(opt microservices.ConnectOptions) microservices.ClientProxy {
	conn, err := net.Dial("tcp", opt.Addr)
	if err != nil {
		panic(err)
	}
	client := &Client{
		Conn:         conn,
		Serializer:   json.Marshal,
		Deserializer: json.Unmarshal,
	}
	if opt.Deserializer != nil {
		client.Deserializer = opt.Deserializer
	}
	if opt.Serializer != nil {
		client.Serializer = opt.Serializer
	}

	return client
}

func (client *Client) Send(event string, data interface{}) error {
	writer := bufio.NewWriter(client.Conn)

	message := microservices.Message{Type: microservices.RPC, Event: event, Data: data}
	jsonData, err := client.Serializer(message)
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

func (client *Client) Publish(event string, data interface{}) error {
	writer := bufio.NewWriter(client.Conn)

	message := microservices.Message{Type: microservices.PubSub, Event: event, Data: data}
	jsonData, err := client.Serializer(message)
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
