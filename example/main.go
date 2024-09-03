package main

import (
	"github.com/tinh-tinh/tinhtinh/config"
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/dto/transform"
	"github.com/tinh-tinh/tinhtinh/example/app"
	"github.com/tinh-tinh/tinhtinh/swagger"
)

func main() {
	server := core.CreateFactory(app.NewModule, "api")

	document := swagger.NewSpecBuilder()
	document.SetHost("http://localhost:" + config.GetRaw("PORT")).SetBasePath("/")

	swagger.SetUp("docs", server, document)

	port := config.GetRaw("PORT")
	server.Listen(transform.StringToInt(port))
}
