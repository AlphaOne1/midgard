package basic_auth

import "errors"

type MapAuthenticator struct {
	Auths map[string]string
}

func (m *MapAuthenticator) Authorize(username, password string) (bool, error) {
	if len(m.Auths) == 0 {
		return false, errors.New("no auths configured")
	}

	pass, passFound := m.Auths[username]

	return passFound && pass == password, nil
}
