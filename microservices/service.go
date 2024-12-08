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
	Conn net.Listener
}

type Options struct {
	Addr string
}

func Create(opt Options) *Service {
	listener, err := net.Listen("tcp", opt.Addr)
	if err != nil {
		panic(err)
	}
	return &Service{Conn: listener}
}

func HandlerConnect(conn net.Conn) {
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

		fmt.Printf("Received message: %v from event %v\n", msg.Data, msg.Event)
	}
}

func StartMicroservice(module core.ModuleParam) *HybridApp {
	app := &HybridApp{
		App: core.App{
			Module: module(),
			Mux:    http.NewServeMux(),
		},
	}

	return app
}

type HybridApp struct {
	core.App
	Conn net.Listener
}

func (hybrid *HybridApp) Listen(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	hybrid.Conn = listener

	http.Handle("/", hybrid.PrepareBeforeListen())
	go http.Serve(hybrid.Conn, nil)
}
