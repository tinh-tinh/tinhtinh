package validator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tinh-tinh/tinhtinh/v2/dto/validator"
)

func Test_IsBool(t *testing.T) {
	t.Parallel()

	assert.True(t, validator.IsBool(true))
	assert.False(t, validator.IsBool(123))

	assert.True(t, validator.IsBool("true"))
	assert.False(t, validator.IsBool("123"))

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
