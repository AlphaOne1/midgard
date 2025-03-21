// Copyright the midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package method_filter

import (
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
		{ // 0
			filter:   []string{http.MethodGet},
			method:   http.MethodGet,
			wantCode: http.StatusOK,
		},
		{ // 1
			filter:   []string{http.MethodOptions, http.MethodGet},
			method:   http.MethodGet,
			wantCode: http.StatusOK,
		},
		{ // 2
			filter:   []string{http.MethodGet},
			method:   http.MethodPost,
			wantCode: http.StatusMethodNotAllowed,
		},
		{ // 3
			filter:   []string{http.MethodGet},
			method:   " ",
			wantCode: http.StatusMethodNotAllowed,
		},
		{ // 4
			filter:   []string{http.MethodGet},
			method:   "",
			wantCode: http.StatusMethodNotAllowed,
		},
		{ // 5
			filter:   []string{},
			method:   http.MethodGet,
			wantCode: http.StatusMethodNotAllowed,
		},
	}

	for k, v := range tests {
		req, _ := http.NewRequest(http.MethodGet, "", strings.NewReader(""))
		// set method after, as Go could change it
		req.Method = v.method
		rec := httptest.NewRecorder()

		mw := util.Must(New(WithMethods(v.filter)))(http.HandlerFunc(util.DummyHandler))

		mw.ServeHTTP(rec, req)

		if rec.Code != v.wantCode {
			t.Errorf("%v: method filter did not work as expected, wanted %v but got %v", k, v.wantCode, rec.Code)
		}
	}
}

func TestMethodFilterUninitialized(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "", strings.NewReader(""))
	// set method after, as Go could change it
	req.Method = http.MethodGet
	rec := httptest.NewRecorder()

	mw := util.Must(New())(http.HandlerFunc(util.DummyHandler))

	mw.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("method filter did not work as expected, wanted %v but got %v", http.StatusServiceUnavailable, rec.Code)
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
		req, _ := http.NewRequest(http.MethodGet, "", strings.NewReader(""))
		req.Method = method
		rec := httptest.NewRecorder()

		mw.ServeHTTP(rec, req)

		if activeFilter[method] && rec.Code != http.StatusOK ||
			!activeFilter[method] && rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("method filter did not work as expected, method %v got %v", method, rec.Code)
		}
	})
}
