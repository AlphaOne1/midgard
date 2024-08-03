// Copyright the midgard contributors.
// SPDX-License-Identifier: MPL-2.0
package method_filter

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AlphaOne1/midgard/util"
)

func TestMethodFilter(t *testing.T) {
	tests := []struct {
		filter   []string
		method   string
		wantCode int
	}{
		{
			filter:   []string{http.MethodGet},
			method:   http.MethodGet,
			wantCode: http.StatusOK,
		},
		{
			filter:   []string{http.MethodOptions, http.MethodGet},
			method:   http.MethodGet,
			wantCode: http.StatusOK,
		},
		{
			filter:   []string{http.MethodGet},
			method:   http.MethodPost,
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			filter:   []string{http.MethodGet},
			method:   " ",
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			// specialty of go, that treats "" as GET
			filter:   []string{http.MethodGet},
			method:   "",
			wantCode: http.StatusOK,
		},
	}

	for k, v := range tests {
		req, _ := http.NewRequest(v.method, "", strings.NewReader(""))
		rec := httptest.NewRecorder()

		mw := util.Must(New(WithMethods(v.filter)))(http.HandlerFunc(util.DummyHandler))

		mw.ServeHTTP(rec, req)

		if rec.Code != v.wantCode {
			t.Errorf("%v: method filter did not work as expected, wanted %v but got %v", k, v.wantCode, rec.Code)
		}
	}
}

func TestHandlerNil(t *testing.T) {
	var subject *Handler = nil

	rec := httptest.NewRecorder()

	subject.ServeHTTP(rec, nil)

	if rec.Result().StatusCode != http.StatusServiceUnavailable {
		t.Errorf("ServeHTTP on nil handler should give error state")
	}

	body := bytes.Buffer{}

	_, _ = io.Copy(&body, rec.Body)

	if body.String() != "service not available" {
		t.Errorf("expected 'service not available' but got '%s'", body.String())
	}
}

func TestOptionError(t *testing.T) {
	errOpt := func(h *Handler) error {
		return errors.New("testerror")
	}

	_, err := New(errOpt)

	if err == nil {
		t.Errorf("expected middleware creation to fail")
	}
}

func FuzzMethodFilter(f *testing.F) {

	f.Add(http.MethodDelete)
	f.Add(http.MethodGet)
	f.Add(http.MethodOptions)
	f.Add(http.MethodPost)
	f.Add(http.MethodPut)

	activeFilter := map[string]bool{http.MethodOptions: true, http.MethodGet: true}
	mw := util.Must(New(WithMethods(util.MapKeys(activeFilter))))(http.HandlerFunc(util.DummyHandler))

	f.Fuzz(func(t *testing.T, method string) {
		if method == "" {
			// compensate Go NewRequest specialty, that treats "" as GET
			method = http.MethodGet
		}

		req, _ := http.NewRequest(method, "", strings.NewReader(""))
		rec := httptest.NewRecorder()

		mw.ServeHTTP(rec, req)

		if activeFilter[method] && rec.Code != http.StatusOK ||
			!activeFilter[method] && rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("method filter did not work as expected, method %v got %v", method, rec.Code)
		}
	})
}
