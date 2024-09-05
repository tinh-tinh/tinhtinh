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
	port := config.GetRaw("PORT")

	document := swagger.NewSpecBuilder()
	document.SetHost("localhost:" + port).SetBasePath("/api")

	swagger.SetUp("docs", server, document)
	server.Listen(transform.StringToInt(port))
}
