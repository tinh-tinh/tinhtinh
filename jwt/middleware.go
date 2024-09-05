package jwt

import (
	"strings"

	"github.com/tinh-tinh/tinhtinh/core"
)

const USER core.Provide = "USER"

func Guard(ctrl *core.DynamicController, ctx core.Ctx) bool {
	jwtService := ctrl.Inject(JWT).(Service)
	header := ctx.Headers("Authorization")
	if header == "" {
		return false
	}
	token := strings.Split(header, " ")[1]

	payload, err := jwtService.VerifyToken(token)
	if err != nil {
		return false
	}

	ctx.Set(USER, payload)
	return true
}
