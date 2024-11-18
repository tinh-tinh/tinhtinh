package validator

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type UserDetail struct {
	ID       string    `validate:"isObjectId"`
	Name     string    `validate:"required,isAlpha"`
	Username string    `validate:"isAlphaNumeric"`
	Email    string    `validate:"required,isEmail"`
	Password string    `validate:"isStrongPassword"`
	UserID   string    `validate:"isUUID"`
	Point    string    `validate:"isFloat"`
	Height   float64   `validate:"isFloat"`
	Age      int       `validate:"isInt"`
	Active   bool      `validate:"isBool"`
	IsAdmin  string    `validate:"isBool"`
	Birth    string    `validate:"isDateString"`
	Total    int       `validate:"isNumber"`
	Join     time.Time `validate:"isDate"`
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

	userDetails3 := &UserDetail{
		Name:     "haha",
		Email:    "babaddok@gmail.com",
		Password: "12345678@Tc",
		UserID:   "1234",
		Point:    "true",
		ID:       "1234",
		Username: "##$$",
		Age:      1,
		Active:   true,
		Birth:    "true",
		Height:   1.5,
		Total:    1,
		IsAdmin:  "abc",
		Join:     time.Now(),
	}
	user3 := &User{
		ID:     3,
		Detail: userDetails3,
	}
	err5 := Scanner(user3)
	require.NotNil(t, err5)
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

func Test_Scanner2(t *testing.T) {
	t.Parallel()

	require.Panics(t, func() { Scanner(UserDetail{}) })
	require.Panics(t, func() { Scanner(nil) })
}

func TestDefault(t *testing.T) {
	type Pagination struct {
		Page  int `validate:"isInt" default:"1"`
		Limit int `validate:"isInt" default:"10"`
	}

	pagin := &Pagination{}
	require.Nil(t, Scanner(pagin))
	require.Equal(t, 1, pagin.Page)
	require.Equal(t, 10, pagin.Limit)

	pagin2 := &Pagination{Page: 4, Limit: 20}
	require.Nil(t, Scanner(pagin2))
	require.Equal(t, 4, pagin2.Page)
	require.Equal(t, 20, pagin2.Limit)
}
