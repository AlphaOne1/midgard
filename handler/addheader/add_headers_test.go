// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package addheader_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlphaOne1/midgard/handler/addheader"
	"github.com/AlphaOne1/midgard/helper"
)

func TestAddHeaders(t *testing.T) {
	t.Parallel()

	handler := helper.Must(addheader.New(addheader.WithHeaders([][2]string{
		{"X-Test-Header", "testValue"},
	})))(http.HandlerFunc(helper.DummyHandler))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Header().Get("X-Test-Header") != "testValue" {
		t.Errorf("X-Test-Header header not added correctly to request")
	}
}
