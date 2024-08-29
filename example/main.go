package main

import (
	"github.com/tinh-tinh/tinhtinh/api"
	"github.com/tinh-tinh/tinhtinh/database/sql"
	"github.com/tinh-tinh/tinhtinh/example/app"
	"gorm.io/driver/postgres"
)

func init() {
	dsn := "host=localhost user=postgres password=postgres dbname=tester port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	sql.ForRoot(postgres.Open(dsn))
}

func main() {
	server := api.New(app.NewModule())
	server.SetGlobalPrefix("api")

	server.Listen(3000)
}
