package util

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestMustGood(t *testing.T) {
	got := Must("pass", nil)

	if got != "pass" {
		t.Errorf(`expected "pass" but got %v`, got)
	}
}

func TestMustBad(t *testing.T) {
	outbuf := bytes.Buffer{}

	slog.SetDefault(slog.New(slog.NewTextHandler(&outbuf, &slog.HandlerOptions{})))

	got := Must("nopass", errors.New("testerror"))

	if got != "nopass" {
		t.Errorf("got %v but wanted `nopass`", got)
	}

	outputMatch := regexp.MustCompile(`^time=[^ ]+ level=ERROR msg="must-condition not met" error=testerror\n`)

	if !outputMatch.Match(outbuf.Bytes()) {
		t.Errorf("output does not match format %v, got %v",
			outputMatch.String(),
			outbuf.String())
	}
}

func TestDummyHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	DummyHandler(rec, req)

	if rec.Body.String() != "dummy" {
		t.Errorf("wanted Dummy but got %v",
			slog.String("body", rec.Body.String()))
	}
}
