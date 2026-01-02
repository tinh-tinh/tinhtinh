package validator_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tinh-tinh/tinhtinh/v2/dto/validator"
)

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
