// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package locallimit_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/AlphaOne1/midgard/handler/ratelimit/locallimit"
	"github.com/AlphaOne1/midgard/helper"
)

func TestLocalLimitRate(t *testing.T) {
	t.Parallel()

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

	for k, test := range tests {
		t.Run(fmt.Sprintf("TestLocalLimitRate-%d", k), func(t *testing.T) {
			t.Parallel()

			got := 0

			limiter := helper.Must(locallimit.New(
				locallimit.WithTargetRate(test.TargetRate),
				locallimit.WithSleepInterval(test.SleepTime)))

			startTime := time.Now()

			for time.Since(startTime) < test.TestDuration {
				if limiter.Limit() {
					got++
				}
			}

			if got < test.WantDrops-1 || got > test.WantDrops+1 {
				t.Errorf("got %v drops but wanted %v", got, test.WantDrops)
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
		h := helper.Must(locallimit.New(locallimit.WithMaxDropsAbsolute(v.want)))
		t.Cleanup(h.Stop)

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
		t.Run(fmt.Sprintf("TestMaxDropsInterval-%d", k), func(t *testing.T) {
			t.Parallel()

			h := helper.Must(locallimit.New(
				locallimit.WithTargetRate(100),
				locallimit.WithMaxDropsInterval(v.want)))

			wantAbsolute := int64(v.want.Seconds() * h.TargetRate)

			if h.MaxDrops != wantAbsolute {
				t.Errorf("got %v wanted %v", h.MaxDrops, wantAbsolute)
			}
		})
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

	for k, test := range tests {
		t.Run(fmt.Sprintf("TestWithSleepInterval-%d", k), func(t *testing.T) {
			t.Parallel()

			h, hErr := locallimit.New(locallimit.WithSleepInterval(test.want))

			if (hErr != nil) != test.wantErr {
				t.Errorf("got error %v, want error %v", hErr, test.wantErr)
			}

			if hErr == nil {
				t.Cleanup(h.Stop)

				if h.SleepInterval != test.want {
					t.Errorf("got %v wanted %v", h.SleepInterval, test.want)
				}
			}
		})
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

	for k, test := range tests {
		t.Run(fmt.Sprintf("TestWithDropTimeout-%d", k), func(t *testing.T) {
			t.Parallel()

			h := helper.Must(locallimit.New(
				locallimit.WithTargetRate(0.1),
				locallimit.WithDropTimeout(test.want)))

			startTime := time.Now()
			h.Limit()
			duration := time.Since(startTime)
			h.Stop()

			if duration < time.Duration(float64(test.want)*0.95) ||
				duration > time.Duration(float64(test.want)*1.05) {

				t.Errorf("used %v but the timeout was %v", duration, test.want)
			}
		})
	}
}
