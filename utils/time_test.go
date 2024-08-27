package utils

import (
	"testing"
	"time"
)

func checkTimestamp(tb testing.TB, expectedCurrent, actualCurrent uint32) {
	isRight := actualCurrent >= expectedCurrent-1 || actualCurrent <= expectedCurrent+1
	if !isRight {
		tb.Errorf("testing error\n Expect value: %d\nActual value: %d", expectedCurrent, actualCurrent)
	}
}

func Test_TimestampUpdater(t *testing.T) {
	t.Parallel()

	StartTimeStampUpdater()

	now := uint32(time.Now().Unix())
	checkTimestamp(t, now, Timestamp())

	// one second later
	time.Sleep(1 * time.Second)
	checkTimestamp(t, now+1, Timestamp())

	// two second later
	time.Sleep(1 * time.Second)
	checkTimestamp(t, now+2, Timestamp())
}

func Benchmark_CalculateTimestamp(b *testing.B) {
	var res uint32
	StartTimeStampUpdater()

	b.Run("Test Benchmark Timestamp", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		for n := 0; n < bb.N; n++ {
			_ = Timestamp()
		}
	})
	b.Run("default", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		for n := 0; n < bb.N; n++ {
			_ = uint32(time.Now().Unix())
		}
	})

	b.Run("Test Benchmark Timestamp Asserted", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		for n := 0; n < bb.N; n++ {
			res = Timestamp()
			checkTimestamp(bb, uint32(time.Now().Unix()), res)
		}
	})
	b.Run("default asserted", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		for n := 0; n < bb.N; n++ {
			res = uint32(time.Now().Unix())
			checkTimestamp(bb, uint32(time.Now().Unix()), res)
		}
	})
}
