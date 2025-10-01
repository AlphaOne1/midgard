// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package helper_test

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/AlphaOne1/midgard/defs"
	"github.com/AlphaOne1/midgard/helper"
)

func TestMustGood(t *testing.T) {
	t.Parallel()

	got := helper.Must("pass", nil)

	if got != "pass" {
		t.Errorf(`expected "pass" but got %v`, got)
	}
}

//nolint:paralleltest // manipulating global exit function
func TestMustBad(t *testing.T) {
	outbuf := bytes.Buffer{}
	*(helper.TexitFunc) = func(_ int) {}
	defer func() { *(helper.TexitFunc) = os.Exit }()

	slog.SetDefault(slog.New(slog.NewTextHandler(&outbuf, &slog.HandlerOptions{})))

	got := helper.Must("nopass", errors.New("testerror"))

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

func TestWrapIfError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		err        error
		text       string
		wantErr    bool
		wantErrMsg string
	}{
		{ // 0
			text:       "wrap this",
			err:        errors.New("original error"),
			wantErr:    true,
			wantErrMsg: "wrap this: original error",
		},
		{ // 1
			text:       "wrap that: %w",
			err:        errors.New("original error"),
			wantErr:    true,
			wantErrMsg: "wrap that: %w: original error",
		},
		{ // 2
			text:    "wrap this",
			err:     nil,
			wantErr: false,
		},
		{ // 3
			text:       "",
			err:        errors.New("error case"),
			wantErr:    true,
			wantErrMsg: "error case",
		},
		{ // 4
			text:    "",
			err:     nil,
			wantErr: false,
		},
	}

	for testIndex, test := range tests {
		t.Run(fmt.Sprintf("WrapError-%d", testIndex), func(t *testing.T) {
			t.Parallel()

			got := helper.WrapIfError(test.text, test.err)

			if (got != nil) != test.wantErr {
				t.Errorf(`got error "%v", but wanted "%v"`, got, test.wantErr)
			}

			if test.wantErr && got.Error() != test.wantErrMsg {
				t.Errorf(`got error "%v" but wanted "%v"`, got.Error(), test.wantErrMsg)
			}
		})
	}
}

func TestGetOrCreateID(t *testing.T) {
	t.Parallel()

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
		got := helper.GetOrCreateID(v.in)

		if v.wantNew && got == v.in {
			t.Errorf("%v: wanted new UUID but got old one", k)
		}

		if !v.wantNew && got != v.in {
			t.Errorf("%v: wanted old UUID but got new one", k)
		}
	}
}

func TestDummyHandler(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	helper.DummyHandler(rec, req)

	if rec.Body.String() != "dummy" {
		t.Errorf("wanted Dummy but got %v", rec.Body.String())
	}
}

func TestWriteState(t *testing.T) {
	t.Parallel()

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
		t.Run(fmt.Sprintf("TestWriteState-%d", k), func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			helper.WriteState(rec, slog.Default(), v.state)

			if rec.Body.String() != http.StatusText(v.state) {
				t.Errorf("wanted %v but got %v", http.StatusText(v.state), rec.Body.String())
			}

			if ct := rec.Result().Header["Content-Type"]; len(ct) == 0 || ct[0] != "text/plain; charset=utf-8" {
				t.Errorf("content type not set correctly, set to %v", ct)
			}

			if cto := rec.Result().Header["X-Content-Type-Options"]; len(cto) == 0 || cto[0] != "nosniff" {
				t.Errorf("content type options not set correctly, set to %v", cto)
			}
		})
	}
}

type MWTest struct {
	defs.MWBase
}

func (h *MWTest) GetMWBase() *defs.MWBase {
	return &h.MWBase
}

func TestIntroCheck(t *testing.T) {
	t.Parallel()

	tests := []struct {
		h    *MWTest
		req  *http.Request
		want bool
	}{
		{
			h:    &MWTest{},
			req:  helper.Must(http.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)),
			want: true,
		},
		{
			h:    nil,
			req:  helper.Must(http.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)),
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

	for k, v := range tests {
		t.Run(fmt.Sprintf("TestIntroCheck-%d", k), func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			got := helper.IntroCheck(v.h, rec, v.req)

			if got != v.want {
				t.Errorf("got %v but wanted %v", got, v.want)
			}
		})
	}
}

func TestMapKeys(t *testing.T) {
	t.Parallel()

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

	for k, test := range tests {
		t.Run(fmt.Sprintf("TestMapKeys-%d", k), func(t *testing.T) {
			t.Parallel()

			got := helper.MapKeys(test.in)
			slices.Sort(got)
			slices.Sort(test.want)

			if test.in == nil && got != nil {
				t.Errorf("got non-nil result %v but wanted nil", got)
			}

			if len(test.in) == 0 && test.in != nil && (len(got) != 0 || got == nil) {
				t.Errorf("got %v but wanted zero length non-nil result", got)
			}

			if len(test.in) > 0 && strings.Join(got, ",") != strings.Join(test.want, ",") {
				t.Errorf("got %v but wanted %v", got, test.want)
			}
		})
	}
}
