// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package basic_auth_test

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AlphaOne1/midgard"
	"github.com/AlphaOne1/midgard/defs"
	"github.com/AlphaOne1/midgard/handler/basic_auth"
	"github.com/AlphaOne1/midgard/util"
)

type AuthTest struct{}

func (a *AuthTest) Authenticate(username, password string) (bool, error) {
	if password == "generr" {
		return false, errors.New("generated")
	}

	return username == "testuser" && password == "testpass", nil
}

func TestBasicAuth(t *testing.T) {
	tests := []struct {
		User      string
		Pass      string
		WantState int
	}{
		{
			User:      "testuser",
			Pass:      "testpass",
			WantState: http.StatusOK,
		},
		{
			User:      "testuser",
			Pass:      "testwrong",
			WantState: http.StatusUnauthorized,
		},
		{
			User:      "testuser",
			Pass:      "",
			WantState: http.StatusUnauthorized,
		},
		{
			User:      "",
			Pass:      "testpass",
			WantState: http.StatusUnauthorized,
		},
		{
			User:      "",
			Pass:      "",
			WantState: http.StatusUnauthorized,
		},
		{
			User:      "testuser",
			Pass:      "generr",
			WantState: http.StatusUnauthorized,
		},
	}

	handler := midgard.StackMiddlewareHandler(
		[]defs.Middleware{
			util.Must(basic_auth.New(
				basic_auth.WithAuthenticator(&AuthTest{}),
				basic_auth.WithRealm("testrealm"))),
		},
		http.HandlerFunc(util.DummyHandler),
	)

	for k, v := range tests {
		t.Run(fmt.Sprintf("TestBasicAuth-%d", k), func(t *testing.T) {
			t.Parallel()

			req, _ := http.NewRequestWithContext(t.Context(), "GET", "/", nil)
			rec := httptest.NewRecorder()

			req.Header.Add(
				"Authorization",
				"Basic "+
					base64.StdEncoding.EncodeToString([]byte(v.User+":"+v.Pass)))

			handler.ServeHTTP(rec, req)

			if rec.Result().StatusCode != v.WantState {
				t.Errorf("got state %v but wanted %v",
					rec.Result().StatusCode,
					v.WantState)
			}
		})
	}
}

func TestBasicAuthDecodeError(t *testing.T) {
	t.Parallel()

	handler := midgard.StackMiddlewareHandler(
		[]defs.Middleware{
			util.Must(basic_auth.New(
				basic_auth.WithAuthenticator(&AuthTest{}),
				basic_auth.WithRealm("testrealm"))),
		},
		http.HandlerFunc(util.DummyHandler),
	)

	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	req.Header.Add(
		"Authorization",
		"Basic absoluteNonsense==")

	handler.ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusUnauthorized {
		t.Errorf("got state %v but wanted %v",
			rec.Result().StatusCode,
			http.StatusUnauthorized)
	}
}

func TestBasicAuthNoAuthenticator(t *testing.T) {
	t.Parallel()

	_, mwErr := basic_auth.New(
		basic_auth.WithRealm("testrealm"))

	if mwErr == nil || mwErr.Error() != "no authenticator configured" {
		t.Errorf("uncought undefined authenticator")
	}
}

func TestBasicAuthDefaultRealm(t *testing.T) {
	t.Parallel()

	handler := midgard.StackMiddlewareHandler(
		[]defs.Middleware{
			util.Must(basic_auth.New(basic_auth.WithAuthenticator(&AuthTest{}))),
		},
		http.HandlerFunc(util.DummyHandler),
	)

	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusUnauthorized {
		t.Errorf("got state %v but wanted %v",
			rec.Result().StatusCode,
			http.StatusUnauthorized)
	}

	authHeader := rec.Result().Header.Get("WWW-Authenticate")

	if !strings.Contains(authHeader, `Basic realm="Restricted"`) {
		t.Errorf("default realm not set correctly: %v", authHeader)
	}
}

func TestBasicAuthRedirect(t *testing.T) {
	t.Parallel()

	handler := midgard.StackMiddlewareHandler(
		[]defs.Middleware{
			util.Must(basic_auth.New(
				basic_auth.WithAuthenticator(&AuthTest{}),
				basic_auth.WithRedirect("/login.html"),
			)),
		},
		http.HandlerFunc(util.DummyHandler),
	)

	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusFound {
		t.Errorf("got state %v but wanted %v",
			rec.Result().StatusCode,
			http.StatusFound)
	}

	relocHeader := rec.Result().Header.Get("Location")

	if !strings.Contains(relocHeader, `/login.html`) {
		t.Errorf("redirect not set correctly: %v", relocHeader)
	}
}
