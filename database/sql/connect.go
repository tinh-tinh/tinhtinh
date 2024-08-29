package sql

import (
	"fmt"

	"gorm.io/gorm"
)

type DB struct {
	conn   *gorm.DB
	models []interface{}
}

var db DB

func ForRoot(dialect gorm.Dialector) {
	conn, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println("connected to database migrating...")
	err = conn.AutoMigrate(db.models...)
	if err != nil {
		panic(err)
	}
	db.conn = conn
}

func ForFeature(models ...interface{}) {
	db.models = append(db.models, models...)
}
