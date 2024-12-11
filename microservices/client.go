package microservices

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"

	"github.com/tinh-tinh/tinhtinh/v2/core"
)

const CLIENT core.Provide = "CLIENT"

type Client struct {
	Conn         net.Conn
	Serializer   core.Encode
	Deserializer core.Decode
}

func RegisterClient(addr string, opts ...Options) core.Modules {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	client := &Client{
		Conn:         conn,
		Serializer:   json.Marshal,
		Deserializer: json.Unmarshal,
	}

	if len(opts) > 0 {
		if opts[0].Serializer != nil {
			client.Serializer = opts[0].Serializer
		}

		if opts[0].Deserializer != nil {
			client.Deserializer = opts[0].Deserializer
		}
	}

	return func(module core.Module) core.Module {
		clientModule := module.New(core.NewModuleOptions{})

		clientModule.NewProvider(core.ProviderOptions{
			Name:  CLIENT,
			Value: client,
		})

		clientModule.Export(CLIENT)
		return clientModule
	}
}

func (client *Client) Close() {
	client.Conn.Close()
}

func (client *Client) Send(event string, data interface{}) error {
	writer := bufio.NewWriter(client.Conn)

	message := Message{Event: event, Data: data}
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

func Inject(module core.Module) *Client {
	conn, ok := module.Ref(CLIENT).(*Client)
	if !ok {
		return nil
	}
	return conn
}
