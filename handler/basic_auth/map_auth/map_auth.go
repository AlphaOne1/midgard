// Copyright the midgard contributors.
// SPDX-License-Identifier: MPL-2.0
package map_auth

import "errors"

// MapAuthenticator holds the authentication relevant data
type MapAuthenticator struct {
	auths map[string]string // map containing username-password pairs
}

// Authorize checks if a given username has the given password entry
// identical in the internal auths map.
func (a *MapAuthenticator) Authenticate(username, password string) (bool, error) {

	if a == nil {
		return false, errors.New("map auth not initialized")
	}

	if len(a.auths) == 0 {
		return false, errors.New("no auths configured")
	}

	pass, passFound := a.auths[username]

	return passFound && pass == password, nil
}

// WithAuths sets the allowed username-password combinations.
func WithAuths(auths map[string]string) func(a *MapAuthenticator) error {
	return func(a *MapAuthenticator) error {
		if len(auths) == 0 {
			return errors.New("no authorizations configured")
		}

		if len(a.auths) == 0 {
			a.auths = make(map[string]string, len(auths))
		}

		for k, v := range auths {
			a.auths[k] = v
		}

		return nil
	}
}

// New creates a new MapAuthenticator with the given configuration.
func New(options ...func(a *MapAuthenticator) error) (*MapAuthenticator, error) {
	a := MapAuthenticator{}

	for _, opt := range options {
		if err := opt(&a); err != nil {
			return nil, err
		}
	}

	if len(a.auths) == 0 {
		return nil, errors.New("no auths configured")
	}

	return &a, nil
}
