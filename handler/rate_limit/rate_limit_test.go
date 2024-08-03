// Copyright the midgard contributors.
// SPDX-License-Identifier: MPL-2.0
package rate_limit

import (
	"bytes"
	"errors"
	"io"
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

func TestHandlerNil(t *testing.T) {
	var subject *Handler = nil

	rec := httptest.NewRecorder()

	subject.ServeHTTP(rec, nil)

	if rec.Result().StatusCode != http.StatusServiceUnavailable {
		t.Errorf("ServeHTTP on nil handler should give error state")
	}

	body := bytes.Buffer{}

	_, _ = io.Copy(&body, rec.Body)

	if body.String() != "service not available" {
		t.Errorf("expected 'service not available' but got '%s'", body.String())
	}
}

func TestOptionError(t *testing.T) {
	errOpt := func(h *Handler) error {
		return errors.New("testerror")
	}

	_, err := New(errOpt)

	if err == nil {
		t.Errorf("expected middleware creation to fail")
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
