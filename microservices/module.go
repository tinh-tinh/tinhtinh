package microservices

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"

	"github.com/tinh-tinh/tinhtinh/core"
)

type ClientOptions struct {
	Addr string
}

const CLIENT core.Provide = "client"

type Client struct {
	Conn net.Conn
}

func RegisterClient(opt ClientOptions) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		conn, err := net.Dial("tcp", opt.Addr)
		if err != nil {
			panic(err)
		}
		client := &Client{Conn: conn}
		clientModule := module.New(core.NewModuleOptions{})

		clientModule.NewProvider(core.ProviderOptions{Name: CLIENT, Value: client})
		clientModule.Export(CLIENT)

		return clientModule
	}
}

func (client *Client) Close() {
	client.Conn.Close()
}

func (client *Client) Send(event string, data interface{}) error {
	writer := bufio.NewWriter(client.Conn)

	message := Package{Event: event, Data: data}
	jsonData, err := json.Marshal(message)
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

func Inject(module *core.DynamicModule) *Client {
	conn, ok := module.Ref(CLIENT).(*Client)
	if !ok {
		return nil
	}
	return conn
}
