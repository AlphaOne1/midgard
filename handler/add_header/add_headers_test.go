// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package add_header

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlphaOne1/midgard/util"
)

func TestAddHeaders(t *testing.T) {
	t.Parallel()

	handler := util.Must(New(WithHeaders([][2]string{
		{"X-Test-Header", "testValue"},
	})))(http.HandlerFunc(util.DummyHandler))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Header().Get("X-Test-Header") != "testValue" {
		t.Errorf("X-Test-Header header not added correctly to request")
	}
}
