package token

import (
	"strings"

	"github.com/tinh-tinh/tinhtinh/core"
)

const USER core.CtxKey = "USER"

func Guard(ctrl *core.DynamicController, ctx core.Ctx) bool {
	tokenService := ctrl.Inject(TOKEN).(Provider)
	header := ctx.Headers("Authorization")
	if header == "" {
		return false
	}
	token := strings.Split(header, " ")[1]

	payload, err := tokenService.Verify(token)
	if err != nil {
		return false
	}

	ctx.Set(USER, payload)
	return true
}
