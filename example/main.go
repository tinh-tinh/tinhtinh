package main

import (
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/example/app"
	"github.com/tinh-tinh/tinhtinh/utils"
)

func main() {
	utils.PrintAlloc()
	server := core.CreateFactory(app.NewModule, "api")
	// port := config.GetRaw("PORT")

	// document := swagger.NewSpecBuilder()
	// document.SetHost("localhost:" + port).SetBasePath("/api").AddSecurity(&swagger.SecuritySchemeObject{
	// 	Type: "apiKey",
	// 	In:   "header",
	// 	Name: "Authorization",
	// })

	// swagger.SetUp("docs", server, document)
	utils.PrintAlloc()
	server.Listen(3000)
}
