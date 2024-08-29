package app

import (
	"github.com/tinh-tinh/tinhtinh/database/sql"
	"gorm.io/gorm"
)

type UserService interface {
	GetAll() []User
}

func NewService() UserService {
	return &UserServiceImpl{
		model: sql.GetModel[User](),
	}
}

type UserServiceImpl struct {
	model *gorm.DB
}

func (u *UserServiceImpl) GetAll() []User {
	var users []User
	u.model.First(&users)

	return users
}
