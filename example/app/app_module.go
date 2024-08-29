package app

import (
	"github.com/tinh-tinh/tinhtinh/api"
	"github.com/tinh-tinh/tinhtinh/database/sql"
)

func NewModule() *api.Module {
	sql.ForFeature(&User{})

	appModule := api.NewModule()
	appModule.Controllers(NewController())

	return appModule
}
