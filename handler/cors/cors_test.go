package cors

import (
	"bytes"
	"errors"
	"io"
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
		header      map[string]string
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
			header:      map[string]string{"Origin": "localhost"},
			wantCode:    http.StatusOK,
			wantHeader:  map[string]string{"Access-Control-Allow-Origin": "*"},
			wantContent: "dummy",
		}, { // 3
			cssMethods:  []string{http.MethodGet},
			cssOrigins:  []string{"dummy.com", "dummy1.com"},
			method:      http.MethodGet,
			header:      map[string]string{"Origin": "dummy.com"},
			wantCode:    http.StatusOK,
			wantHeader:  map[string]string{"Access-Control-Allow-Origin": "dummy.com"},
			wantContent: "dummy",
		}, { // 4
			cssMethods:  []string{http.MethodGet},
			cssOrigins:  []string{"dummy0.com", "dummy1.com"},
			method:      http.MethodGet,
			header:      map[string]string{"Origin": "dummy.com"},
			wantCode:    http.StatusForbidden,
			wantHeader:  nil,
			wantContent: "origin [dummy.com] not allowed",
		},
	}

	for k, v := range tests {
		req, _ := http.NewRequest(v.method, "http://dummy.com:8080", strings.NewReader(""))

		for hk, hv := range v.header {
			req.Header.Set(hk, hv)
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

func TestOptionError(t *testing.T) {
	errOpt := func(h *Handler) error {
		return errors.New("testerror")
	}

	_, err := New(errOpt)

	if err == nil {
		t.Errorf("expected middleware creation to fail")
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
