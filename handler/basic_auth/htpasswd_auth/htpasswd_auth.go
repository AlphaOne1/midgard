// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package htpasswd_auth

import (
	"errors"
	"io"
	"os"

	"github.com/tg123/go-htpasswd"
)

// HTPassWDAuth holds the htpasswd relevant data
type HTPassWDAuth struct {
	auth *htpasswd.File
}

// Authorize checks if for a given username the password hash matches the
// one stored in the used htpasswd file.
func (a *HTPassWDAuth) Authorize(username, password string) (bool, error) {
	if a == nil {
		return false, errors.New("htpasswd auth not initialized")
	}

	if a.auth == nil {
		return false, errors.New("htpasswd not initialized")
	}

	return a.auth.Match(username, password), nil
}

// WithAuthInput configures the htpasswd file to be read from the
// given io.Reader.
func WithAuthInput(in io.Reader) func(a *HTPassWDAuth) error {
	return func(a *HTPassWDAuth) error {
		if in == nil {
			return errors.New("input is nil")
		}

		var err error

		a.auth, err = htpasswd.NewFromReader(in, htpasswd.DefaultSystems, nil)
		return err
	}
}

// WithAuthFile configures the htpasswd file to be read from the
// filesystem with the given name.
func WithAuthFile(fileName string) func(a *HTPassWDAuth) error {
	return func(a *HTPassWDAuth) error {
		if len(fileName) == 0 {
			return errors.New("input file name is necessary")
		}

		input, err := os.Open(fileName)

		if err != nil {
			return err
		}

		defer func() { _ = input.Close() }()

		return WithAuthInput(input)(a)
	}
}

// New creates a new htpasswd authenticator.
func New(options ...func(*HTPassWDAuth) error) (*HTPassWDAuth, error) {
	a := HTPassWDAuth{}

	for _, opt := range options {
		if err := opt(&a); err != nil {
			return nil, err
		}
	}

	if a.auth == nil {
		return nil, errors.New("htpasswd input is required")
	}

	return &a, nil
}
