package htpasswd_auth

import (
	"errors"
	"io"
	"os"

	"github.com/tg123/go-htpasswd"
)

type HTPassWDAuth struct {
	auth *htpasswd.File
}

func (a *HTPassWDAuth) Authorize(username, password string) (bool, error) {
	if a == nil {
		return false, errors.New("htpasswd auth not initialized")
	}

	if a.auth == nil {
		return false, errors.New("htpasswd: not initialized")
	}

	return a.auth.Match(username, password), nil
}

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
