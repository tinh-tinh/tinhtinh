package microservices

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/tinh-tinh/tinhtinh/core"
)

type Service struct {
	Conn         net.Listener
	Module       *core.DynamicModule
	Serializer   core.Encode
	Deserializer core.Decode
	RetryAttempt int
	handlers     map[string]HandlerFnc
}

type Options struct {
	Serializer   core.Encode
	Deserializer core.Decode
	RetryAttempt int
}

func Start(module core.ModuleParam, options ...Options) *Service {
	svc := &Service{
		Module:       module(),
		Serializer:   json.Marshal,
		Deserializer: json.Unmarshal,
		handlers:     make(map[string]HandlerFnc),
	}

	if len(options) > 0 {
		if options[0].Serializer != nil {
			svc.Serializer = options[0].Serializer
		}
		if options[0].Deserializer != nil {
			svc.Deserializer = options[0].Deserializer
		}
		if options[0].RetryAttempt > 0 {
			svc.RetryAttempt = options[0].RetryAttempt
		}
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
		conn, err := svc.Conn.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go svc.HandlerConnect(conn)
	}
}

func (svc *Service) HandlerConnect(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}

		var msg Package
		err = json.Unmarshal([]byte(message), &msg)
		if err != nil {
			fmt.Println("Error unmarshalling message:", err)
			return
		}

		if svc.handlers[msg.Event] != nil {
			svc.handlers[msg.Event](msg.Data)
		}

		fmt.Printf("Received message: %v from event %v\n", msg.Data, msg.Event)
	}
}
