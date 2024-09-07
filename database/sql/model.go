package sql

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	CreatedAt *time.Time     `gorm:"not null;default:now()"`
	UpdatedAt *time.Time     `gorm:"not null;default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
