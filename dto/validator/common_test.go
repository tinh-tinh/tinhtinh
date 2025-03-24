package validator_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/dto/validator"
)

func Test_IsAlpha(t *testing.T) {
	t.Parallel()

	require.True(t, validator.IsAlpha("test"))
	require.True(t, validator.IsAlpha("Abcd"))
	require.True(t, validator.IsAlpha("teDFD"))
	require.True(t, validator.IsAlpha("ZERERV"))

	require.False(t, validator.IsAlpha("12334"))
	require.False(t, validator.IsAlpha("Testdvdv1"))
	require.False(t, validator.IsAlpha("Testdvdv!"))
	require.False(t, validator.IsAlpha("Testdvdv "))

	require.False(t, validator.IsAlpha(123))
}

func Test_IsAlphanumeric(t *testing.T) {
	t.Parallel()

	require.True(t, validator.IsAlphanumeric("test"))
	require.True(t, validator.IsAlphanumeric("12333"))

	require.False(t, validator.IsAlphanumeric("test1243!"))
	require.False(t, validator.IsAlphanumeric("!test1243"))
	require.False(t, validator.IsAlphanumeric("test 1243"))

	require.False(t, validator.IsAlphanumeric(true))
}

func Test_IsEmail(t *testing.T) {
	t.Parallel()

	require.True(t, validator.IsEmail("abc@gmail.com"))
	require.True(t, validator.IsEmail("abc@mailinator.ai"))

	require.False(t, validator.IsEmail("abc"))
	require.False(t, validator.IsEmail("abc@gmail"))
	require.False(t, validator.IsEmail("abc@gmailcom"))
	require.False(t, validator.IsEmail("abc@gmail."))
	require.False(t, validator.IsEmail("abcgmail@.com"))
	require.False(t, validator.IsEmail("abc@gmail. "))

	require.False(t, validator.IsEmail(123))
}

func Test_IsStrongPassword(t *testing.T) {
	t.Parallel()

	require.False(t, validator.IsStrongPassword("12345678"))
	require.False(t, validator.IsStrongPassword("abcderf"))
	require.False(t, validator.IsStrongPassword("nf38yhg847"))
	require.False(t, validator.IsStrongPassword("@##$%$#$%$#"))
	require.False(t, validator.IsStrongPassword("cdncndndndc"))
	require.False(t, validator.IsStrongPassword("1a@"))
	require.False(t, validator.IsStrongPassword("1aABCD@"))

	require.True(t, validator.IsStrongPassword("12345678@Ab"))
	require.True(t, validator.IsStrongPassword("ACDBD@123def"))
	require.True(t, validator.IsStrongPassword("1adsd#@@Ab"))
	require.True(t, validator.IsStrongPassword("23Lsdvn@!bcd"))
	require.True(t, validator.IsStrongPassword("@#%#$^%^YS@a12"))
	require.True(t, validator.IsStrongPassword("@#54353FVERGEededr"))

	require.False(t, validator.IsStrongPassword(123))
}

func Test_IsUUID(t *testing.T) {
	t.Parallel()

	require.True(t, validator.IsUUID("550e8400-e29b-41d4-a716-446655440000"))

	require.False(t, validator.IsUUID("a0eebc99"))
	require.False(t, validator.IsUUID(123))
}

func Test_IsObjectId(t *testing.T) {
	t.Parallel()

	require.True(t, validator.IsObjectId("6592008029c8c3e4dc76256c"))

	require.False(t, validator.IsObjectId("a0eebc99"))
	require.False(t, validator.IsObjectId("6592dj8029c8c3e4dc76256c"))
	require.False(t, validator.IsObjectId(123))
}

func Test_IsRegexMatch(t *testing.T) {
	t.Parallel()

	require.True(t, validator.IsRegexMatch("^[a-f0-9]{24}$", "6592008029c8c3e4dc76256c"))
	require.False(t, validator.IsRegexMatch("^[a-f0-9]{24}$", "a0eebc99"))
	require.False(t, validator.IsRegexMatch("^[a-f0-9]{24}$", true))
}

func Test_IsInt(t *testing.T) {
	t.Parallel()

	require.True(t, validator.IsInt(123))
	require.False(t, validator.IsInt(123.123))
	require.True(t, validator.IsInt("123"))

	require.False(t, validator.IsInt(true))

	type Level int
	const (
		Low Level = iota
		Medium
		High
	)

	var level Level = High
	require.True(t, validator.IsInt(level))
}

func Test_IsFloat(t *testing.T) {
	t.Parallel()

	require.True(t, validator.IsFloat(123.123))
	require.False(t, validator.IsFloat(123))
	require.True(t, validator.IsFloat("123.123"))

	require.False(t, validator.IsFloat(true))

	type GepCode float32
	const (
		Asia    GepCode = 1.1
		Europe  GepCode = 1.2
		America GepCode = 1.3
	)

	var gepCode GepCode = Asia
	require.True(t, validator.IsFloat(gepCode))
}

func Test_IsNumber(t *testing.T) {
	t.Parallel()

	require.True(t, validator.IsNumber(123))
	require.True(t, validator.IsNumber(123.123))
	require.True(t, validator.IsNumber("123"))

	require.False(t, validator.IsNumber(true))
}

func Test_IsDateString(t *testing.T) {
	t.Parallel()

	require.True(t, validator.IsDate("2020-01-01"))
	require.True(t, validator.IsDate(time.Now()))

	require.False(t, validator.IsDate(123))
}

func Test_IsBool(t *testing.T) {
	t.Parallel()

	require.True(t, validator.IsBool(true))
	require.False(t, validator.IsBool(123))

	require.True(t, validator.IsBool("true"))
	require.False(t, validator.IsBool(123))

	type Status bool
	const (
		Active   Status = true
		Inactive Status = false
	)

	var status Status = Active
	require.True(t, validator.IsBool(status))
}

func Test_IsNil(t *testing.T) {
	t.Parallel()

	require.True(t, validator.IsNil(nil))
	require.True(t, validator.IsNil([]string{}))

	a := []interface{}{}
	require.True(t, validator.IsNil(a))

	b := make(map[string]interface{})
	require.True(t, validator.IsNil(b))
}

func Test_Validators_Edge_Cases(t *testing.T) {
	t.Parallel()

	// Test nil cases
	require.False(t, validator.IsAlpha(nil))
	require.False(t, validator.IsAlphanumeric(nil))
	require.False(t, validator.IsEmail(nil))
	require.False(t, validator.IsStrongPassword(nil))
	require.False(t, validator.IsUUID(nil))
	require.False(t, validator.IsObjectId(nil))
	require.False(t, validator.IsInt(nil))
	require.False(t, validator.IsFloat(nil))
	require.False(t, validator.IsNumber(nil))
	require.False(t, validator.IsDate(nil))
	require.False(t, validator.IsBool(nil))

	// Test empty string cases
	require.False(t, validator.IsAlpha(""))
	require.False(t, validator.IsAlphanumeric(""))
	require.False(t, validator.IsEmail(""))
	require.False(t, validator.IsStrongPassword(""))
	require.False(t, validator.IsUUID(""))
	require.False(t, validator.IsObjectId(""))
	require.False(t, validator.IsInt(""))
	require.False(t, validator.IsFloat(""))
	require.False(t, validator.IsNumber(""))
	require.False(t, validator.IsDate(""))
	require.False(t, validator.IsBool(""))
}
