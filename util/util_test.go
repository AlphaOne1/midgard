// Copyright the midgard contributors.
// SPDX-License-Identifier: MPL-2.0
package util

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/AlphaOne1/midgard/defs"
)

func TestMustGood(t *testing.T) {
	got := Must("pass", nil)

	if got != "pass" {
		t.Errorf(`expected "pass" but got %v`, got)
	}
}

func TestMustBad(t *testing.T) {
	outbuf := bytes.Buffer{}
	exitFunc = func(_ int) {}
	defer func() { exitFunc = os.Exit }()

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

func TestGetOrCreateID(t *testing.T) {
	tests := []struct {
		in      string
		wantNew bool
	}{
		{
			in:      "",
			wantNew: true,
		},
		{
			in:      "nonsense",
			wantNew: false,
		},
	}

	for k, v := range tests {
		got := GetOrCreateID(v.in)

		if v.wantNew == true && got == v.in {
			t.Errorf("%v: wanted new UUID but got old one", k)
		}

		if !v.wantNew == true && got != v.in {
			t.Errorf("%v: wanted old UUID but got new one", k)
		}
	}
}

func TestDummyHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	DummyHandler(rec, req)

	if rec.Body.String() != "dummy" {
		t.Errorf("wanted Dummy but got %v", rec.Body.String())
	}
}

func TestWriteState(t *testing.T) {
	tests := []struct {
		state int
	}{
		{
			state: http.StatusOK,
		},
		{
			state: http.StatusBadRequest,
		},
		{
			state: http.StatusAccepted,
		},
	}

	for k, v := range tests {
		rec := httptest.NewRecorder()
		WriteState(rec, slog.Default(), v.state)

		if rec.Body.String() != http.StatusText(v.state) {
			t.Errorf("%v: wanted %v but got %v", k, http.StatusText(v.state), rec.Body.String())
		}

		if ct := rec.Result().Header["Content-Type"]; len(ct) == 0 || ct[0] != "text/plain; charset=utf-8" {
			t.Errorf("%v: content type not set correctly, set to %v", k, ct)
		}

		if cto := rec.Result().Header["X-Content-Type-Options"]; len(cto) == 0 || cto[0] != "nosniff" {
			t.Errorf("%v: content type options not set correctly, set to %v", k, cto)
		}
	}
}

type MWTest struct {
	defs.MWBase
}

func (h *MWTest) GetMWBase() *defs.MWBase {
	return &h.MWBase
}

func TestIntroCheck(t *testing.T) {
	tests := []struct {
		h    *MWTest
		req  *http.Request
		want bool
	}{
		{
			h:    &MWTest{},
			req:  Must(http.NewRequest(http.MethodGet, "/", nil)),
			want: true,
		},
		{
			h:    nil,
			req:  Must(http.NewRequest(http.MethodGet, "/", nil)),
			want: false,
		},
		{
			h:    &MWTest{},
			req:  nil,
			want: false,
		},
		{
			h:    nil,
			req:  nil,
			want: false,
		},
	}

	rec := httptest.NewRecorder()

	for k, v := range tests {
		got := IntroCheck(v.h, rec, v.req)

		if got != v.want {
			t.Errorf("%v: got %v but wanted %v", k, got, v.want)
		}
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
