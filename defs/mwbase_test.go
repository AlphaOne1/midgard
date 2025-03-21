// Copyright the midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package defs

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"testing"
)

func TestLogNil(t *testing.T) {
	var m *MWBase

	if m.Log() != slog.Default() {
		t.Errorf("expected default logger on nil MWBase")
	}
}

func TestSetLogNil(t *testing.T) {
	var m *MWBase

	if err := m.SetLog(slog.Default()); err == nil {
		t.Errorf("expected error on setting logger on nil MWBase")
	}
}

func TestSetLogger(t *testing.T) {
	var m MWBase
	var testLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	err := (&m).SetLog(testLogger)

	if err != nil || (&m).Log() != testLogger {
		t.Errorf("expected set testLogger on MWBase")
	}
}

func TestSetNilLogger(t *testing.T) {
	var m MWBase

	err := (&m).SetLog(nil)

	if err == nil {
		t.Errorf("expected setting nil logger on MWBase to fail")
	}
}

func TestLogLevelNil(t *testing.T) {
	var m *MWBase

	if m.LogLevel() != slog.LevelInfo {
		t.Errorf("expected INFO level on nil MWBase")
	}
}

func TestSetLogLevelNil(t *testing.T) {
	var m *MWBase

	if err := m.SetLogLevel(slog.LevelInfo); err == nil {
		t.Errorf("expected error on setting log level on nil MWBase")
	}
}

func TestSetLogLevel(t *testing.T) {
	var m MWBase

	err := (&m).SetLogLevel(slog.LevelWarn)

	if err != nil || (&m).LogLevel() != slog.LevelWarn {
		t.Errorf("expected set correct logLevel on MWBase")
	}
}

func TestNextNil(t *testing.T) {
	var m *MWBase

	if m.Next() != nil {
		t.Errorf("expected nil next on nil MWBase")
	}
}

func TestSetNextNil(t *testing.T) {
	var m *MWBase

	if err := m.SetNext(nil); err == nil {
		t.Errorf("expected error on setting next on nil MWBase")
	}
}

func TestSetNext(t *testing.T) {
	var m MWBase
	testHandler := http.Handler(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))

	err := (&m).SetNext(testHandler)

	if err != nil || fmt.Sprintf("%p", (&m).Next()) != fmt.Sprintf("%p", testHandler) {
		t.Errorf("expected correctly set next on MWBase")
	}
}

func TestSetNilNext(t *testing.T) {
	var m MWBase

	err := (&m).SetNext(nil)

	if err == nil {
		t.Errorf("expected error on setting nil next on MWBase")
	}
}

type TestMWBaser struct {
	MWBase
}

func (h *TestMWBaser) GetMWBase() *MWBase {
	return &h.MWBase
}

func TestWithLogger(t *testing.T) {
	testLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	option := WithLogger[*TestMWBaser](testLogger)
	testHandler := TestMWBaser{}

	err := option(&testHandler)

	if err != nil {
		t.Errorf("expected no error on setting logger on MWBase")
	}
}

func TestWithLoggerNil(t *testing.T) {
	testLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	option := WithLogger[MWBaser](testLogger)
	var testHandler MWBaser

	err := option(testHandler)

	if err == nil {
		t.Errorf("expected error on setting logger on nil MWBase")
	}
}

func TestWithLogLevel(t *testing.T) {
	option := WithLogLevel[*TestMWBaser](slog.LevelWarn)
	testHandler := TestMWBaser{}

	err := option(&testHandler)

	if err != nil {
		t.Errorf("expected no error on setting loglevel option on MWBase")
	}
}

func TestWithLogLevelNil(t *testing.T) {
	option := WithLogLevel[MWBaser](slog.LevelWarn)
	var testHandler MWBaser

	err := option(testHandler)

	if err == nil {
		t.Errorf("expected error on setting loglevel option on nil MWBase")
	}
}
