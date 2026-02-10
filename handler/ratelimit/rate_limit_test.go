// SPDX-FileCopyrightText: 2026 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package ratelimit_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlphaOne1/midgard"
	"github.com/AlphaOne1/midgard/defs"
	"github.com/AlphaOne1/midgard/handler/ratelimit"
	"github.com/AlphaOne1/midgard/handler/ratelimit/locallimit"
	"github.com/AlphaOne1/midgard/helper"
)

func TestRateLimit(t *testing.T) {
	t.Parallel()

	handler := midgard.StackMiddlewareHandler(
		[]defs.Middleware{
			helper.Must(ratelimit.New(ratelimit.WithLimiter(
				helper.Must(locallimit.New(
					locallimit.WithTargetRate(20),
					locallimit.WithDropTimeout(15*time.Millisecond),
					locallimit.WithSleepInterval(100*time.Millisecond)))))),
		},
		http.HandlerFunc(helper.DummyHandler))

	got := 0

	for range 30 {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Result().StatusCode == http.StatusOK {
			got++
		}
	}

	if got != 6 {
		t.Errorf("got %d, want %d", got, 6)
	}
}

func TestNoOptions(t *testing.T) {
	t.Parallel()

	_, err := ratelimit.New()

	if err == nil {
		t.Errorf("expected middleware creation to fail")
	}
}

func TestNilLimiter(t *testing.T) {
	t.Parallel()

	_, err := ratelimit.New(ratelimit.WithLimiter(nil))

	if err == nil {
		t.Errorf("expected middleware creation to fail")
	}
}
