// SPDX-FileCopyrightText: 2026 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package correlation_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlphaOne1/midgard/handler/correlation"
	"github.com/AlphaOne1/midgard/helper"
)

func TestCorrelationNewID(t *testing.T) {
	t.Parallel()

	var gotCorrelationHeaderInside bool

	insideHandler := func(_ http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Correlation-ID") != "" {
			gotCorrelationHeaderInside = true
		}
	}

	handler := helper.Must(correlation.New())(http.HandlerFunc(insideHandler))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !gotCorrelationHeaderInside {
		t.Errorf("no X-Correlation-ID header added to request")
	}

	if rec.Header().Get("X-Correlation-ID") == "" {
		t.Errorf("no X-Correlation-ID header in response")
	}
}

func TestCorrelationSuppliedID(t *testing.T) {
	t.Parallel()

	var gotCorrelationHeaderInside bool

	insideHandler := func(_ http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Correlation-ID") == "setOutside" {
			gotCorrelationHeaderInside = true
		}
	}

	handler := helper.Must(correlation.New())(http.HandlerFunc(insideHandler))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Correlation-ID", "setOutside")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !gotCorrelationHeaderInside {
		t.Errorf("X-Correlation-ID header not added correctly to request")
	}

	if rec.Header().Get("X-Correlation-ID") != "setOutside" {
		t.Errorf("X-Correlation-ID header not set correctly in response")
	}
}
