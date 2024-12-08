package main

import (
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
	app := microservices.Start(appModule)
	app.Listen("localhost:8080")
}
