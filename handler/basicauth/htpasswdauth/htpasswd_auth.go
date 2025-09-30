// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

// Package htpasswdauth implements the basic auth functionality using a htpasswd file.
package htpasswdauth

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/tg123/go-htpasswd"
)

// ErrEmptyInput is returned when the input is empty.
var ErrEmptyInput = errors.New("input is empty")

// ErrNotInitialized is returned when the htpasswd authenticator is not initialized.
var ErrNotInitialized = errors.New("htpasswd auth not initialized")

// HTPassWDAuth holds the htpasswd relevant data.
type HTPassWDAuth struct {
	auth *htpasswd.File
}

// Authorize checks if for a given username the password hash matches the
// one stored in the used htpasswd file.
func (a *HTPassWDAuth) Authorize(username, password string) (bool, error) {
	if a == nil {
		return false, ErrNotInitialized
	}

	if a.auth == nil {
		return false, ErrNotInitialized
	}

	return a.auth.Match(username, password), nil
}

// WithAuthInput configures the htpasswd file to be read from the
// given io.Reader.
func WithAuthInput(in io.Reader) func(a *HTPassWDAuth) error {
	return func(a *HTPassWDAuth) error {
		if in == nil {
			return ErrEmptyInput
		}

		var err error

		a.auth, err = htpasswd.NewFromReader(in, htpasswd.DefaultSystems, nil)

		return fmt.Errorf("could not read htpasswd input: %w", err)
	}
}

// WithAuthFile configures the htpasswd file to be read from the
// filesystem with the given name.
func WithAuthFile(fileName string) func(a *HTPassWDAuth) error {
	return func(auth *HTPassWDAuth) error {
		if len(fileName) == 0 {
			return ErrEmptyInput
		}

		input, err := os.Open(filepath.Clean(fileName))

		if err != nil {
			return fmt.Errorf("could not open auth file: %w", err)
		}

		defer func() { _ = input.Close() }()

		return WithAuthInput(input)(auth)
	}
}

// New creates a new htpasswd authenticator.
func New(options ...func(*HTPassWDAuth) error) (*HTPassWDAuth, error) {
	auth := HTPassWDAuth{}

	for _, opt := range options {
		if err := opt(&auth); err != nil {
			return nil, err
		}
	}

	if auth.auth == nil {
		return nil, ErrEmptyInput
	}

	return &auth, nil
}
