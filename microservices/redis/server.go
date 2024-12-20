package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	redis_store "github.com/redis/go-redis/v9"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

type Server struct {
	Context      context.Context
	Addr         string
	Module       core.Module
	Serializer   core.Encode
	Deserializer core.Decode
}

func New(module core.ModuleParam, opts ...microservices.ConnectOptions) microservices.Service {
	svc := Open(opts...)
	svc.Create(module())

	return svc
}

func Open(opts ...microservices.ConnectOptions) core.Service {
	svc := &Server{
		Serializer:   json.Marshal,
		Deserializer: json.Unmarshal,
		Context:      context.Background(),
	}

	if len(opts) > 0 {
		if opts[0].Serializer != nil {
			svc.Serializer = opts[0].Serializer
		}

		if opts[0].Deserializer != nil {
			svc.Deserializer = opts[0].Deserializer
		}
	}

	return svc
}

func (svc *Server) Create(module core.Module) {
	svc.Module = module
}

func (svc *Server) Listen() {
	rdb := redis_store.NewClient(&redis_store.Options{
		Addr:     svc.Addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()
	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	for _, prd := range svc.Module.GetDataProviders() {
		var subscriber *redis_store.PubSub
		if strings.HasSuffix(string(prd.GetName()), "*") {
			subscriber = rdb.PSubscribe(svc.Context, string(prd.GetName()))
		} else {
			subscriber = rdb.Subscribe(svc.Context, string(prd.GetName()))
		}
		go svc.Handler(subscriber, prd.GetFactory())
	}
}

func (svc *Server) Handler(params ...interface{}) {
	subscriber := params[0].(*redis_store.PubSub)
	factory := params[1].(core.Factory)
	for {
		msg, err := subscriber.ReceiveMessage(svc.Context)
		if err != nil {
			fmt.Println("Error reading message: ", err)
			return
		}

		fmt.Printf("Received message: %s from event %s\n", msg.Payload, msg.Channel)
		data := microservices.ParseCtx(msg.Payload)
		factory(data)
	}
}
