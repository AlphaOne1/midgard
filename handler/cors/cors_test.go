// Copyright the midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package cors

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/AlphaOne1/midgard/util"
)

func TestEvalCSSHandler(t *testing.T) {
	tests := []struct {
		cssMethods  []string
		cssOrigins  []string
		method      string
		header      map[string][]string
		wantCode    int
		wantHeader  map[string]string
		wantContent string
	}{
		{ // 0
			cssMethods:  []string{http.MethodGet},
			cssOrigins:  []string{"*"},
			method:      http.MethodGet,
			header:      nil,
			wantCode:    http.StatusOK,
			wantHeader:  nil,
			wantContent: "dummy",
		}, { // 1
			cssMethods:  []string{http.MethodGet},
			cssOrigins:  []string{"*"},
			method:      http.MethodOptions,
			header:      nil,
			wantCode:    http.StatusOK,
			wantHeader:  map[string]string{"Access-Control-Allow-Origin": "*"},
			wantContent: "",
		}, { // 2
			cssMethods:  []string{http.MethodGet},
			cssOrigins:  []string{"*"},
			method:      http.MethodGet,
			header:      map[string][]string{"Origin": {"localhost"}},
			wantCode:    http.StatusOK,
			wantHeader:  map[string]string{"Access-Control-Allow-Origin": "*"},
			wantContent: "dummy",
		}, { // 3
			cssMethods:  []string{http.MethodGet},
			cssOrigins:  []string{"dummy.com", "dummy1.com"},
			method:      http.MethodGet,
			header:      map[string][]string{"Origin": {"dummy.com"}},
			wantCode:    http.StatusOK,
			wantHeader:  map[string]string{"Access-Control-Allow-Origin": "dummy.com"},
			wantContent: "dummy",
		}, { // 4
			cssMethods:  []string{http.MethodGet},
			cssOrigins:  []string{"dummy0.com", "dummy1.com"},
			method:      http.MethodGet,
			header:      map[string][]string{"Origin": {"dummy.com"}},
			wantCode:    http.StatusForbidden,
			wantHeader:  nil,
			wantContent: http.StatusText(http.StatusForbidden),
		}, { // 5
			cssMethods:  []string{http.MethodGet},
			cssOrigins:  []string{"dummy0.com", "dummy1.com"},
			method:      http.MethodPost,
			header:      map[string][]string{"Origin": {"dummy0.com"}},
			wantCode:    http.StatusMethodNotAllowed,
			wantHeader:  nil,
			wantContent: http.StatusText(http.StatusMethodNotAllowed),
		}, { // 6
			cssMethods: []string{http.MethodGet},
			cssOrigins: []string{"dummy0.com", "dummy1.com"},
			method:     http.MethodGet,
			header: map[string][]string{
				"Origin":      {"dummy0.com"},
				"X-Forbidden": {"forbidden"},
			},
			wantCode:    http.StatusForbidden,
			wantHeader:  nil,
			wantContent: http.StatusText(http.StatusForbidden),
		}, { // 7
			cssMethods: []string{http.MethodGet},
			cssOrigins: []string{"dummy0.com", "dummy1.com"},
			method:     http.MethodGet,
			header: map[string][]string{
				"Origin": {"dummy0.com", "dummy1.com"},
			},
			wantCode:    http.StatusOK,
			wantHeader:  nil,
			wantContent: "dummy",
		}, { // 8
			cssMethods: []string{http.MethodGet},
			cssOrigins: []string{"dummy0.com", "dummy1.com"},
			method:     http.MethodGet,
			header: map[string][]string{
				"Origin": {"", "dummy0.com"},
			},
			wantCode:    http.StatusOK,
			wantHeader:  nil,
			wantContent: "dummy",
		}, { // 9
			cssMethods: []string{http.MethodGet},
			cssOrigins: []string{"dummy0.com", "dummy1.com"},
			method:     http.MethodGet,
			header: map[string][]string{
				"Origin": {},
			},
			wantCode:    http.StatusForbidden,
			wantHeader:  nil,
			wantContent: http.StatusText(http.StatusForbidden),
		},
	}

	for k, v := range tests {
		req, _ := http.NewRequest(v.method, "http://dummy.com:8080", strings.NewReader(""))

		for hk, hv := range v.header {
			for _, hvi := range hv {
				req.Header.Add(hk, hvi)
			}
		}

		rec := httptest.NewRecorder()

		mw := util.Must(New(
			WithMethods(v.cssMethods),
			WithHeaders(MinimumAllowHeaders()),
			WithOrigins(v.cssOrigins)))(http.HandlerFunc(util.DummyHandler))

		mw.ServeHTTP(rec, req)

		if rec.Code != v.wantCode {
			t.Errorf("%v: css filter did not work as expected, wanted %v but got %v", k, v.wantCode, rec.Code)
		}

		if rec.Body.String() != v.wantContent {
			t.Errorf("%v: wanted '%v' in body, but got '%v'", k, v.wantContent, rec.Body.String())
		}

		for wk, wv := range v.wantHeader {
			if val, found := rec.Result().Header[wk]; !found || !slices.Contains(val, wv) {
				t.Errorf("%v: wanted [%v:%v] but did not find it", k, wk, wv)
			}
		}
	}
}
