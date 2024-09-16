package validator

import (
	"testing"

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
}

func Test_IsAlphanumeric(t *testing.T) {
	t.Parallel()

	require.True(t, IsAlphanumeric("test"))
	require.True(t, IsAlphanumeric("12333"))

	require.False(t, IsAlphanumeric("test1243!"))
	require.False(t, IsAlphanumeric("!test1243"))
	require.False(t, IsAlphanumeric("test 1243"))
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
}
