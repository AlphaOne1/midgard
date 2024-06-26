package map_auth

import "errors"

type MapAuthenticator struct {
	auths map[string]string
}

func (a *MapAuthenticator) Authorize(username, password string) (bool, error) {

	if len(a.auths) == 0 {
		return false, errors.New("no auths configured")
	}

	pass, passFound := a.auths[username]

	return passFound && pass == password, nil
}

func WithAuths(auths map[string]string) func(a *MapAuthenticator) error {
	return func(a *MapAuthenticator) error {
		if len(auths) == 0 {
			return errors.New("no authorizations configured")
		}

		for k, v := range auths {
			a.auths[k] = v
		}

		return nil
	}
}

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
