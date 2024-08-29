package sql

import (
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `gorm:"embedded"`
	Name       string `gorm:"type:varchar(100)"`
}

func Test_Connect(t *testing.T) {
	t.Run("test case", func(t *testing.T) {
		ForFeature(&User{})
		dsn := "host=localhost user=postgres password=postgres dbname=tester port=5432 sslmode=disable TimeZone=Asia/Shanghai"
		ForRoot(postgres.Open(dsn))

		if db.conn == nil {
			t.Error("expect not nil, but got", db)
		}
	})
}
