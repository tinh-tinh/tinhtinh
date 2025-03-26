package validator_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tinh-tinh/tinhtinh/v2/dto/validator"
)

func Test_IsAlpha(t *testing.T) {
	t.Parallel()

	assert.True(t, validator.IsAlpha("test"))
	assert.True(t, validator.IsAlpha("Abcd"))
	assert.True(t, validator.IsAlpha("teDFD"))
	assert.True(t, validator.IsAlpha("ZERERV"))

	assert.False(t, validator.IsAlpha("12334"))
	assert.False(t, validator.IsAlpha("Testdvdv1"))
	assert.False(t, validator.IsAlpha("Testdvdv!"))
	assert.False(t, validator.IsAlpha("Testdvdv "))

	assert.False(t, validator.IsAlpha(123))
	assert.False(t, validator.IsAlpha(true))

	assert.True(t, validator.IsAlpha([]string{"test", "Abcd"}))
	assert.False(t, validator.IsAlpha([]string{"test", "123"}))
	assert.True(t, validator.IsAlpha("test"))
	assert.False(t, validator.IsAlpha("123"))
}

func Test_IsAlphanumeric(t *testing.T) {
	t.Parallel()

	assert.True(t, validator.IsAlphanumeric("test"))
	assert.True(t, validator.IsAlphanumeric("12333"))

	assert.False(t, validator.IsAlphanumeric("test1243!"))
	assert.False(t, validator.IsAlphanumeric("!test1243"))
	assert.False(t, validator.IsAlphanumeric("test 1243"))

	assert.False(t, validator.IsAlphanumeric(true))

	assert.True(t, validator.IsAlphanumeric([]string{"test", "123"}))
	assert.False(t, validator.IsAlphanumeric([]string{"$#^$%^", "123"}))
	assert.True(t, validator.IsAlphanumeric("test"))
	assert.False(t, validator.IsAlphanumeric("$#^$%^"))
}

func Test_IsEmail(t *testing.T) {
	t.Parallel()

	assert.True(t, validator.IsEmail("abc@gmail.com"))
	assert.True(t, validator.IsEmail("abc@mailinator.ai"))

	assert.False(t, validator.IsEmail("abc"))
	assert.False(t, validator.IsEmail("abc@gmail"))
	assert.False(t, validator.IsEmail("abc@gmailcom"))
	assert.False(t, validator.IsEmail("abc@gmail."))
	assert.False(t, validator.IsEmail("abcgmail@.com"))
	assert.False(t, validator.IsEmail("abc@gmail. "))

	assert.False(t, validator.IsEmail(123))

	assert.True(t, validator.IsEmail([]string{"abc@gmail.com", "abc@mailinator.ai"}))
	assert.False(t, validator.IsEmail([]string{"abc", "abc@gmail.com"}))
	assert.True(t, validator.IsEmail("abc@gmail.com"))
	assert.False(t, validator.IsEmail("abc"))
}

func Test_IsStrongPassword(t *testing.T) {
	t.Parallel()

	assert.False(t, validator.IsStrongPassword("12345678"))
	assert.False(t, validator.IsStrongPassword("abcderf"))
	assert.False(t, validator.IsStrongPassword("nf38yhg847"))
	assert.False(t, validator.IsStrongPassword("@##$%$#$%$#"))
	assert.False(t, validator.IsStrongPassword("cdncndndndc"))
	assert.False(t, validator.IsStrongPassword("1a@"))
	assert.False(t, validator.IsStrongPassword("1aABCD@"))

	assert.True(t, validator.IsStrongPassword("12345678@Ab"))
	assert.True(t, validator.IsStrongPassword("ACDBD@123def"))
	assert.True(t, validator.IsStrongPassword("1adsd#@@Ab"))
	assert.True(t, validator.IsStrongPassword("23Lsdvn@!bcd"))
	assert.True(t, validator.IsStrongPassword("@#%#$^%^YS@a12"))
	assert.True(t, validator.IsStrongPassword("@#54353FVERGEededr"))

	assert.False(t, validator.IsStrongPassword(123))
}

func Test_IsUUID(t *testing.T) {
	t.Parallel()

	assert.True(t, validator.IsUUID("550e8400-e29b-41d4-a716-446655440000"))

	assert.False(t, validator.IsUUID("a0eebc99"))
	assert.False(t, validator.IsUUID(123))

	assert.True(t, validator.IsUUID([]string{"550e8400-e29b-41d4-a716-446655440000"}))
	assert.False(t, validator.IsUUID([]string{"a0eebc99", "550e8400-e29b-41d4-a716-446655440000"}))
	assert.True(t, validator.IsUUID("550e8400-e29b-41d4-a716-446655440000"))
	assert.False(t, validator.IsUUID("a0eebc99"))
}

