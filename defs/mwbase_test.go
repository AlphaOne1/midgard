package defs

import (
	"log/slog"
	"testing"
)

func TestLogNil(t *testing.T) {
	var m *MWBase

	if m.Log() != nil {
		t.Errorf("expected nil logger on nil MWBase")
	}
}

func TestSetLogNil(t *testing.T) {
	var m *MWBase

	if err := m.SetLog(slog.Default()); err == nil {
		t.Errorf("expected error on setting logger on nil MWBase")
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

func TestNextNil(t *testing.T) {
	var m *MWBase

	if m.Next() != nil {
		t.Errorf("exptected nil next on nil MWBase")
	}
}

func TestSetNextNil(t *testing.T) {
	var m *MWBase

	if err := m.SetNext(nil); err == nil {
		t.Errorf("expected error on setting next on nil MWBase")
	}
}
