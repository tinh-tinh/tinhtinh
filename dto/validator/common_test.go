package validator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_IsAlpha(t *testing.T) {
	t.Parallel()

	require.True(t, IsAlpha("test"))
	require.True(t, IsAlpha("Abcd"))
	require.True(t, IsAlpha("teDFD"))
	require.True(t, IsAlpha("ZERERV"))

	require.False(t, IsAlpha("12334"))
	require.False(t, IsAlpha("Testdvdv1"))
	require.False(t, IsAlpha("Testdvdv!"))
	require.False(t, IsAlpha("Testdvdv "))

	require.False(t, IsAlpha(123))
}

func Test_IsAlphanumeric(t *testing.T) {
	t.Parallel()

	require.True(t, IsAlphanumeric("test"))
	require.True(t, IsAlphanumeric("12333"))

	require.False(t, IsAlphanumeric("test1243!"))
	require.False(t, IsAlphanumeric("!test1243"))
	require.False(t, IsAlphanumeric("test 1243"))

	require.False(t, IsAlphanumeric(true))
}

func Test_IsEmail(t *testing.T) {
	t.Parallel()

	require.True(t, IsEmail("abc@gmail.com"))
	require.True(t, IsEmail("abc@mailinator.ai"))

	require.False(t, IsEmail("abc"))
	require.False(t, IsEmail("abc@gmail"))
	require.False(t, IsEmail("abc@gmailcom"))
	require.False(t, IsEmail("abc@gmail."))
	require.False(t, IsEmail("abcgmail@.com"))
	require.False(t, IsEmail("abc@gmail. "))

	require.False(t, IsEmail(123))
}

func Test_IsStrongPassword(t *testing.T) {
	t.Parallel()

	require.False(t, IsStrongPassword("12345678"))
	require.False(t, IsStrongPassword("abcderf"))
	require.False(t, IsStrongPassword("nf38yhg847"))
	require.False(t, IsStrongPassword("@##$%$#$%$#"))
	require.False(t, IsStrongPassword("cdncndndndc"))
	require.False(t, IsStrongPassword("1a@"))
	require.False(t, IsStrongPassword("1aABCD@"))

	require.True(t, IsStrongPassword("12345678@Ab"))
	require.True(t, IsStrongPassword("ACDBD@123def"))
	require.True(t, IsStrongPassword("1adsd#@@Ab"))
	require.True(t, IsStrongPassword("23Lsdvn@!bcd"))
	require.True(t, IsStrongPassword("@#%#$^%^YS@a12"))
	require.True(t, IsStrongPassword("@#54353FVERGEededr"))

	require.False(t, IsStrongPassword(123))
}

func Test_IsUUID(t *testing.T) {
	t.Parallel()

	require.True(t, IsUUID("550e8400-e29b-41d4-a716-446655440000"))

	require.False(t, IsUUID("a0eebc99"))
	require.False(t, IsUUID(123))
}

func Test_IsObjectId(t *testing.T) {
	t.Parallel()

	require.True(t, IsObjectId("6592008029c8c3e4dc76256c"))

	require.False(t, IsObjectId("a0eebc99"))
	require.False(t, IsObjectId("6592dj8029c8c3e4dc76256c"))
	require.False(t, IsObjectId(123))
}

func Test_IsRegexMatch(t *testing.T) {
	t.Parallel()

	require.True(t, IsRegexMatch("^[a-f0-9]{24}$", "6592008029c8c3e4dc76256c"))
	require.False(t, IsRegexMatch("^[a-f0-9]{24}$", "a0eebc99"))
	require.False(t, IsRegexMatch("^[a-f0-9]{24}$", true))
}

func Test_IsInt(t *testing.T) {
	t.Parallel()

	require.True(t, IsInt(123))
	require.False(t, IsInt(123.123))
	require.True(t, IsInt("123"))

	require.False(t, IsInt(true))
}

func Test_IsFloat(t *testing.T) {
	t.Parallel()

	require.True(t, IsFloat(123.123))
	require.False(t, IsFloat(123))
	require.True(t, IsFloat("123.123"))

	require.False(t, IsFloat(true))
}

func Test_IsNumber(t *testing.T) {
	t.Parallel()

	require.True(t, IsNumber(123))
	require.True(t, IsNumber(123.123))
	require.True(t, IsNumber("123"))

	require.False(t, IsNumber(true))
}

func Test_IsDateString(t *testing.T) {
	t.Parallel()

	require.True(t, IsDateString("2020-01-01"))
	require.True(t, IsDateString(time.Now()))

	require.False(t, IsDateString(123))
}

func Test_IsBool(t *testing.T) {
	t.Parallel()

	require.True(t, IsBool(true))
	require.False(t, IsBool(123))

	require.True(t, IsBool("true"))
	require.False(t, IsBool(123))
}

func Test_IsNil(t *testing.T) {
	t.Parallel()

	require.True(t, IsNil(nil))
	require.True(t, IsNil([]string{}))

	a := []interface{}{}
	require.True(t, IsNil(a))

	b := make(map[string]interface{})
	require.True(t, IsNil(b))
}
