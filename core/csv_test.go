package core_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func Test_ParseCsv(t *testing.T) {
	body := core.ParseCsv(nil, nil)
	require.Empty(t, body)
}

func Test_Csv(t *testing.T) {
	type User struct {
		UserID   string
		FullName string
		Email    string
	}

	controller := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx core.Ctx) error {
			user1 := &User{"1", "Jack Johnson", "jack@hotmail.com"}
			user2 := &User{"2", "Jill Smith", "jill@hotmail.com"}
			user3 := &User{"3", "James Murphy", "james@hotmail.com"}

			users := []*User{user1, user2, user3}

			data := core.ParseCsv(users, []string{"UserID", "FullName", "Email"})

			fmt.Println(data)
			return ctx.ExportCSV("users.csv", data)
		})

		ctrl.Get("error", func(ctx core.Ctx) error {
			data := core.ParseCsv(3, []string{"UserID", "FullName", "Email"})

			return ctx.ExportCSV("users.csv", data)
		})

		ctrl.Post("", func(ctx core.Ctx) error {
			return ctx.ExportCSV("data", nil)
		})

		return ctrl
	}

	module := func() *core.DynamicModule {
		appModule := core.NewModule(core.NewModuleOptions{
			Controllers: []core.Controllers{controller},
		})

		return appModule
	}

	app := core.CreateFactory(module)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()
	testClient := testServer.Client()

	resp, err := testClient.Get(testServer.URL + "/api/test")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	require.Equal(t, "UserID,FullName,Email\n1,Jack Johnson,jack@hotmail.com\n2,Jill Smith,jill@hotmail.com\n3,James Murphy,james@hotmail.com\n", string(data))

	resp, err = testClient.Post(testServer.URL+"/api/test", "application/json", nil)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/test/error")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err = io.ReadAll(resp.Body)
	require.Nil(t, err)

	require.Equal(t, "UserID,FullName,Email\n", string(data))

}
