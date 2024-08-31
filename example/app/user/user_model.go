package user

import "github.com/tinh-tinh/tinhtinh/database/sql"

type User struct {
	sql.Model `gorm:"embedded"`
	Name      string `gorm:"type:varchar(100)"`
	Email     string `gorm:"type:varchar(100);uniqueIndex:idx_email"`
	Password  string `gorm:"type:varchar(255)"`
	Role      string `gorm:"type:varchar(100)"`
	Active    bool   `gorm:"type:bool"`
}
