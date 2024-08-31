package user

import (
	"github.com/tinh-tinh/tinhtinh/api"
	"github.com/tinh-tinh/tinhtinh/database/sql"
	"github.com/tinh-tinh/tinhtinh/example/app/user/dto"
	"gorm.io/gorm"
)

type Service interface {
	GetAll() []User
	Create(dto.SignUpUser) error
}

func service() api.Provider {
	return api.Provider{
		Name: "USER",
		Value: &ServiceImpl{
			model: sql.GetModel[User](),
		},
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
