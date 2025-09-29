// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package defs_test

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"testing"

	"github.com/AlphaOne1/midgard/defs"
)

func TestLogNil(t *testing.T) {
	t.Parallel()

	var m *defs.MWBase

	if m.Log() != slog.Default() {
		t.Errorf("expected default logger on nil MWBase")
	}
}

func TestSetLogNil(t *testing.T) {
	t.Parallel()

	var m *defs.MWBase

	if err := m.SetLog(slog.Default()); err == nil {
		t.Errorf("expected error on setting logger on nil MWBase")
	}
}

func TestSetLogger(t *testing.T) {
	t.Parallel()

	var m defs.MWBase
	var testLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	err := (&m).SetLog(testLogger)

	if err != nil || (&m).Log() != testLogger {
		t.Errorf("expected set testLogger on MWBase")
	}
}

func TestSetNilLogger(t *testing.T) {
	t.Parallel()

	var m defs.MWBase

	err := (&m).SetLog(nil)

	if err == nil {
		t.Errorf("expected setting nil logger on MWBase to fail")
	}
}

func TestLogLevelNil(t *testing.T) {
	t.Parallel()

	var m *defs.MWBase

	if m.LogLevel() != slog.LevelInfo {
		t.Errorf("expected INFO level on nil MWBase")
	}
}

func TestSetLogLevelNil(t *testing.T) {
	t.Parallel()

	var m *defs.MWBase

	if err := m.SetLogLevel(slog.LevelInfo); err == nil {
		t.Errorf("expected error on setting log level on nil MWBase")
	}
}

func TestSetLogLevel(t *testing.T) {
	t.Parallel()

	var m defs.MWBase

	err := (&m).SetLogLevel(slog.LevelWarn)

	if err != nil || (&m).LogLevel() != slog.LevelWarn {
		t.Errorf("expected set correct logLevel on MWBase")
	}
}

func TestNextNil(t *testing.T) {
	t.Parallel()

	var m *defs.MWBase

	if m.Next() != nil {
		t.Errorf("expected nil next on nil MWBase")
	}
}

func TestSetNextNil(t *testing.T) {
	t.Parallel()

	var m *defs.MWBase

	if err := m.SetNext(nil); err == nil {
		t.Errorf("expected error on setting next on nil MWBase")
	}
}

func TestSetNext(t *testing.T) {
	t.Parallel()

	var m defs.MWBase
	testHandler := http.Handler(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))

	err := (&m).SetNext(testHandler)

	if err != nil || fmt.Sprintf("%p", (&m).Next()) != fmt.Sprintf("%p", testHandler) {
		t.Errorf("expected correctly set next on MWBase")
	}
}

func TestSetNilNext(t *testing.T) {
	t.Parallel()

	var m defs.MWBase

	err := (&m).SetNext(nil)

	if err == nil {
		t.Errorf("expected error on setting nil next on MWBase")
	}
}

type TestMWBaser struct {
	defs.MWBase
}

func (h *TestMWBaser) GetMWBase() *defs.MWBase {
	return &h.MWBase
}

func TestWithLogger(t *testing.T) {
	t.Parallel()

	testLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	option := defs.WithLogger[*TestMWBaser](testLogger)
	testHandler := TestMWBaser{}

	err := option(&testHandler)

	if err != nil {
		t.Errorf("expected no error on setting logger on MWBase")
	}
}

func TestWithLoggerNil(t *testing.T) {
	t.Parallel()

	testLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	option := defs.WithLogger[defs.MWBaser](testLogger)
	var testHandler defs.MWBaser

	err := option(testHandler)

	if err == nil {
		t.Errorf("expected error on setting logger on nil MWBase")
	}
}

func TestWithLogLevel(t *testing.T) {
	t.Parallel()

	option := defs.WithLogLevel[*TestMWBaser](slog.LevelWarn)
	testHandler := TestMWBaser{}

	err := option(&testHandler)

	if err != nil {
		t.Errorf("expected no error on setting loglevel option on MWBase")
	}
}

func TestWithLogLevelNil(t *testing.T) {
	t.Parallel()

	option := defs.WithLogLevel[defs.MWBaser](slog.LevelWarn)
	var testHandler defs.MWBaser

	err := option(testHandler)

	if err == nil {
		t.Errorf("expected error on setting loglevel option on nil MWBase")
	}
}
