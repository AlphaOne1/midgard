package util

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"regexp"
	"slices"
	"strings"
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

func TestMapKeys(t *testing.T) {
	tests := []struct {
		in   map[string]int
		want []string
	}{
		{
			in:   map[string]int{"1": 1, "2": 2, "3": 3},
			want: []string{"1", "2", "3"},
		},
		{
			in:   nil,
			want: nil,
		},
		{
			in:   map[string]int{},
			want: []string{},
		},
	}

	for k, v := range tests {
		got := MapKeys(v.in)
		slices.Sort(got)
		slices.Sort(v.want)

		if v.in == nil && got != nil {
			t.Errorf("%v: got non-nil result %v but wanted nil", k, got)
		}

		if len(v.in) == 0 && v.in != nil && (len(got) != 0 || got == nil) {
			t.Errorf("%v: got %v but wanted zero length non-nil result", k, got)
		}

		if len(v.in) > 0 && strings.Join(got, ",") != strings.Join(v.want, ",") {
			t.Errorf("%v: got %v but wanted %v", k, got, v.want)
		}
	}
}
