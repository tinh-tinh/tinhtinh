package main

import (
	"fmt"

	"github.com/tinh-tinh/tinhtinh/api"
	"github.com/tinh-tinh/tinhtinh/config"
	"github.com/tinh-tinh/tinhtinh/database/sql"
	"github.com/tinh-tinh/tinhtinh/example/app"
	"gorm.io/driver/postgres"
)

type Config struct {
	Port    int    `mapstructure:"PORT"`
	NodeEnv string `mapstructure:"NODE_ENV"`

	DBHost string `mapstructure:"DB_HOST"`
	DBPort int    `mapstructure:"DB_PORT"`
	DBUser string `mapstructure:"DB_USER"`
	DBPass string `mapstructure:"DB_PASS"`
	DBName string `mapstructure:"DB_NAME"`
}

func init() {
	config.Register[Config]("")
	cfg := config.Get[Config]()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort)
	// "host=localhost user=postgres password=postgres dbname=tester port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	sql.ForRoot(postgres.Open(dsn))
}

func main() {
	server := api.New(app.NewModule())
	server.SetGlobalPrefix("api")

	server.Listen(config.Get[Config]().Port)
}
