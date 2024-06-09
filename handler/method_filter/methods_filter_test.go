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
