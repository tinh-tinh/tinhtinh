package user

import (
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/example/app/user/dto"
	"gorm.io/gorm"
)

const USER_SERVICE core.Provide = "UserService"

type Service interface {
	GetAll() []User
	Create(dto.SignUpUser) error
}

func userService(module *core.DynamicModule) *core.DynamicProvider {
	pd := core.NewProvider(USER_SERVICE, &ServiceImpl{
		model: module.Ref("user").(*gorm.DB),
	})

	return pd
}

type ServiceImpl struct {
	model *gorm.DB
}

func (s *ServiceImpl) GetAll() []User {
	var user []User
	s.model.First(&user)

	return user
}

func (s *ServiceImpl) Create(input dto.SignUpUser) error {
	newUser := User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Role:     "user",
		Active:   true,
	}

	result := s.model.Create(&newUser)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
