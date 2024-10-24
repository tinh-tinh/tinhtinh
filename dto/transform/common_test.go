package transform

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_ToBool(t *testing.T) {
	require.True(t, ToBool("true"))
	require.False(t, ToBool(false))

	require.Panics(t, func() { ToBool(123) })
}

func Test_ToInt(t *testing.T) {
	require.Equal(t, 123, ToInt("123"))
	require.Equal(t, 456, ToInt(456))

	require.Panics(t, func() { ToInt(true) })
}

func Test_ToFloat(t *testing.T) {
	require.Equal(t, 0.45, ToFloat("0.45"))
	require.Equal(t, float32(0.45), ToFloat(float32(0.45)))
	require.Equal(t, float64(0.45), ToFloat(float64(0.45)))
	require.Equal(t, 1.0, ToFloat(1))

	require.Panics(t, func() { ToFloat(true) })
}

func Test_StringToDate(t *testing.T) {
	require.Equal(t, time.Date(2024, 9, 27, 0, 0, 0, 0, time.UTC), StringToDate("2024-09-27"))
}

func Test_StringToBool(t *testing.T) {
	require.True(t, StringToBool("true"))
	require.False(t, StringToBool("false"))
}

func Test_StringToInt(t *testing.T) {
	require.Equal(t, int64(123), StringToInt64("123"))
}

func Test_StringToTimeDuration(t *testing.T) {
	require.Equal(t, time.Second, StringToTimeDuration("1s"))
	require.Equal(t, time.Minute, StringToTimeDuration("1m"))
	require.Equal(t, time.Hour, StringToTimeDuration("1h"))
}
