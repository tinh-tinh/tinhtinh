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
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

type Server struct {
	Addr   string
	Module core.Module
	config microservices.Config
}

func New(module core.Module, opts ...Options) microservices.Service {
	svc := &Server{
		Module: module,
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

	var subscribers []*microservices.SubscribeHandler
	store, ok := svc.Module.Ref(microservices.STORE).(*microservices.Store)
	if ok && store != nil {
		subscribers = append(subscribers, store.Subscribers...)
	}
	tcpStore, ok := svc.Module.Ref(microservices.ToTransport(microservices.TCP)).(*microservices.Store)
	if ok && tcpStore != nil {
		subscribers = append(subscribers, tcpStore.Subscribers...)
	}

	if store == nil && tcpStore == nil {
		panic("store required")
	}

	go http.Serve(listener, nil)
	for {
		conn, errr := listener.Accept()
		if errr != nil {
			panic(errr)
		}
		go svc.handler(conn, subscribers)
	}
}

func (svc *Server) handler(conn net.Conn, subscribers []*microservices.SubscribeHandler) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message: ", err)
			return
		}

		msg := microservices.DecodeMessage(svc, []byte(message))
		svc.handlerPubSub(subscribers, msg)
	}
}

func (svc *Server) handlerPubSub(handlers []*microservices.SubscribeHandler, msg microservices.Message) {
	if msg.Event == "*" {
		for _, sub := range handlers {
			go sub.Handle(svc, msg)
		}
	} else if strings.ContainsAny(msg.Event, "*") {
		prefix := strings.TrimSuffix(msg.Event, "*")
		for _, sub := range handlers {
			if strings.HasPrefix(string(sub.Name), prefix) {
				go sub.Handle(svc, msg)
			}
		}
	} else {
		findEvent := slices.IndexFunc(handlers, func(e *microservices.SubscribeHandler) bool {
			return string(e.Name) == msg.Event
		})
		if findEvent != -1 {
			sub := handlers[findEvent]
			go sub.Handle(svc, msg)
		}
	}
}

func (svc *Server) Config() microservices.Config {
	return svc.config
}
