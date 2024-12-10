package validator_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/dto/validator"
)

func Test_Scanner(t *testing.T) {
	require.Panics(t, func() {
		_ = validator.Scanner(nil)
	})
	type Nested struct {
		IsAlpha string `validate:"isAlpha"`
	}
	type Input struct {
		Required         string    `validate:"required"`
		IsAlpha          string    `validate:"isAlpha"`
		IsAlphanumeric   string    `validate:"isAlphaNumeric"`
		IsEmail          string    `validate:"isEmail"`
		IsStrongPassword string    `validate:"isStrongPassword"`
		IsUUID           string    `validate:"isUUID"`
		IsFloat          float64   `validate:"isFloat"`
		IsInt            int       `validate:"isInt"`
		IsBool           bool      `validate:"isBool"`
		IsDateString     time.Time `validate:"isDateString"`
		IsNumber         int       `validate:"isNumber"`
		IsNumber2        float64   `validate:"isNumber"`
		IsObjectId       string    `validate:"isObjectId"`
		Nested           *Nested   `validate:"nested"`
		Slice            []*Nested `validate:"nested"`
		Lala             string    `validate:"isAlpha"`
	}
	require.Panics(t, func() {
		_ = validator.Scanner(Input{})
	})

	happyCase := &Input{
		Required:         "required",
		IsAlpha:          "abc",
		IsAlphanumeric:   "abc123",
		IsEmail:          "0K9oE@example.com",
		IsStrongPassword: "12345678@Tc",
		IsUUID:           "550e8400-e29b-41d4-a716-446655440000",
		IsFloat:          123.123,
		IsInt:            123,
		IsBool:           true,
		IsDateString:     time.Now(),
		IsNumber:         123,
		IsNumber2:        39.49,
		IsObjectId:       "5e9bf1f6d3d2d3d3d3d3d3d3",
		Nested: &Nested{
			IsAlpha: "abc",
		},
		Slice: []*Nested{
			{IsAlpha: "avc"},
		},
	}

	err := validator.Scanner(happyCase)
	require.Nil(t, err)

	badCaseStr := &Input{
		IsAlpha:          "$#%",
		IsAlphanumeric:   "#$%#^^%$",
		IsEmail:          "abc",
		IsStrongPassword: "mno",
		IsUUID:           "abc",
		IsObjectId:       "fvddf",
		Nested:           &Nested{IsAlpha: "757557"},
		Slice: []*Nested{{
			IsAlpha: "455455445",
		}},
	}
	err = validator.Scanner(badCaseStr)
	require.NotNil(t, err)
	require.Equal(t, "Required is required\nIsAlpha is not a valid alpha\nIsAlphanumeric is not a valid alpha numeric\nIsEmail is not a valid email\nIsStrongPassword is not a valid strong password\nIsUUID is not a valid UUID\nIsObjectId is not a valid ObjectID\nIsAlpha is not a valid alpha\nIsAlpha is not a valid alpha", err.Error())

	type BadCase struct {
		IsFloat      interface{} `validate:"isFloat"`
		IsInt        interface{} `validate:"isInt"`
		IsBool       interface{} `validate:"isBool"`
		IsDateString interface{} `validate:"isDateString"`
		IsNumber     interface{} `validate:"isNumber"`
	}

	badCaseNum := &BadCase{
		IsFloat:      true,
		IsInt:        "Abc",
		IsBool:       23,
		IsDateString: "33",
		IsNumber:     "fff",
	}
	err = validator.Scanner(badCaseNum)
	require.NotNil(t, err)
	require.Equal(t, "IsFloat is not a valid float\nIsInt is not a valid int\nIsBool is not a valid bool\nIsDateString is not a valid date time\nIsNumber is not a valid number", err.Error())
}

func Benchmark_Scanner(b *testing.B) {

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
			require.Nil(b, validator.Scanner(user))
			userList = append(userList, user)
		}

		require.Nil(b, validator.Scanner(&UserList{
			Users: userList,
		}))
	})
}

func TestDefault(t *testing.T) {
	type Pagination struct {
		Page  int `validate:"isInt" default:"1"`
		Limit int `validate:"isInt" default:"10"`
	}

	pagin := &Pagination{}
	require.Nil(t, validator.Scanner(pagin))
	require.Equal(t, 1, pagin.Page)
	require.Equal(t, 10, pagin.Limit)

	pagin2 := &Pagination{Page: 4, Limit: 20}
	require.Nil(t, validator.Scanner(pagin2))
	require.Equal(t, 4, pagin2.Page)
	require.Equal(t, 20, pagin2.Limit)
}
