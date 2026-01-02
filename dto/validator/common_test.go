package validator_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/dto/validator"
)

func Test_IsEmpty(t *testing.T) {
	t.Parallel()

	assert.True(t, validator.IsEmpty(nil))
	assert.True(t, validator.IsEmpty([]string{}))

	a := []interface{}{}
	assert.True(t, validator.IsEmpty(a))

	b := make(map[string]interface{})
	assert.True(t, validator.IsEmpty(b))

	type Pointer struct{}
	var c *Pointer
	assert.True(t, validator.IsEmpty(c))
}

func Test_Validators_Edge_Cases(t *testing.T) {
	t.Parallel()

	// Test nil cases
	assert.False(t, validator.IsAlpha(nil))
	assert.False(t, validator.IsAlphanumeric(nil))
	assert.False(t, validator.IsEmail(nil))
	assert.False(t, validator.IsStrongPassword(nil))
	assert.False(t, validator.IsUUID(nil))
	assert.False(t, validator.IsObjectId(nil))
	assert.False(t, validator.IsInt(nil))
	assert.False(t, validator.IsFloat(nil))
	assert.False(t, validator.IsNumber(nil))
	assert.False(t, validator.IsDate(nil))
	assert.False(t, validator.IsBool(nil))

	// Test empty string cases
	assert.False(t, validator.IsAlpha(""))
	assert.False(t, validator.IsAlphanumeric(""))
	assert.False(t, validator.IsEmail(""))
	assert.False(t, validator.IsStrongPassword(""))
	assert.False(t, validator.IsUUID(""))
	assert.False(t, validator.IsObjectId(""))
	assert.False(t, validator.IsInt(""))
	assert.False(t, validator.IsFloat(""))
	assert.False(t, validator.IsNumber(""))
	assert.False(t, validator.IsDate(""))
	assert.False(t, validator.IsBool(""))
}

func Test_MinLength(t *testing.T) {
	t.Parallel()

	assert.True(t, validator.MinLength("test", 4))
	assert.False(t, validator.MinLength("test", 5))

	assert.True(t, validator.MinLength([]string{"test", "test2"}, 2))
	assert.False(t, validator.MinLength([]string{"test", "test2"}, 5))

	assert.False(t, validator.MinLength(nil, 0))

	assert.False(t, validator.MinLength(123, 0))
	assert.False(t, validator.MinLength(123.123, 0))
	assert.False(t, validator.MinLength(true, 0))
	assert.False(t, validator.MinLength(false, 0))
	assert.False(t, validator.MinLength(map[string]interface{}{}, 0))
	assert.True(t, validator.MinLength([]interface{}{}, 0))
}

func Test_MaxLength(t *testing.T) {
	t.Parallel()

	assert.True(t, validator.MaxLength("test", 4))
	assert.False(t, validator.MaxLength("test", 3))

	assert.True(t, validator.MaxLength([]string{"test", "test2"}, 2))
	assert.False(t, validator.MaxLength([]string{"test", "test2"}, 1))

	assert.False(t, validator.MaxLength(nil, 0))

	assert.False(t, validator.MaxLength(123, 0))
	assert.False(t, validator.MaxLength(123.123, 0))
	assert.False(t, validator.MaxLength(true, 0))
	assert.False(t, validator.MaxLength(false, 0))
	assert.False(t, validator.MaxLength(map[string]interface{}{}, 0))
	assert.True(t, validator.MaxLength([]interface{}{}, 0))
}

func Test_DDOS(t *testing.T) {
	t.Parallel()

	bigString := randomBigStr()
	require.False(t, validator.IsAlphanumeric(bigString))
}

func randomBigStr() string {
	var bigString strings.Builder
	// Define the number of repetitions
	repeat := 100000000
	smallString := "Hello, Go! "

	// Append the small string multiple times
	for i := 0; i < repeat; i++ {
		bigString.WriteString(smallString)
	}

	// Convert the builder to a string
	result := bigString.String()
	return result
}
