package tcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"slices"
	"strings"

	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

type Server struct {
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

func (svc *Server) Create(module core.Module) {
	svc.Module = module
}

func (svc *Server) Listen() {
	listener, err := net.Listen("tcp", svc.Addr)
	if err != nil {
		panic(err)
	}
	store := svc.Module.Ref(microservices.STORE).(*microservices.Store)
	if store == nil {
		panic("store not found")
	}

	go http.Serve(listener, nil)
	for {
		conn, errr := listener.Accept()
		if errr != nil {
			panic(errr)
		}
		go svc.handler(conn, store)
	}
}

func (svc *Server) handler(conn net.Conn, store *microservices.Store) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message: ", err)
			return
		}

		var msg microservices.Message
		err = svc.deserializer([]byte(message), &msg)
		if err != nil {
			fmt.Println("Error deserializing message: ", err)
			return
		}
		if msg.Type == microservices.RPC {
			svc.handlerRPC(store.Subscribers[string(microservices.RPC)], msg)
		} else if msg.Type == microservices.PubSub {
			svc.handlerPubSub(store.Subscribers[string(microservices.PubSub)], msg)
		}
	}
}

func (svc *Server) handlerRPC(handlers []microservices.SubscribeHandler, msg microservices.Message) {
	data := microservices.ParseCtx(msg.Data, svc)
	subscriber := common.Filter(handlers, func(e microservices.SubscribeHandler) bool {
		return e.Name == msg.Event
	})
	for _, sub := range subscriber {
		sub.Factory(data)
	}
}

func (svc *Server) handlerPubSub(handlers []microservices.SubscribeHandler, msg microservices.Message) {
	data := microservices.ParseCtx(msg.Data, svc)
	if msg.Event == "*" {
		for _, provider := range handlers {
			provider.Factory(data)
		}
	} else if strings.ContainsAny(msg.Event, "*") {
		prefix := strings.TrimSuffix(msg.Event, "*")
		fmt.Println(prefix)
		for _, provider := range handlers {
			if strings.HasPrefix(string(provider.Name), prefix) {
				provider.Factory(data)
			}
		}
	} else {
		findEvent := slices.IndexFunc(handlers, func(e microservices.SubscribeHandler) bool {
			return string(e.Name) == msg.Event
		})
		if findEvent != -1 {
			provider := handlers[findEvent]
			provider.Factory(data)
		}
	}
}

func (svc *Server) Serializer(v interface{}) ([]byte, error) {
	return svc.serializer(v)
}

func (svc *Server) Deserializer(data []byte, v interface{}) error {
	return svc.deserializer(data, v)
}
