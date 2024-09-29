package validator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserDetail struct {
	Name     string `validate:"required,isAlpha"`
	Email    string `validate:"required,isEmail"`
	Password string `validate:"isStrongPassword"`
}

type User struct {
	ID     int         `validate:"required,isInt"`
	Detail *UserDetail `validate:"nested"`
	Banner string      `validate:"isAlpha"`
}

type UserList struct {
	Users []*User `validate:"nested"`
}

type Post struct {
	Author *User  `validate:"nested"`
	Title  string `validate:"required"`
}

func Test_Scanner(t *testing.T) {
	t.Parallel()

	userDetails1 := &UserDetail{
		Email:    "tinh@gmail.com",
		Password: "12345678@Tc",
	}
	user1 := &User{
		ID:     1,
		Detail: userDetails1,
		Banner: "true",
	}

	err1 := Scanner(user1)
	require.NotNil(t, err1)
	require.Equal(t, "Name is required\nName is not a valid alpha", err1.Error())

	userDetails2 := &UserDetail{
		Name:     "haha",
		Email:    "babaddok@gmail.com",
		Password: "12345678@Tc",
	}
	user2 := &User{
		ID:     2,
		Detail: userDetails2,
	}

	err2 := Scanner(user2)
	require.Nil(t, err2)

	userList := &UserList{
		Users: []*User{user1, user2},
	}
	err3 := Scanner(userList)
	require.NotNil(t, err3)
	require.Equal(t, "Name is required\nName is not a valid alpha", err3.Error())

	post := &Post{
		Author: user2,
		Title:  "",
	}
	err4 := Scanner(post)
	require.NotNil(t, err4)
	require.Equal(t, "Title is required", err4.Error())
}

func Benchmark_Scanner(b *testing.B) {
	b.Run("test_validator", func(b *testing.B) {
		var userList []*User
		count := b.N
		fmt.Printf("No of test case %d\n", count)
		for n := 0; n < count; n++ {
			userDetail := &UserDetail{
				Name:     "Haha",
				Email:    fmt.Sprintf("abc%d@gmail.com", n),
				Password: fmt.Sprintf("1234567%d@Abc", n),
			}
			user := &User{
				ID:     n,
				Detail: userDetail,
			}
			require.Nil(b, Scanner(user))
			userList = append(userList, user)
		}

		require.Nil(b, Scanner(&UserList{
			Users: userList,
		}))
	})
}
