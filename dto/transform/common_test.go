package transform_test

import (
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tinh-tinh/tinhtinh/v2/dto/transform"
)

func Test_ToBool(t *testing.T) {
	assert.Equal(t, true, transform.ToBool("true"))
	assert.Equal(t, false, transform.ToBool("false"))

	assert.Panics(t, func() { transform.ToBool(123) })
	type Status bool
	const (
		Active   Status = true
		Inactive Status = false
	)
	var status Status = Active
	assert.Equal(t, Active, transform.ToBool(status))

	assert.Equal(t, []bool{true, false}, transform.ToBool([]string{"true", "false"}))
	assert.Panics(t, func() { transform.ToBool([]string{"true", "123"}) })
}

func Test_ToInt(t *testing.T) {
	assert.Equal(t, 123, transform.ToInt("123"))
	assert.Equal(t, 456, transform.ToInt(456))

	assert.Panics(t, func() { transform.ToInt(true) })

	type Level int
	const (
		Low Level = iota
		Medium
		High
	)
	var level Level = High
	assert.Equal(t, High, transform.ToInt(level))

	assert.Equal(t, []int{123, 456}, transform.ToInt([]string{"123", "456"}))
	assert.Panics(t, func() { transform.ToInt([]string{"123", "true"}) })
}

func Test_ToFloat(t *testing.T) {
	assert.Equal(t, 0.45, transform.ToFloat("0.45"))
	assert.Equal(t, float32(0.45), transform.ToFloat(float32(0.45)))
	assert.Equal(t, float64(0.45), transform.ToFloat(float64(0.45)))
	assert.Equal(t, 1.0, transform.ToFloat(1))

	assert.Panics(t, func() { transform.ToFloat(true) })

	type GepCode float32
	const (
		Asia    GepCode = 1.1
		Europe  GepCode = 1.2
		America GepCode = 1.3
	)
	var gepCode GepCode = Asia
	assert.Equal(t, Asia, transform.ToFloat(gepCode))

	assert.Equal(t, []float64{0.45, 1.0}, transform.ToFloat([]string{"0.45", "1"}))
	assert.Panics(t, func() { transform.ToFloat([]string{"0.45", "true"}) })
}

func Test_ToDate(t *testing.T) {
	date, _ := time.Parse("2006-01-02", "2024-01-01")
	assert.Equal(t, date, transform.ToDate("2024-01-01"))

	assert.Panics(t, func() { transform.ToDate(true) })

	assert.Equal(t, []time.Time{date}, transform.ToDate([]string{"2024-01-01"}))
	assert.Panics(t, func() { transform.ToDate([]string{"2024-01-01", "true"}) })
}

func Test_ToString(t *testing.T) {
	assert.Equal(t, "true", transform.ToString(true))
	assert.Equal(t, "false", transform.ToString(false))
	assert.Equal(t, "123", transform.ToString(123))
	assert.Equal(t, "123", transform.ToString(int64(123)))
	assert.Equal(t, "-0.45", transform.ToString(-0.45))
	assert.Equal(t, "0.1", transform.ToString(float32(0.1)))
	assert.Equal(t, "2024-01-01", transform.ToString("2024-01-01"))

	current := time.Now()
	assert.Equal(t, current.String(), transform.ToString(current))

	assert.Panics(t, func() { transform.ToString(make(map[string]interface{})) })

	type Status string
	const (
		Pending    Status = "pending"
		Processing Status = "processing"
		Completed  Status = "completed"
		Failed     Status = "failed"
	)
	var status Status = Pending
	assert.Equal(t, Pending, transform.ToString(status))

	assert.Equal(t, strconv.Itoa(math.MaxInt64), transform.ToString(math.MaxInt64))
	assert.Equal(t, strconv.Itoa(math.MinInt64), transform.ToString(math.MinInt64))

	loc, _ := time.LoadLocation("UTC")
	utcTime := time.Now().In(loc)
	assert.Equal(t, utcTime.String(), transform.ToString(utcTime))

	assert.Equal(t, []string{"true", "false"}, transform.ToString([]bool{true, false}))
	assert.Equal(t, []string{"123", "456"}, transform.ToString([]int{123, 456}))
	assert.Equal(t, []string{"0.45", "1"}, transform.ToString([]float64{0.45, 1.0}))

	assert.Panics(t, func() { transform.ToString([]any{make(map[string]interface{})}) })
}
