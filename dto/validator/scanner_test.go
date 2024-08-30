package validator

import (
	"testing"
)

type User struct {
	Name     string `validate:"required,isAlpha"`
	Email    string `validate:"required,isEmail"`
	Password string `validate:"isStrongPassword"`
}

func Test_Scanner(t *testing.T) {
	testcases := []struct {
		input User
		want  string
	}{
		{input: User{
			Email:    "tinh@gmail.com",
			Password: "12345678@Tc",
		}, want: "Name is not a valid alpha"},
		{input: User{
			Name:     "Baba",
			Email:    "babaddok@gmail.com",
			Password: "12345678@Tc",
		}, want: ""},
	}

	for _, tc := range testcases {
		t.Run("test case", func(t *testing.T) {
			errMsg := Scanner(&tc.input)

			if errMsg.Error() != tc.want {
				t.Errorf("expect %s, but got %s", tc.want, errMsg.Error())
			}
		})
	}
}
