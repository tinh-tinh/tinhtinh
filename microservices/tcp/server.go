package tcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"slices"
	"strings"

	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

type Server struct {
	Addr         string
	Module       core.Module
	Serializer   core.Encode
	Deserializer core.Decode
}

// New creates a new TCP server with the given module and options.
//
// The module is the module that the server uses to initialize itself.
// The options are the options that the server uses to initialize itself.
// The options can be used to override the default encoder, decoder, and address.
//
// The server is created with a default encoder and decoder of json.Marshal and
// json.Unmarshal respectively.
// The server is created with a default address of ":8080".
// The server is initialized by calling the init method of the module.
// The server is then returned.
func New(module core.ModuleParam, opts ...microservices.ConnectOptions) microservices.Service {
	svc := &Server{
		Module:       module(),
		Serializer:   json.Marshal,
		Deserializer: json.Unmarshal,
	}

	if len(opts) > 0 {
		if opts[0].Serializer != nil {
			svc.Serializer = opts[0].Serializer
		}

		if opts[0].Deserializer != nil {
			svc.Deserializer = opts[0].Deserializer
		}
		if opts[0].Addr != "" {
			svc.Addr = opts[0].Addr
		}
	}

	return svc
}

func Open(opts ...microservices.ConnectOptions) core.Service {
	svc := &Server{
		Serializer:   json.Marshal,
		Deserializer: json.Unmarshal,
	}

	if len(opts) > 0 {
		if opts[0].Serializer != nil {
			svc.Serializer = opts[0].Serializer
		}

		if opts[0].Deserializer != nil {
			svc.Deserializer = opts[0].Deserializer
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

// Listen starts the TCP server on the address specified by the Addr field.
//
// It will listen for incoming connections and handle each connection in a
// separate goroutine. The handler function will be called for each incoming
// connection. The handler function will receive the connection as a parameter.
//
// The server will panic if there is an error when listening for connections or
// if there is an error when accepting a connection.
func (svc *Server) Listen() {
	listener, err := net.Listen("tcp", svc.Addr)
	if err != nil {
		panic(err)
	}

	go http.Serve(listener, nil)
	for {
		conn, errr := listener.Accept()
		if errr != nil {
			panic(errr)
		}
		go svc.handler(conn)
	}
}

// handler processes an incoming TCP connection represented by the param.
// It reads messages from the connection, deserializes them into a microservices.Message
// structure, and finds the corresponding event provider in the module's data providers.
// If the event provider is found, it executes the provider's factory function with the
// message data. The function handles errors in reading and deserializing the message
// by printing error messages and terminating the connection handling.

func (svc *Server) handler(param interface{}) {
	conn := param.(net.Conn)
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message: ", err)
			return
		}

		var msg microservices.Message
		err = svc.Deserializer([]byte(message), &msg)
		if err != nil {
			fmt.Println("Error deserializing message: ", err)
			return
		}
		fmt.Printf("Received message: %v from event %v\n", msg.Data, msg.Event)

		data := microservices.ParseCtx(msg.Data)
		if msg.Event == "*" {
			for _, provider := range svc.Module.GetDataProviders() {
				fnc := provider.GetFactory()
				fnc(data)
			}
		} else if strings.ContainsAny(msg.Event, "*") {
			prefix := strings.TrimSuffix(msg.Event, "*")
			fmt.Println(prefix)
			for _, provider := range svc.Module.GetDataProviders() {
				if strings.HasPrefix(string(provider.GetName()), prefix) {
					fnc := provider.GetFactory()
					fnc(data)
				}
			}
		} else {
			findEvent := slices.IndexFunc(svc.Module.GetDataProviders(), func(e core.Provider) bool {
				return string(e.GetName()) == msg.Event
			})
			if findEvent != -1 {
				provider := svc.Module.GetDataProviders()[findEvent]
				fnc := provider.GetFactory()
				fnc(data)
			}
		}
	}
}
