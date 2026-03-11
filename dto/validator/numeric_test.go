package validator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tinh-tinh/tinhtinh/v2/dto/validator"
)

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