func Test_IsObjectId(t *testing.T) {
	t.Parallel()

	assert.True(t, validator.IsObjectId("6592008029c8c3e4dc76256c"))

	assert.False(t, validator.IsObjectId("a0eebc99"))
	assert.False(t, validator.IsObjectId("6592dj8029c8c3e4dc76256c"))
	assert.False(t, validator.IsObjectId(123))

	assert.True(t, validator.IsObjectId([]string{"6592008029c8c3e4dc76256c"}))
	assert.False(t, validator.IsObjectId([]string{"a0eebc99", "6592008029c8c3e4dc76256c"}))
	assert.True(t, validator.IsObjectId("6592008029c8c3e4dc76256c"))
	assert.False(t, validator.IsObjectId("a0eebc99"))
}

func Test_IsRegexMatch(t *testing.T) {
	t.Parallel()

	assert.True(t, validator.IsRegexMatch("^[a-f0-9]{24}$", "6592008029c8c3e4dc76256c"))
	assert.False(t, validator.IsRegexMatch("^[a-f0-9]{24}$", "a0eebc99"))
	assert.False(t, validator.IsRegexMatch("^[a-f0-9]{24}$", true))
}

func Test_IsInt(t *testing.T) {
	t.Parallel()

	assert.True(t, validator.IsInt(123))
	assert.False(t, validator.IsInt(123.123))
	assert.True(t, validator.IsInt("123"))

	assert.False(t, validator.IsInt(true))

	type Level int
	const (
		Low Level = iota
		Medium
		High
	)

	var level Level = High
	assert.True(t, validator.IsInt(level))

	assert.True(t, validator.IsInt([]string{"123", "456"}))
	assert.False(t, validator.IsInt([]string{"123", "true"}))
	assert.True(t, validator.IsInt("123"))
	assert.False(t, validator.IsInt("true"))
}

func Test_IsFloat(t *testing.T) {
	t.Parallel()

	assert.True(t, validator.IsFloat(123.123))
	assert.False(t, validator.IsFloat(123))
	assert.True(t, validator.IsFloat("123.123"))

	assert.False(t, validator.IsFloat(true))

	type GepCode float32
	const (
		Asia    GepCode = 1.1
		Europe  GepCode = 1.2
		America GepCode = 1.3
	)

	var gepCode GepCode = Asia
	assert.True(t, validator.IsFloat(gepCode))

	assert.True(t, validator.IsFloat([]string{"123.123", "456.456"}))
	assert.False(t, validator.IsFloat([]string{"123", "true"}))
	assert.True(t, validator.IsFloat("123.123"))
	assert.False(t, validator.IsFloat("true"))
}

func Test_IsNumber(t *testing.T) {
	t.Parallel()

	assert.True(t, validator.IsNumber(123))
	assert.True(t, validator.IsNumber(123.123))
	assert.True(t, validator.IsNumber("123"))

	assert.False(t, validator.IsNumber(true))

	assert.True(t, validator.IsNumber([]string{"123", "456"}))
	assert.False(t, validator.IsNumber([]string{"123", "true"}))
	assert.True(t, validator.IsNumber("123"))
	assert.False(t, validator.IsNumber("true"))
}

func Test_IsDateString(t *testing.T) {
	t.Parallel()

	assert.True(t, validator.IsDate("2020-01-01"))
	assert.True(t, validator.IsDate(time.Now()))

	assert.False(t, validator.IsDate(123))

	assert.True(t, validator.IsDate([]string{"2020-01-01", "2020-01-02"}))
	assert.False(t, validator.IsDate([]string{"2020-01-01", "123"}))
	assert.True(t, validator.IsDate("2020-01-01"))
	assert.False(t, validator.IsDate("123"))
}

func Test_IsBool(t *testing.T) {
	t.Parallel()

	assert.True(t, validator.IsBool(true))
	assert.False(t, validator.IsBool(123))

	assert.True(t, validator.IsBool("true"))
	assert.False(t, validator.IsBool(123))

	type Status bool
	const (
		Active   Status = true
		Inactive Status = false
	)

	var status Status = Active
	assert.True(t, validator.IsBool(status))

	assert.True(t, validator.IsBool([]string{"true", "false"}))
	assert.False(t, validator.IsBool([]string{"true", "123"}))
	assert.True(t, validator.IsBool("true"))
	assert.False(t, validator.IsBool("123"))
}

func Test_IsNil(t *testing.T) {
	t.Parallel()

	assert.True(t, validator.IsNil(nil))
	assert.True(t, validator.IsNil([]string{}))

	a := []interface{}{}
	assert.True(t, validator.IsNil(a))

	b := make(map[string]interface{})
	assert.True(t, validator.IsNil(b))
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
