package tcp

import (
	"net"
	"time"

	"github.com/tinh-tinh/tinhtinh/microservices"
)

type RetryOptions struct {
	Retry int
	Delay time.Duration
}

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
		if opt.RetryOptions.Retry != 0 {
			time.Sleep(opt.RetryOptions.Delay)
			opt.RetryOptions.Retry--
			return NewClient(opt)
		}
		panic(err)
	}

	if opt.Timeout > 0 {
		conn.SetDeadline(time.Now().Add(opt.Timeout))
	}

	client := &Client{
		Conn:   conn,
		config: microservices.NewConfig(opt.Config),
	}

	return client
}

func (client *Client) Config() microservices.Config {
	return client.config
}

func (client *Client) Send(event string, data interface{}, headers ...microservices.Header) error {
	return microservices.DefaultSend(client)(event, data, headers...)
}

func (client *Client) Publish(event string, data interface{}, headers ...microservices.Header) error {
	return microservices.DefaultPublish(client)(event, data, headers...)
}

func (client *Client) Timeout(duration time.Duration) microservices.ClientProxy {
	err := client.Conn.SetWriteDeadline(time.Now().Add(duration))
	if err != nil {
		panic(err)
	}
	return client
}

func (client *Client) Emit(event string, message microservices.Message) error {
	payload, err := microservices.EncodeMessage(client, message)
	if err != nil {
		client.config.ErrorHandler(err)
		return err
	}

	payload = append(payload, '\n')
	_, err = client.Conn.Write(payload)
	if err != nil {
		client.config.ErrorHandler(err)
		return err
	}

	return nil
}
