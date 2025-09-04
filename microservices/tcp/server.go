package tcp

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"slices"
	"strings"

	"github.com/tinh-tinh/tinhtinh/microservices"
	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

type Server struct {
	Addr   string
	Module core.Module
	config microservices.Config
}

func New(module core.ModuleParam, opts ...Options) microservices.Service {
	svc := &Server{
		Module: module(),
		config: microservices.DefaultConfig(),
	}

	if len(opts) > 0 {
		if !reflect.ValueOf(opts[0].Config).IsZero() {
			svc.config = microservices.ParseConfig(opts[0].Config)
		}
		if opts[0].Addr != "" {
			svc.Addr = opts[0].Addr
		}
	}

	return svc
}

func Open(opts ...Options) core.Service {
	svc := &Server{
		config: microservices.DefaultConfig(),
	}

	if len(opts) > 0 {
		if !reflect.ValueOf(opts[0].Config).IsZero() {
			svc.config = microservices.ParseConfig(opts[0].Config)
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
	store, ok := svc.Module.Ref(microservices.STORE).(*microservices.Store)
	if !ok {
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

		msg := microservices.DecodeMessage(svc, []byte(message))
		switch msg.Type {
		case microservices.RPC:
			svc.handlerRPC(store.GetRPC(), msg)
		case microservices.PubSub:
			svc.handlerPubSub(store.GetPubSub(), msg)
		}
	}
}

func (svc *Server) handlerRPC(handlers []*microservices.SubscribeHandler, msg microservices.Message) {
	subscriber := common.Filter(handlers, func(e *microservices.SubscribeHandler) bool {
		return e.Name == msg.Event
	})
	for _, sub := range subscriber {
		sub.Handle(svc, msg)
	}
}

func (svc *Server) handlerPubSub(handlers []*microservices.SubscribeHandler, msg microservices.Message) {
	if msg.Event == "*" {
		for _, sub := range handlers {
			sub.Handle(svc, msg)
		}
	} else if strings.ContainsAny(msg.Event, "*") {
		prefix := strings.TrimSuffix(msg.Event, "*")
		for _, sub := range handlers {
			if strings.HasPrefix(string(sub.Name), prefix) {
				sub.Handle(svc, msg)
			}
		}
	} else {
		findEvent := slices.IndexFunc(handlers, func(e *microservices.SubscribeHandler) bool {
			return string(e.Name) == msg.Event
		})
		if findEvent != -1 {
			sub := handlers[findEvent]
			sub.Handle(svc, msg)
		}
	}
}

func (svc *Server) Config() microservices.Config {
	return svc.config
}
