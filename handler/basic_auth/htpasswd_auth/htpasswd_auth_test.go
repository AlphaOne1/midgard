// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package htpasswd_auth_test

import (
	"testing"

	"github.com/AlphaOne1/midgard/handler/basic_auth/htpasswd_auth"
	"github.com/AlphaOne1/midgard/util"
)

func TestHtpasswdAuth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Username string
		Password string
		Valid    bool
	}{
		{
			Username: "user0",
			Password: "pass0",
			Valid:    true,
		},
		{
			Username: "user1",
			Password: "pass1",
			Valid:    true,
		},
		{
			Username: "user0",
			Password: "wrong0",
			Valid:    false,
		},
	}

	a := util.Must(htpasswd_auth.New(htpasswd_auth.WithAuthFile("testwd")))

	for k, v := range tests {
		gotAuth, gotErr := a.Authorize(v.Username, v.Password)

		if gotErr != nil {
			t.Errorf("%v: got error, but did not expect any: %v", k, gotErr)
		}

		if gotAuth != v.Valid {
			t.Errorf("%v: got auth %v but wanted %v", k, v.Valid, gotAuth)
		}
	}
}

func TestHtpasswdNil(t *testing.T) {
	t.Parallel()

	var subject *htpasswd_auth.HTPassWDAuth

	if _, err := subject.Authorize("u", "p"); err == nil {
		t.Errorf("authorize on nil authorizer should give error")
	}
}

func TestHtpasswdNonExistingFile(t *testing.T) {
	t.Parallel()

	_, err := htpasswd_auth.New(htpasswd_auth.WithAuthFile("IDoNotExistNowhereInThisWorldForgetIt"))

	if err == nil {
		t.Errorf("authorizer initialization with non-existent file should give error")
	}
}

func TestHtpasswdNoOptions(t *testing.T) {
	t.Parallel()

	_, err := htpasswd_auth.New()

	if err == nil {
		t.Errorf("authorizer initialization without options should give error")
	}
}

func TestHtpasswdWrongReader(t *testing.T) {
	t.Parallel()

	_, err := htpasswd_auth.New(htpasswd_auth.WithAuthInput(nil))

	if err == nil {
		t.Errorf("authorizer initialization nil reader should give error")
	}
}

func TestHtpasswdEmptyFilename(t *testing.T) {
	t.Parallel()

	_, err := htpasswd_auth.New(htpasswd_auth.WithAuthFile(""))

	if err == nil {
		t.Errorf("authorizer initialization with empty filename should give error")
	}
}
