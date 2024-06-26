package htpasswd_auth

import (
	"errors"
	"io"
	"os"

	"github.com/tg123/go-htpasswd"
)

type HTPassWDAuth struct {
	input io.Reader
	auth  *htpasswd.File
}

func (a *HTPassWDAuth) Authorize(username, password string) (bool, error) {
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

		a.input = in
		return nil
	}
}

func WithAuthFile(fileName string) func(a *HTPassWDAuth) error {
	return func(a *HTPassWDAuth) error {
		if len(fileName) == 0 {
			return errors.New("input file name is necessary")
		}

		var err error
		a.input, err = os.Open(fileName)

		if err != nil {
			return err
		}

		return nil
	}
}

func New(options ...func(*HTPassWDAuth) error) (*HTPassWDAuth, error) {
	a := HTPassWDAuth{}

	for _, opt := range options {
		if err := opt(&a); err != nil {
			return nil, err
		}
	}

	if a.input == nil {
		return nil, errors.New("htpasswd input is required")
	}

	var err error

	if a.auth, err = htpasswd.NewFromReader(a.input, htpasswd.DefaultSystems, nil); err != nil {
		return nil, err
	}

	return &a, nil
}
