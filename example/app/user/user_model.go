package user

import "github.com/tinh-tinh/tinhtinh/database/sql"

type User struct {
	sql.Model `gorm:"embedded"`
	Name      string `gorm:"type:varchar(100)"`
}
