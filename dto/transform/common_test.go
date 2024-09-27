package transform

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_ToBool(t *testing.T) {
	require.True(t, ToBool("true"))
	require.False(t, ToBool(false))
}

func Test_ToInt(t *testing.T) {
	require.Equal(t, 123, ToInt("123"))
	require.Equal(t, 456, ToInt(456))
}

func Test_ToFloat(t *testing.T) {
	require.Equal(t, 0.45, ToFloat("0.45"))
	require.Equal(t, float32(0.45), ToFloat(float32(0.45)))
	require.Equal(t, float64(0.45), ToFloat(float64(0.45)))
	require.Equal(t, 1.0, ToFloat(1))
}

func Test_StringToDate(t *testing.T) {
	require.Equal(t, time.Date(2024, 9, 27, 0, 0, 0, 0, time.UTC), StringToDate("2024-09-27"))
}
