package main

import (
	"github.com/tinh-tinh/tinhtinh/config"
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/dto/transform"
	"github.com/tinh-tinh/tinhtinh/example/app"
	"github.com/tinh-tinh/tinhtinh/swagger"
)

func main() {
	server := core.CreateFactory(app.NewModule)
	server.SetGlobalPrefix("api")

	document := swagger.NewSpecBuilder().SetTitle("Swagger UI").SetDescription("Swagger Description").Build()
	swagger.SetUp(server, document)

	port := config.GetRaw("PORT")
	server.Listen(transform.StringToInt(port))
}
