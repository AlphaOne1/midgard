// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package cors

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/AlphaOne1/midgard/util"
)

//
// Basic Handler
//

func TestHandlerNil(t *testing.T) {
	t.Parallel()

	var handler *Handler

	if got := handler.GetMWBase(); got != nil {
		t.Errorf("MWBase of nil must be nil, but got non-nil")
	}

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("expected %v but got %v", http.StatusInternalServerError, rec.Result().StatusCode)
	}
}

//
// Generic Options
//

func TestOptionError(t *testing.T) {
	t.Parallel()

	errOpt := func(h *Handler) error {
		return errors.New("testerror")
	}

	_, err := New(errOpt)

	if err == nil {
		t.Errorf("expected middleware creation to fail")
	}
}

func TestOptionNil(t *testing.T) {
	t.Parallel()

	_, err := New(nil)

	if err == nil {
		t.Errorf("expected middleware creation to fail")
	}
}

func TestHandlerNextNil(t *testing.T) {
	t.Parallel()

	h := util.Must(New(WithLogLevel(slog.LevelDebug)))(nil)

	if h != nil {
		t.Errorf("expected handler to be nil")
	}
}

//
// WithLevel
//

func TestOptionWithLevel(t *testing.T) {
	t.Parallel()

	h := util.Must(New(WithLogLevel(slog.LevelDebug)))(http.HandlerFunc(util.DummyHandler))

	if h.(*Handler).LogLevel() != slog.LevelDebug {
		t.Errorf("wanted loglevel debug not set")
	}
}

func TestOptionWithLevelOnNil(t *testing.T) {
	t.Parallel()

	err := WithLogLevel(slog.LevelDebug)(nil)

	if err == nil {
		t.Errorf("expted error on configuring nil handler")
	}
}

//
// WithLogger
//

func TestOptionWithLogger(t *testing.T) {
	t.Parallel()

	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	h := util.Must(New(WithLogger(l)))(http.HandlerFunc(util.DummyHandler))

	if h.(*Handler).Log() != l {
		t.Errorf("logger not set correctly")
	}
}

func TestOptionWithLoggerOnNil(t *testing.T) {
	t.Parallel()

	err := WithLogger(slog.Default())(nil)

	if err == nil {
		t.Errorf("expted error on configuring nil handler")
	}
}

func TestOptionWithNilLogger(t *testing.T) {
	t.Parallel()

	var l *slog.Logger = nil
	_, hErr := New(WithLogger(l))

	if hErr == nil {
		t.Errorf("expected error on configuration with nil logger")
	}
}
