package transform_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/dto/transform"
)

func Test_ToBool(t *testing.T) {
	require.True(t, transform.ToBool("true"))
	require.False(t, transform.ToBool(false))

	require.Panics(t, func() { transform.ToBool(123) })
}

func Test_ToInt(t *testing.T) {
	require.Equal(t, 123, transform.ToInt("123"))
	require.Equal(t, 456, transform.ToInt(456))

	require.Panics(t, func() { transform.ToInt(true) })
}

func Test_ToFloat(t *testing.T) {
	require.Equal(t, 0.45, transform.ToFloat("0.45"))
	require.Equal(t, float32(0.45), transform.ToFloat(float32(0.45)))
	require.Equal(t, float64(0.45), transform.ToFloat(float64(0.45)))
	require.Equal(t, 1.0, transform.ToFloat(1))

	require.Panics(t, func() { transform.ToFloat(true) })
}

func Test_ToDate(t *testing.T) {
	date, _ := time.Parse("2006-01-02", "2024-01-01")
	require.Equal(t, date, transform.ToDate("2024-01-01"))

	require.Panics(t, func() { transform.ToDate(true) })
}

func Test_ToString(t *testing.T) {
	require.Equal(t, "true", transform.ToString(true))
	require.Equal(t, "false", transform.ToString(false))
	require.Equal(t, "123", transform.ToString(123))
	require.Equal(t, "0.45", transform.ToString(0.45))
	require.Equal(t, "0.1", transform.ToString(float32(0.1)))
	require.Equal(t, "2024-01-01", transform.ToString("2024-01-01"))

	current := time.Now()
	require.Equal(t, current.String(), transform.ToString(current))

	require.Panics(t, func() { transform.ToString(make(map[string]interface{})) })
}
