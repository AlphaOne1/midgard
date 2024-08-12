// Copyright the midgard contributors.
// SPDX-License-Identifier: MPL-2.0
package correlation

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/AlphaOne1/midgard/util"
)

func TestCorrelationNewID(t *testing.T) {
	var gotCorrelationHeaderInside bool

	insideHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Correlation-ID") != "" {
			gotCorrelationHeaderInside = true
		}
	}

	handler := util.Must(New())(http.HandlerFunc(insideHandler))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !gotCorrelationHeaderInside {
		t.Errorf("no X-Correlation-ID header added to request")
	}

	if rec.Header().Get("X-Correlation-iD") == "" {
		t.Errorf("no X-Correlation-ID header in response")
	}
}

func TestCorrelationSuppliedID(t *testing.T) {
	var gotCorrelationHeaderInside bool

	insideHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Correlation-ID") == "setOutside" {
			gotCorrelationHeaderInside = true
		}
	}

	handler := util.Must(New())(http.HandlerFunc(insideHandler))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Correlation-ID", "setOutside")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !gotCorrelationHeaderInside {
		t.Errorf("X-Correlation-ID header not added correctly to request")
	}

	if rec.Header().Get("X-Correlation-iD") != "setOutside" {
		t.Errorf("X-Correlation-ID header not set correctly in response")
	}
}

//
// WithLevel
//

func TestOptionWithLevel(t *testing.T) {
	h := util.Must(New(WithLogLevel(slog.LevelDebug)))(http.HandlerFunc(util.DummyHandler))

	if h.(*Handler).LogLevel() != slog.LevelDebug {
		t.Errorf("wanted loglevel debug not set")
	}
}

func TestOptionWithLevelOnNil(t *testing.T) {
	err := WithLogLevel(slog.LevelDebug)(nil)

	if err == nil {
		t.Errorf("expted error on configuring nil handler")
	}
}

//
// WithLogger
//

func TestOptionWithLogger(t *testing.T) {
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	h := util.Must(New(WithLogger(l)))(http.HandlerFunc(util.DummyHandler))

	if h.(*Handler).Log() != l {
		t.Errorf("logger not set correctly")
	}
}

func TestOptionWithLoggerOnNil(t *testing.T) {
	err := WithLogger(slog.Default())(nil)

	if err == nil {
		t.Errorf("expted error on configuring nil handler")
	}
}

func TestOptionWithNilLogger(t *testing.T) {
	var l *slog.Logger = nil
	_, hErr := New(WithLogger(l))

	if hErr == nil {
		t.Errorf("expected error on configuration with nil logger")
	}
}
