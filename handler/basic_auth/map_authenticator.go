package basic_auth

type MapAuthenticator struct {
	Auths map[string]string
}

func (m *MapAuthenticator) Authorize(username, password string) (bool, error) {
	pass, passFound := m.Auths[username]

	return passFound && pass == password, nil
}
