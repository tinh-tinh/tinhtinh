package user

import (
	"github.com/tinh-tinh/tinhtinh/database/sql"
	"gorm.io/gorm"
)

type Service interface {
	GetAll() []User
}

func NewService() Service {
	return &ServiceImpl{
		model: sql.GetModel[User](),
	}
}

type ServiceImpl struct {
	model *gorm.DB
}

func (s *ServiceImpl) GetAll() []User {
	var user []User
	s.model.First(&user)

	return user
}
