// Copyright the midgard contributors.
// SPDX-License-Identifier: MPL-2.0

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

func TestMaxDropsAbsolute(t *testing.T) {
	tests := []struct {
		want int64
	}{
		{want: 23},
		{want: 42},
		{want: 123},
	}

	for k, v := range tests {
		h := util.Must(New(WithMaxDropsAbsolute(v.want)))

		if h.MaxDrops != v.want {
			t.Errorf("%v: got %v wanted %v", k, h.MaxDrops, v.want)
		}
	}
}

func TestMaxDropsInterval(t *testing.T) {
	tests := []struct {
		want time.Duration
	}{
		{want: 1 * time.Second},
		{want: 10 * time.Second},
		{want: 1400 * time.Millisecond},
	}

	for k, v := range tests {
		h := util.Must(New(
			WithTargetRate(100),
			WithMaxDropsInterval(v.want)))

		wantAbsolute := int64(v.want.Seconds() * h.TargetRate)

		if h.MaxDrops != wantAbsolute {
			t.Errorf("%v: got %v wanted %v", k, h.MaxDrops, wantAbsolute)
		}
	}
}

func TestWithSleepInterval(t *testing.T) {
	tests := []struct {
		want    time.Duration
		wantErr bool
	}{
		{
			want:    10 * time.Millisecond,
			wantErr: false,
		},
		{
			want:    200 * time.Millisecond,
			wantErr: false,
		},
		{
			want:    1400 * time.Millisecond,
			wantErr: false,
		},
		{
			want:    0,
			wantErr: true,
		},
	}

	for k, v := range tests {
		h, hErr := New(WithSleepInterval(v.want))

		if (hErr != nil) != v.wantErr {
			t.Errorf("%v: got error %v, want error %v", k, hErr, v.wantErr)
		}

		if hErr == nil && h.SleepInterval != v.want {
			t.Errorf("%v: got %v wanted %v", k, h.SleepInterval, v.want)
		}
	}
}

func TestWithDropTimeout(t *testing.T) {
	tests := []struct {
		want time.Duration
	}{
		{want: 1 * time.Second},
		{want: 2 * time.Second},
		{want: 400 * time.Millisecond},
	}

	for k, v := range tests {
		h := util.Must(New(
			WithTargetRate(0.1),
			WithDropTimeout(v.want)))

		startTime := time.Now()
		h.Limit()
		duration := time.Since(startTime)
		h.Stop()

		if duration < time.Duration(float64(v.want)*0.95) ||
			duration > time.Duration(float64(v.want)*1.05) {
			t.Errorf("%v: used %v but the timeout was %v", k, duration, v.want)
		}

	}
}
