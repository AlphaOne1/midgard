// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package correlation_test

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/AlphaOne1/midgard/handler/correlation"
	"github.com/AlphaOne1/midgard/helper"
)

//
// Basic Handler
//

func TestHandlerNil(t *testing.T) {
	t.Parallel()

	var handler *correlation.Handler

	if got := handler.GetMWBase(); got != nil {
		t.Errorf("MWBase of nil must be nil, but got non-nil")
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	//goland:noinspection GoMaybeNil
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

	errOpt := func( /* h */ *correlation.Handler) error {
		return errors.New("testerror")
	}

	_, err := correlation.New(errOpt)

	if err == nil {
		t.Errorf("expected middleware creation to fail")
	}
}

func TestOptionNil(t *testing.T) {
	t.Parallel()

	_, err := correlation.New(nil)

	if err == nil {
		t.Errorf("expected middleware creation to fail")
	}
}

func TestHandlerNextNil(t *testing.T) {
	t.Parallel()

	h := helper.Must(correlation.New(correlation.WithLogLevel(slog.LevelDebug)))(nil)

	if h != nil {
		t.Errorf("expected handler to be nil")
	}
}

//
// WithLevel
//

func TestOptionWithLevel(t *testing.T) {
	t.Parallel()

	h := helper.Must(correlation.New(correlation.WithLogLevel(slog.LevelDebug)))(http.HandlerFunc(helper.DummyHandler))

	val, isValid := h.(*correlation.Handler)

	if !isValid {
		t.Fatalf("wrong type")
	}

	if val.LogLevel() != slog.LevelDebug {
		t.Errorf("wanted loglevel debug not set")
	}
}

func TestOptionWithLevelOnNil(t *testing.T) {
	t.Parallel()

	err := correlation.WithLogLevel(slog.LevelDebug)(nil)

	if err == nil {
		t.Errorf("expected error on configuring nil handler")
	}
}

//
// WithLogger
//

func TestOptionWithLogger(t *testing.T) {
	t.Parallel()

	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	h := helper.Must(correlation.New(correlation.WithLogger(l)))(http.HandlerFunc(helper.DummyHandler))

	val, isValid := h.(*correlation.Handler)

	if !isValid {
		t.Fatalf("wrong type")
	}

	if val.Log() != l {
		t.Errorf("logger not set correctly")
	}
}

func TestOptionWithLoggerOnNil(t *testing.T) {
	t.Parallel()

	err := correlation.WithLogger(slog.Default())(nil)

	if err == nil {
		t.Errorf("expected error on configuring nil handler")
	}
}

func TestOptionWithNilLogger(t *testing.T) {
	t.Parallel()

	var l *slog.Logger
	_, hErr := correlation.New(correlation.WithLogger(l))

	if hErr == nil {
		t.Errorf("expected error on configuration with nil logger")
	}
}
