// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package mapauth_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/AlphaOne1/midgard/handler/basicauth/mapauth"
)

func TestMapAuthenticator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Auths      map[string]string
		User       string
		Pass       string
		WantNewErr bool
		Want       bool
		WantError  bool
	}{
		{
			Auths:      map[string]string{"testuser": "testpass"},
			User:       "testuser",
			Pass:       "testpass",
			WantNewErr: false,
			Want:       true,
			WantError:  false,
		},
		{
			Auths:      map[string]string{"testuser": "testpass"},
			User:       "testuser",
			Pass:       "testwrong",
			WantNewErr: false,
			Want:       false,
			WantError:  false,
		},
		{
			Auths:      map[string]string{},
			User:       "testuser",
			Pass:       "testpass",
			WantNewErr: true,
			Want:       false,
			WantError:  true,
		},
	}

	for k, test := range tests {
		t.Run(fmt.Sprintf("TestMapAuthenticator-%d", k), func(t *testing.T) {
			t.Parallel()

			auth, newErr := mapauth.New(mapauth.WithAuths(test.Auths))

			if newErr != nil {
				if !test.WantNewErr {
					t.Errorf("got error on creation, but wanted none")
				}

				return
			}

			if test.WantNewErr {
				t.Errorf("wanted error on creation, but got none")

				return
			}

			gotAuth, gotErr := auth.Authenticate(test.User, test.Pass)

			if gotErr != nil {
				if !test.WantError {
					t.Errorf("did not expect error, but got: %v", gotErr)
				}
				if gotAuth {
					t.Errorf("got error, so auth should not work, but got: %v", gotAuth)
				}
			} else {
				if test.WantError {
					t.Errorf("did expect error, but got none")
				}
				if gotAuth != test.Want {
					t.Errorf("got auth %v but wanted %v", gotAuth, test.Want)
				}
			}
		})
	}
}

func TestNoAuthorizations(t *testing.T) {
	t.Parallel()

	_, authErr := mapauth.New()

	if authErr == nil {
		t.Errorf("expected error on creation, but got none")
	}

	if !errors.Is(authErr, mapauth.ErrNoAuthorizations) {
		t.Errorf("expected error to be ErrNoAuthorizations, but got: %v", authErr)
	}
}

func TestUninitializedAuthenticator(t *testing.T) {
	t.Parallel()

	var auth *mapauth.MapAuthenticator

	_, authErr := auth.Authenticate("top", "secret")

	if authErr == nil {
		t.Errorf("expected error on creation, but got none")
	}

	if !errors.Is(authErr, mapauth.ErrNotInitialized) {
		t.Errorf("expected error to be ErrUninitialized, but got: %v", authErr)
	}
}
