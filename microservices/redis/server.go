package redis

import (
	"context"
	"encoding/json"
	"strings"

	redis_store "github.com/redis/go-redis/v9"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

type Server struct {
	Context      context.Context
	Addr         string
	Module       core.Module
	serializer   core.Encode
	deserializer core.Decode
}

func New(module core.ModuleParam, opts ...microservices.ConnectOptions) microservices.Service {
	svc := &Server{
		Module:       module(),
		serializer:   json.Marshal,
		deserializer: json.Unmarshal,
	}

	if len(opts) > 0 {
		if opts[0].Serializer != nil {
			svc.serializer = opts[0].Serializer
		}

		if opts[0].Deserializer != nil {
			svc.deserializer = opts[0].Deserializer
		}
		if opts[0].Addr != "" {
			svc.Addr = opts[0].Addr
		}
	}

	return svc
}

func Open(opts ...microservices.ConnectOptions) core.Service {
	svc := &Server{
		serializer:   json.Marshal,
		deserializer: json.Unmarshal,
		Context:      context.Background(),
	}

	if len(opts) > 0 {
		if opts[0].Serializer != nil {
			svc.serializer = opts[0].Serializer
		}

		if opts[0].Deserializer != nil {
			svc.deserializer = opts[0].Deserializer
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
		return
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
			return
		}

		data := microservices.ParseCtx(msg.Payload, svc)
		factory(data)
	}
}

func (svc *Server) Serializer(v interface{}) ([]byte, error) {
	return svc.serializer(v)
}

func (svc *Server) Deserializer(data []byte, v interface{}) error {
	return svc.deserializer(data, v)
}
