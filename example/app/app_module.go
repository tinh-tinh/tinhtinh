package app

import (
	"github.com/tinh-tinh/tinhtinh/api"
	"github.com/tinh-tinh/tinhtinh/example/app/user"
)

func NewModule() *api.Module {
	appModule := api.NewModule()
	appModule.Import(user.NewModule())

	return appModule
}
