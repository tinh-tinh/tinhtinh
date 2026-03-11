package validator_test

import (
	"testing"

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
