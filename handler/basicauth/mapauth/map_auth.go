// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

// Package mapauth implements the basic auth functionality using a user-pass-map.
package mapauth

import "errors"

// ErrNoAuthorizations is returned when no authorizations are configured.
var ErrNoAuthorizations = errors.New("no authorizations configured")

// ErrNotInitialized is returned when the map authenticator is not initialized.
var ErrNotInitialized = errors.New("mapauth not initialized")

// MapAuthenticator holds the authentication relevant data.
type MapAuthenticator struct {
	auths map[string]string // map containing username-password pairs
}

// Authenticate checks if a given username has the given password entry
// identical in the internal auths map.
func (a *MapAuthenticator) Authenticate(username, password string) (bool, error) {
	if a == nil {
		return false, ErrNotInitialized
	}

	// This state should never happen, as New will not allow
	// returning a MapAuthenticator without authorizations.
	if len(a.auths) == 0 {
		return false, ErrNoAuthorizations
	}

	pass, passFound := a.auths[username]

	return passFound && pass == password, nil
}

// WithAuths sets the allowed username-password combinations.
func WithAuths(auths map[string]string) func(a *MapAuthenticator) error {
	return func(a *MapAuthenticator) error {
		if len(auths) == 0 {
			return ErrNoAuthorizations
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
	authenticator := MapAuthenticator{}

	for _, opt := range options {
		if err := opt(&authenticator); err != nil {
			return nil, err
		}
	}

	if len(authenticator.auths) == 0 {
		return nil, ErrNoAuthorizations
	}

	return &authenticator, nil
}
