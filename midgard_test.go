// Copyright the midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package midgard

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlphaOne1/midgard/defs"
	"github.com/AlphaOne1/midgard/handler/correlation"
	"github.com/AlphaOne1/midgard/handler/method_filter"
	"github.com/AlphaOne1/midgard/util"
)

func TestStackMiddleware(t *testing.T) {
	newMWHandler := StackMiddlewareHandler(
		[]defs.Middleware{
			util.Must(method_filter.New(
				method_filter.WithMethods([]string{http.MethodGet}))),
			util.Must(correlation.New()),
		},
		http.HandlerFunc(util.DummyHandler),
	)

	_ = newMWHandler

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	newMWHandler.ServeHTTP(res, req)

	if res.Result().StatusCode != http.StatusOK {
		t.Errorf("stacking did not work, status not OK")
	}

	if res.Result().Header.Get("X-Correlation-ID") == "" {
		t.Errorf("stacking did not work, missing X-Correlation-ID")
	}

	_ = req.Body.Close()

	req = httptest.NewRequest(http.MethodPut, "/", nil)
	res = httptest.NewRecorder()

	newMWHandler.ServeHTTP(res, req)

	if res.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("stacking did not work, status not MethodNotAllowed")
	}

	_ = req.Body.Close()
}

func TestEmptyMiddlewareHandler(t *testing.T) {
	newMWHandler := StackMiddlewareHandler(
		[]defs.Middleware{},
		http.HandlerFunc(util.DummyHandler))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	newMWHandler.ServeHTTP(res, req)

	if res.Result().StatusCode != http.StatusOK {
		t.Errorf("stacking did not work, status not OK")
	}

	result := make([]byte, 5)

	if _, err := res.Result().Body.Read(result); err != nil {
		t.Errorf("could not read result: %v", err)
	}

	if string(result) != "dummy" {
		t.Errorf(`got wrong result, wanted "dummy" but got "%v"`, string(result))
	}

	_ = req.Body.Close()
}

func TestEmptyMiddleware(t *testing.T) {
	got := StackMiddleware([]defs.Middleware{})

	if got != nil {
		t.Errorf("expected nil on empty middleware stack")
	}
}
