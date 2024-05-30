package handler

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

func TestEvalCSSHandler(t *testing.T) {
	tests := []struct {
		cssMethods  []string
		cssOrigins  []string
		method      string
		header      map[string]string
		wantCode    int
		wantHeader  map[string]string
		wantContent string
	}{
		{
			cssMethods:  []string{http.MethodGet},
			cssOrigins:  []string{"*"},
			method:      http.MethodGet,
			header:      nil,
			wantCode:    http.StatusOK,
			wantHeader:  nil,
			wantContent: "dummy",
		},
		{
			cssMethods:  []string{http.MethodGet},
			cssOrigins:  []string{"*"},
			method:      http.MethodOptions,
			header:      nil,
			wantCode:    http.StatusOK,
			wantHeader:  nil,
			wantContent: "",
		},
		{
			cssMethods:  []string{http.MethodGet},
			cssOrigins:  []string{"*"},
			method:      http.MethodGet,
			header:      map[string]string{"Origin": "localhost"},
			wantCode:    http.StatusOK,
			wantHeader:  map[string]string{"Access-Control-Allow-Origin": "*"},
			wantContent: "dummy",
		},
	}

	for k, v := range tests {
		req, _ := http.NewRequest(v.method, "", strings.NewReader(""))

		for hk, hv := range v.header {
			req.Header.Set(hk, hv)
		}

		rec := httptest.NewRecorder()

		mw := NewEvalCSSHandler(v.cssMethods, v.cssOrigins)(http.HandlerFunc(dummyHandler))

		mw.ServeHTTP(rec, req)

		if rec.Code != v.wantCode {
			t.Errorf("%v: method filter did not work as expected, wanted %v but got %v", k, v.wantCode, rec.Code)
		}

		if rec.Body.String() != v.wantContent {
			t.Errorf("%v: wanted '%v' in body, but got '%v'", k, v.wantContent, rec.Body.String())
		}

		for wk, wv := range v.wantHeader {
			if val, found := rec.Result().Header[wk]; !found || !slices.Contains(val, wv) {
				t.Errorf("%v", k)
			}
		}
	}
}
