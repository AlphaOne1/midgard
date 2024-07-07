package local_limit

import (
	"testing"
	"time"

	"github.com/AlphaOne1/midgard/util"
)

func TestLocalLimitRate(t *testing.T) {
	tests := []struct {
		TargetRate   float64
		SleepTime    time.Duration
		TestDuration time.Duration
		WantDrops    int
	}{
		{
			TargetRate:   1,
			SleepTime:    1 * time.Second,
			TestDuration: 3 * time.Second,
			WantDrops:    3,
		},
		{
			TargetRate:   100,
			SleepTime:    200 * time.Millisecond,
			TestDuration: 420 * time.Millisecond,
			WantDrops:    40,
		},
		{
			TargetRate:   0.5,
			SleepTime:    200 * time.Millisecond,
			TestDuration: 2100 * time.Millisecond,
			WantDrops:    1,
		},
	}

	for k, v := range tests {
		got := 0

		limiter := util.Must(New(
			WithTargetRate(v.TargetRate),
			WithSleepInterval(v.SleepTime)))

		startTime := time.Now()

		for time.Since(startTime) < v.TestDuration {
			if limiter.Limit() {
				got++
			}
		}

		if got < v.WantDrops-1 || got > v.WantDrops+1 {
			t.Errorf("%v: got %v drops but wanted %v", k, got, v.WantDrops)
		}

		limiter.Stop()
	}
}
