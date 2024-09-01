package main

import (
	"github.com/tinh-tinh/tinhtinh/config"
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/dto/transform"
	"github.com/tinh-tinh/tinhtinh/example/app"
)

func main() {
	server := core.CreateFactory(app.NewModule)
	server.SetGlobalPrefix("api")

	port := config.GetRaw("PORT")
	server.Listen(transform.StringToInt(port))
}
