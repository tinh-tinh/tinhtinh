package validator_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/dto/validator"
)

func Test_Scanner_V2(t *testing.T) {
	v := validator.Validator{}
	type Enum int
	const (
		Pending Enum = iota
		Processing
		Completed
		Failed
	)
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
		IsDateString     string    `validate:"isDateString"`
		IsNumber         int       `validate:"isNumber"`
		IsNumber2        float64   `validate:"isNumber"`
		IsObjectId       string    `validate:"isObjectId"`
		Nested           *Nested   `validate:"nested"`
		Slice            []*Nested `validate:"nested"`
		Lala             string    `validate:"isAlpha"`
		Date             time.Time `validate:"isDate"`
		Enum             Enum      `validate:"isInt"`
		MinLength        string    `validate:"minLength=3"`
		MaxLength        string    `validate:"maxLength=10"`
	}
	err := v.Validate(Input{})
	assert.NotNil(t, err)

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
		IsDateString:     time.Now().Format("2006-01-02"),
		IsNumber:         123,
		IsNumber2:        39.49,
		IsObjectId:       "5e9bf1f6d3d2d3d3d3d3d3d3",
		Nested: &Nested{
			IsAlpha: "abc",
		},
		Slice: []*Nested{
			{IsAlpha: "avc"},
		},
		Date:      time.Now(),
		Enum:      Pending,
		MinLength: "abcd",
		MaxLength: "xyz",
	}

	err = v.Validate(happyCase)
	assert.Nil(t, err)

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
		MinLength: "a",
		MaxLength: "qwerteryuiiuoopo[o[bggnfnghmj,sccsbbhmjk,kk]]",
	}
	err = v.Validate(badCaseStr)
	assert.NotNil(t, err)
	assert.Equal(t, "Required is required\nIsAlpha is not a valid alpha\nIsAlphanumeric is not a valid alpha numeric\nIsEmail is not a valid email\nIsStrongPassword is not a valid strong password\nIsUUID is not a valid UUID\nIsObjectId is not a valid ObjectId\nIsAlpha is not a valid alpha\nIsAlpha is not a valid alpha\nMinLength minimum length is 3\nMaxLength maximum length is 10", err.Error())

	type BadCase struct {
		IsFloat      any `validate:"isFloat"`
		IsInt        any `validate:"isInt"`
		IsBool       any `validate:"isBool"`
		IsDateString any `validate:"isDateString"`
		IsNumber     any `validate:"isNumber"`
	}

	badCaseNum := &BadCase{
		IsFloat:      true,
		IsInt:        "Abc",
		IsBool:       23,
		IsDateString: "33",
		IsNumber:     "fff",
	}
	err = v.Validate(badCaseNum)
	assert.NotNil(t, err)
	assert.Equal(t, "IsFloat is not a valid float\nIsInt is not a valid int\nIsBool is not a valid bool\nIsDateString is not a valid date time\nIsNumber is not a valid number", err.Error())

	err = v.Validate(&CustomScanner{
		IsCustom: "abc",
	})
	assert.Nil(t, err)

	customScanner := &CustomScanner{
		IsCustom: "def",
	}
	err = v.Validate(customScanner)
	assert.NotNil(t, err)
	assert.Equal(t, "custom scan error", err.Error())
}

func TestDefault_V2(t *testing.T) {
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

func Test_Array_V2(t *testing.T) {
	type User struct {
		ArrEmail []string  `validate:"isEmail"`
		ArrInt   []int     `validate:"isInt"`
		ArrFloat []float64 `validate:"isFloat"`
		ArrBool  []bool    `validate:"isBool"`
	}

	user := &User{
		ArrEmail: []string{"abc@gmail.com", "abc@mailinator.ai"},
		ArrInt:   []int{1, 2},
		ArrFloat: []float64{1.1, 2.2},
		ArrBool:  []bool{true, false},
	}

	v := validator.Validator{}
	require.Nil(t, v.Validate(user))
}

func TestOptional_V2(t *testing.T) {
	type Contact struct {
		ContactName  string `json:"name" validator:"isAlphanumeric"`
		ContactEmail string `json:"email" validate:"isEmail"`
		ContactPhone string `json:"phone"`
	}

	type Location struct {
		Name     *string  `json:"name,omitempty"`
		Line     *string  `json:"addressLine,omitempty"`
		City     *string  `json:"city,omitempty"`
		State    *string  `json:"state,omitempty"`
		Country  *string  `json:"country,omitempty"`
		Zipcode  *string  `json:"zipCode,omitempty"`
		Timezone *string  `json:"timezone,omitempty"`
		MapUrl   *string  `json:"mapUrl,omitempty"`
		Contact  *Contact `json:"contact,omitempty" validate:"nested"`
	}

	input := &Location{}
	name := "name"
	input.Name = &name

	v := validator.Validator{}
	err := v.Validate(input)
	require.Nil(t, err)
}
