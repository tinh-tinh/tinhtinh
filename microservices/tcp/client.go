package tcp

import (
	"errors"
	"net"
	"net/rpc"
	"time"

	"github.com/tinh-tinh/tinhtinh/microservices"
)

type Options struct {
	microservices.Config
	Addr    string
	Timeout time.Duration
}

type Client struct {
	config    microservices.Config
	eventConn net.Conn
	rpcClient *rpc.Client
	timeout   time.Duration
}

func connect(addr string) (net.Conn, *rpc.Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, nil, err
	}

	rpcClient, err := rpc.Dial("tcp", addr)
	if err != nil {
		return nil, nil, err
	}

	return conn, rpcClient, nil
}

func NewClient(opt Options) microservices.ClientProxy {
	eventConn, rpcClient, err := connect(opt.Addr)
	if err != nil {
		if opt.RetryOptions.Retry != 0 {
			time.Sleep(opt.RetryOptions.Delay)
			opt.RetryOptions.Retry--
			return NewClient(opt)
		}
		panic(err)
	}

	if opt.Timeout > 0 {
		eventConn.SetDeadline(time.Now().Add(opt.Timeout))
	}

	client := &Client{
		eventConn: eventConn,
		rpcClient: rpcClient,
		config:    microservices.NewConfig(opt.Config),
		timeout:   microservices.DEFAULT_TIMEOUT,
	}

	return client
}

func (client *Client) Config() microservices.Config {
	return client.config
}

func (client *Client) Publish(event string, data any, headers ...microservices.Header) error {
	payload, err := microservices.EncodeMessage(client, microservices.Message{
		Event:   event,
		Headers: microservices.AssignHeader(client.Config().Header, headers...),
		Data:    data,
	})
	if err != nil {
		client.config.ErrorHandler(err)
		return err
	}

	payload = append(payload, '\n')
	_, err = client.eventConn.Write(payload)
	if err != nil {
		client.config.ErrorHandler(err)
		return err
	}

	return nil
}

func (client *Client) Timeout(duration time.Duration) microservices.ClientProxy {
	client.timeout = duration
	err := client.eventConn.SetWriteDeadline(time.Now().Add(duration))
	if err != nil {
		panic(err)
	}
	return client
}

func (client *Client) Send(path string, data any, response any, headers ...microservices.Header) error {
	msg := microservices.Message{
		Event:   path,
		Headers: microservices.AssignHeader(client.Config().Header, headers...),
		Data:    data,
	}

	bytes, err := microservices.EncodeMessage(client, msg)
	if err != nil {
		client.config.ErrorHandler(err)
		return err
	}

	resRaw := []byte{}
	call := client.rpcClient.Go("RpcGateway.Call", &bytes, &resRaw, nil)
	select {
	case <-call.Done:
		if call.Error != nil {
			client.config.ErrorHandler(call.Error)
			return call.Error
		}
		err := client.config.Deserializer(resRaw, response)
		if err != nil {
			client.config.ErrorHandler(err)
			return err
		}
		return nil

	case <-time.After(client.timeout):
		err := errors.New("RPC call timed out")
		client.config.ErrorHandler(err)
		return err
	}
}
