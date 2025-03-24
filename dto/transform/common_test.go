package transform_test

import (
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/v2/dto/transform"
)

func Test_ToBool(t *testing.T) {
	require.Equal(t, true, transform.ToBool("true"))
	require.Equal(t, false, transform.ToBool("false"))

	require.Panics(t, func() { transform.ToBool(123) })
	type Status bool
	const (
		Active   Status = true
		Inactive Status = false
	)
	var status Status = Active
	require.Equal(t, Active, transform.ToBool(status))
}

func Test_ToInt(t *testing.T) {
	require.Equal(t, 123, transform.ToInt("123"))
	require.Equal(t, 456, transform.ToInt(456))

	require.Panics(t, func() { transform.ToInt(true) })

	type Level int
	const (
		Low Level = iota
		Medium
		High
	)
	var level Level = High
	require.Equal(t, High, transform.ToInt(level))
}

func Test_ToFloat(t *testing.T) {
	require.Equal(t, 0.45, transform.ToFloat("0.45"))
	require.Equal(t, float32(0.45), transform.ToFloat(float32(0.45)))
	require.Equal(t, float64(0.45), transform.ToFloat(float64(0.45)))
	require.Equal(t, 1.0, transform.ToFloat(1))

	require.Panics(t, func() { transform.ToFloat(true) })

	type GepCode float32
	const (
		Asia    GepCode = 1.1
		Europe  GepCode = 1.2
		America GepCode = 1.3
	)
	var gepCode GepCode = Asia
	require.Equal(t, Asia, transform.ToFloat(gepCode))
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
	require.Equal(t, "123", transform.ToString(int64(123)))
	require.Equal(t, "-0.45", transform.ToString(-0.45))
	require.Equal(t, "0.1", transform.ToString(float32(0.1)))
	require.Equal(t, "2024-01-01", transform.ToString("2024-01-01"))

	current := time.Now()
	require.Equal(t, current.String(), transform.ToString(current))

	require.Panics(t, func() { transform.ToString(make(map[string]interface{})) })

	type Status string
	const (
		Pending    Status = "pending"
		Processing Status = "processing"
		Completed  Status = "completed"
		Failed     Status = "failed"
	)
	var status Status = Pending
	require.Equal(t, Pending, transform.ToString(status))

	require.Equal(t, strconv.Itoa(math.MaxInt64), transform.ToString(math.MaxInt64))
	require.Equal(t, strconv.Itoa(math.MinInt64), transform.ToString(math.MinInt64))

	loc, _ := time.LoadLocation("UTC")
	utcTime := time.Now().In(loc)
	require.Equal(t, utcTime.String(), transform.ToString(utcTime))
}
