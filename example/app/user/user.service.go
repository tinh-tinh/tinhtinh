package user

import (
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/database/sql"
	"github.com/tinh-tinh/tinhtinh/example/app/user/dto"
	"gorm.io/gorm"
)

type CrudService struct {
	model *gorm.DB
}

func (s *CrudService) GetAll() []User {
	var user []User
	s.model.First(&user)

	return user
}

func (s *CrudService) Create(input dto.SignUpUser) error {
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

const USER_SERVICE core.Provide = "UserService"

func service(module *core.DynamicModule) *core.DynamicProvider {
	userSv := core.NewProvider(module)

	userSv.Set(USER_SERVICE, &CrudService{
		model: module.Ref(sql.ConnectDB).(*gorm.DB),
	})

	return userSv
}
