package post

import "github.com/tinh-tinh/tinhtinh/database/sql"

type Post struct {
	sql.Model `gorm:"embedded"`
	Title     string `gorm:"type:varchar(255)"`
	Content   string `gorm:"type:text"`
}
