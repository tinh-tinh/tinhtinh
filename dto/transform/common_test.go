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

func Test_StringToDate(t *testing.T) {
	require.Equal(t, time.Date(2024, 9, 27, 0, 0, 0, 0, time.UTC), transform.StringToDate("2024-09-27"))
}

func Test_StringToBool(t *testing.T) {
	require.True(t, transform.StringToBool("true"))
	require.False(t, transform.StringToBool("false"))
}

func Test_StringToInt(t *testing.T) {
	require.Equal(t, int64(123), transform.StringToInt64("123"))
	require.Equal(t, int(123), transform.StringToInt("123"))
}

func Test_StringToTimeDuration(t *testing.T) {
	require.Equal(t, time.Second, transform.StringToTimeDuration("1s"))
	require.Equal(t, time.Minute, transform.StringToTimeDuration("1m"))
	require.Equal(t, time.Hour, transform.StringToTimeDuration("1h"))
}
