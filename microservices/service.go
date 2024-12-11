package microservices

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"slices"

	"github.com/tinh-tinh/tinhtinh/v2/core"
)

type Service struct {
	Conn         net.Listener
	Module       core.Module
	Serializer   core.Encode
	Deserializer core.Decode
	RetryAttempt int
}

type Options struct {
	Wildcard     bool
	Delimiter    string
	Serializer   core.Encode
	Deserializer core.Decode
	RetryAttempt int
}

func New(module core.ModuleParam, opts ...Options) *Service {
	svc := &Service{
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

		svc.RetryAttempt = opts[0].RetryAttempt
	}

	return svc
}

func (svc *Service) Listen(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	svc.Conn = listener

	go http.Serve(svc.Conn, nil)
	for {
		conn, errr := svc.Conn.Accept()
		if errr != nil {
			panic(errr)
		}
		go svc.Handler(conn)
	}
}

func (svc *Service) Handler(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message: ", err)
			return
		}

		var msg Message
		err = svc.Deserializer([]byte(message), &msg)
		if err != nil {
			fmt.Println("Error deserializing message: ", err)
			return
		}
		fmt.Printf("Received message: %v from event %v\n", msg.Data, msg.Event)

		findEvent := slices.IndexFunc(svc.Module.GetDataProviders(), func(e core.Provider) bool {
			return string(e.GetName()) == msg.Event
		})
		if findEvent != -1 {
			provider := svc.Module.GetDataProviders()[findEvent]
			fnc := provider.GetFactory()
			fnc(msg.Data)
		}
	}
}
