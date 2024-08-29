package sql

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Model struct {
	ID        *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt *time.Time     `gorm:"not null;default:now()"`
	UpdatedAt *time.Time     `gorm:"not null;default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func GetModel[M any]() *gorm.DB {
	var model M
	return db.conn.Model(&model)
}
