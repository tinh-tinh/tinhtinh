package main

import (
	"encoding/json"
	"fmt"

	"github.com/tinh-tinh/tinhtinh/microservices"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

type Message struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	appService := func(module core.Module) core.Provider {
		handler := microservices.NewHandler(module, core.ProviderOptions{})

		handler.OnResponse("user.created", func(param ...interface{}) interface{} {
			if len(param) == 0 {
				return nil
			}
			msg := param[0]
			var decodedData Message
			if msg != nil {
				dataBytes, _ := json.Marshal(msg)
				_ = json.Unmarshal(dataBytes, &decodedData)
				fmt.Println("Decoded Data:", decodedData)
			}

			return nil
		})

		handler.OnResponse("user.updated", func(param ...interface{}) interface{} {
			if len(param) == 0 {
				return nil
			}
			msg := param[0]
			var decodedData Message
			if msg != nil {
				dataBytes, _ := json.Marshal(msg)
				_ = json.Unmarshal(dataBytes, &decodedData)
				fmt.Println("Decoded Data:", decodedData)
			}

			return nil
		})

		return handler
	}

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Providers: []core.Providers{
				appService,
			},
		})
		return module
	}
	app := microservices.New(appModule)
	app.Listen("localhost:8080")
}
