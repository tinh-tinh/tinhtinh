package core

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Metadata(t *testing.T) {
	const role_key = "roles"

	roleFnc := func(roles ...string) *Metadata {
		return SetMetadata(role_key, roles)
	}

	controller := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test").Guard(func(ctrl *DynamicController, ctx Ctx) bool {
			roles, ok := ctrl.GetMetadata(role_key).([]string)
			if !ok {
				return false
			}
			isRole := slices.IndexFunc(roles, func(role string) bool {
				return ctx.Query("role") == role
			})
			return isRole != -1
		})

		ctrl.Metadata(roleFnc("admin")).Get("", func(ctx Ctx) {
			ctx.JSON(Map{
				"data": "ok",
			})
		})

		return ctrl
	}

	module := func() *DynamicModule {
		mod := NewModule(NewModuleOptions{
			Controllers: []Controller{controller},
		})

		return mod
	}

	app := CreateFactory(module)
	app.SetGlobalPrefix("api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()
	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test?role=admin")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
