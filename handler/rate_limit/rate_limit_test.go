// Copyright the midgard contributors.
// SPDX-License-Identifier: MPL-2.0
package rate_limit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlphaOne1/midgard"
	"github.com/AlphaOne1/midgard/defs"
	"github.com/AlphaOne1/midgard/handler/rate_limit/local_limit"
	"github.com/AlphaOne1/midgard/util"
)

func TestRateLimit(t *testing.T) {

	handler := midgard.StackMiddlewareHandler(
		[]defs.Middleware{
			util.Must(New(WithLimiter(
				util.Must(local_limit.New(
					local_limit.WithTargetRate(20),
					local_limit.WithDropTimeout(15*time.Millisecond),
					local_limit.WithSleepInterval(100*time.Millisecond)))))),
		},
		http.HandlerFunc(util.DummyHandler))

	got := 0

	for i := 0; i < 30; i++ {
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
	_, err := New()

	if err == nil {
		t.Errorf("expected middleware creation to fail")
	}
}

func TestNilLimiter(t *testing.T) {
	_, err := New(WithLimiter(nil))

	if err == nil {
		t.Errorf("expected middleware creation to fail")
	}
}
