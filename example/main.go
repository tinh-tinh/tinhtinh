package main

import (
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/example/app"
)

func main() {
	server := core.CreateFactory(app.NewModule, "api")
	// port := config.GetRaw("PORT")

	// document := swagger.NewSpecBuilder()
	// document.SetHost("localhost:" + port).SetBasePath("/api").AddSecurity(&swagger.SecuritySchemeObject{
	// 	Type: "apiKey",
	// 	In:   "header",
	// 	Name: "Authorization",
	// })

	// swagger.SetUp("docs", server, document)
	server.Listen(3000)
}
