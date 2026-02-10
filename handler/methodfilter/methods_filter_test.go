// SPDX-FileCopyrightText: 2026 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package methodfilter_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AlphaOne1/midgard/handler/methodfilter"
	"github.com/AlphaOne1/midgard/helper"
)

func TestMethodFilter(t *testing.T) {
	t.Parallel()

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

	for k, test := range tests {
		t.Run(fmt.Sprintf("TestMethodFilter-%d", k), func(t *testing.T) {
			t.Parallel()

			req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, "", strings.NewReader(""))
			// set the method after, as Go could change it
			req.Method = test.method
			rec := httptest.NewRecorder()

			mw := helper.Must(methodfilter.New(methodfilter.WithMethods(test.filter)))(http.HandlerFunc(helper.DummyHandler))

			mw.ServeHTTP(rec, req)

			if rec.Code != test.wantCode {
				t.Errorf("method filter did not work as expected, wanted %v but got %v", test.wantCode, rec.Code)
			}
		})
	}
}

func TestMethodFilterUninitialized(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, "", strings.NewReader(""))
	// set the method after, as Go could change it
	req.Method = http.MethodGet
	rec := httptest.NewRecorder()

	mw := helper.Must(methodfilter.New())(http.HandlerFunc(helper.DummyHandler))

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
	mw := helper.Must(methodfilter.New(
		methodfilter.WithMethods(helper.MapKeys(activeFilter))))(http.HandlerFunc(helper.DummyHandler))

	f.Fuzz(func(t *testing.T, method string) {
		req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, "", strings.NewReader(""))
		req.Method = method
		rec := httptest.NewRecorder()

		mw.ServeHTTP(rec, req)

		if activeFilter[method] && rec.Code != http.StatusOK ||
			!activeFilter[method] && rec.Code != http.StatusMethodNotAllowed {

			t.Errorf("method filter did not work as expected, method %v got %v", method, rec.Code)
		}
	})
}
