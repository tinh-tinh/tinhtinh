package main

import (
	"fmt"

	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/microservices"
)

type Message struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	appModule := func() *core.DynamicModule {
		module := core.NewModule(core.NewModuleOptions{})
		return module
	}
	app := microservices.StartMicroservice(appModule)
	app.Listen("localhost:8080")
	fmt.Println(app)
	defer app.Conn.Close()

	for {
		conn, err := app.Conn.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go microservices.HandlerConnect(conn)
	}
}
