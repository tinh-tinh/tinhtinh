package redis

import (
	"context"
	"encoding/json"
	"fmt"

	redis_store "github.com/redis/go-redis/v9"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

type Client struct {
	Conn         *redis_store.Client
	Context      context.Context
	Serializer   core.Encode
	Deserializer core.Decode
}

func NewClient(opt microservices.ConnectOptions) microservices.ClientProxy {
	conn := redis_store.NewClient(&redis_store.Options{
		Addr:     opt.Addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	client := &Client{
		Context:      context.Background(),
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

func (client *Client) Close() {
	client.Conn.Close()
}

func (client *Client) Send(event string, data interface{}) error {
	payload, err := client.Serializer(data)
	if err != nil {
		return err
	}
	fmt.Printf("Send message: %v for event %s\n", data, event)
	err = client.Conn.Publish(client.Context, event, payload).Err()
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) Broadcast(data interface{}) error {
	return client.Send("*", data)
}
