// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package local_limit_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/AlphaOne1/midgard/handler/rate_limit/local_limit"
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
		t.Run(fmt.Sprintf("TestLocalLimitRate_%v", k), func(t *testing.T) {
			t.Parallel()

			got := 0

			limiter := util.Must(local_limit.New(
				local_limit.WithTargetRate(v.TargetRate),
				local_limit.WithSleepInterval(v.SleepTime)))

			startTime := time.Now()

			for time.Since(startTime) < v.TestDuration {
				if limiter.Limit() {
					got++
				}
			}

			if got < v.WantDrops-1 || got > v.WantDrops+1 {
				t.Errorf("got %v drops but wanted %v", got, v.WantDrops)
			}

			limiter.Stop()
		})
	}
}

func TestMaxDropsAbsolute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		want int64
	}{
		{want: 23},
		{want: 42},
		{want: 123},
	}

	for k, v := range tests {
		h := util.Must(local_limit.New(local_limit.WithMaxDropsAbsolute(v.want)))

		if h.MaxDrops != v.want {
			t.Errorf("%v: got %v wanted %v", k, h.MaxDrops, v.want)
		}
	}
}

func TestMaxDropsInterval(t *testing.T) {
	t.Parallel()

	tests := []struct {
		want time.Duration
	}{
		{want: 1 * time.Second},
		{want: 10 * time.Second},
		{want: 1400 * time.Millisecond},
	}

	for k, v := range tests {
		h := util.Must(local_limit.New(
			local_limit.WithTargetRate(100),
			local_limit.WithMaxDropsInterval(v.want)))

		wantAbsolute := int64(v.want.Seconds() * h.TargetRate)

		if h.MaxDrops != wantAbsolute {
			t.Errorf("%v: got %v wanted %v", k, h.MaxDrops, wantAbsolute)
		}
	}
}

func TestWithSleepInterval(t *testing.T) {
	t.Parallel()

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
		h, hErr := local_limit.New(local_limit.WithSleepInterval(v.want))

		if (hErr != nil) != v.wantErr {
			t.Errorf("%v: got error %v, want error %v", k, hErr, v.wantErr)
		}

		if hErr == nil && h.SleepInterval != v.want {
			t.Errorf("%v: got %v wanted %v", k, h.SleepInterval, v.want)
		}
	}
}

func TestWithDropTimeout(t *testing.T) {
	t.Parallel()

	tests := []struct {
		want time.Duration
	}{
		{want: 1 * time.Second},
		{want: 2 * time.Second},
		{want: 400 * time.Millisecond},
	}

	for k, v := range tests {
		h := util.Must(local_limit.New(
			local_limit.WithTargetRate(0.1),
			local_limit.WithDropTimeout(v.want)))

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
