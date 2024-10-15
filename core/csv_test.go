package core

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Csv(t *testing.T) {
	type User struct {
		UserID   string
		FullName string
		Email    string
	}

	controller := func(module *DynamicModule) *DynamicController {
		ctrl := module.NewController("test")

		ctrl.Get("", func(ctx Ctx) error {
			user1 := &User{"1", "Jack Johnson", "jack@hotmail.com"}
			user2 := &User{"2", "Jill Smith", "jill@hotmail.com"}
			user3 := &User{"3", "James Murphy", "james@hotmail.com"}

			users := []*User{user1, user2, user3}

			data := ParseCsv(users, []string{"UserID", "FullName", "Email"})

			fmt.Println(data)
			return ctx.ExportCSV("users.csv", data)
		})

		return ctrl
	}

	module := func() *DynamicModule {
		appModule := NewModule(NewModuleOptions{
			Controllers: []Controller{controller},
		})

		return appModule
	}

	app := CreateFactory(module)
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
}
