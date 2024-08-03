// Copyright the midgard contributors.
// SPDX-License-Identifier: MPL-2.0
package map_auth

import (
	"testing"

	"github.com/AlphaOne1/midgard/util"
)

func TestMapAuthenticator(t *testing.T) {
	tests := []struct {
		Auths     map[string]string
		User      string
		Pass      string
		Want      bool
		WantError bool
	}{
		{
			Auths:     map[string]string{"testuser": "testpass"},
			User:      "testuser",
			Pass:      "testpass",
			Want:      true,
			WantError: false,
		},
		{
			Auths:     map[string]string{"testuser": "testpass"},
			User:      "testuser",
			Pass:      "testwrong",
			Want:      false,
			WantError: false,
		},
		{
			Auths:     map[string]string{},
			User:      "testuser",
			Pass:      "testpass",
			Want:      false,
			WantError: true,
		},
	}

	for k, v := range tests {
		auth := util.Must(New(WithAuths(v.Auths)))

		gotAuth, gotErr := auth.Authorize(v.User, v.Pass)

		if gotErr != nil {
			if !v.WantError {
				t.Errorf("%v: did not expect error, but got: %v", k, gotErr)
			}
			if gotAuth {
				t.Errorf("%v: got error, so auth should not work, but got: %v", k, gotAuth)
			}
		} else {
			if v.WantError {
				t.Errorf("%v: did expect error, but got none", k)
			}
			if gotAuth != v.Want {
				t.Errorf("%v: got auth %v but wanted %v", k, gotAuth, v.Want)
			}
		}
	}
}
