package user

import (
	"github.com/tinh-tinh/tinhtinh/api"
	"github.com/tinh-tinh/tinhtinh/database/sql"
)

func NewModule() *api.Module {
	userModule := api.NewModule()

	sql.ForFeature(&User{})
	userModule.Controllers(NewController())

	return userModule
}
