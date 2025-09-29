// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package map_auth_test

import (
	"fmt"
	"testing"

	"github.com/AlphaOne1/midgard/handler/basic_auth/map_auth"
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

	for k, v := range tests {
		t.Run(fmt.Sprintf("TestMapAuthenticator_%v", k), func(t *testing.T) {
			t.Parallel()

			auth, newErr := map_auth.New(map_auth.WithAuths(v.Auths))

			if newErr != nil {
				if !v.WantNewErr {
					t.Errorf("%v: got error on creation, but wanted none", k)
				}

				return
			}

			if v.WantNewErr {
				t.Errorf("%v: wanted error on creation, but got none", k)

				return
			}

			gotAuth, gotErr := auth.Authenticate(v.User, v.Pass)

			if gotErr != nil {
				if !v.WantError {
					t.Errorf("did not expect error, but got: %v", gotErr)
				}
				if gotAuth {
					t.Errorf("got error, so auth should not work, but got: %v", gotAuth)
				}
			} else {
				if v.WantError {
					t.Errorf("did expect error, but got none")
				}
				if gotAuth != v.Want {
					t.Errorf("got auth %v but wanted %v", gotAuth, v.Want)
				}
			}
		})
	}
}
