package htpasswd_auth

import (
	"testing"

	"github.com/AlphaOne1/midgard/util"
)

func TestHtpasswdAuth(t *testing.T) {
	tests := []struct {
		Username string
		Password string
		Valid    bool
	}{
		{
			Username: "user0",
			Password: "pass0",
			Valid:    true,
		},
		{
			Username: "user1",
			Password: "pass1",
			Valid:    true,
		},
		{
			Username: "user0",
			Password: "wrong0",
			Valid:    false,
		},
	}

	a := util.Must(New(WithAuthFile("testwd")))

	for k, v := range tests {
		gotAuth, gotErr := a.Authorize(v.Username, v.Password)

		if gotErr != nil {
			t.Errorf("%v: got error, but did not expect any: %v", k, gotErr)
		}

		if gotAuth != v.Valid {
			t.Errorf("%v: got auth %v but wanted %v", k, v.Valid, gotAuth)
		}
	}
}

func TestHtpasswdNil(t *testing.T) {
	var subject *HTPassWDAuth = nil

	if _, err := subject.Authorize("u", "p"); err == nil {
		t.Errorf("authorize on nil authorizer should give error")
	}
}

func TestHtpasswdNoFile(t *testing.T) {
	_, err := New(WithAuthFile("IDoNotExistNowhereInThisWorldForgetIt"))

	if err == nil {
		t.Errorf("authorizer initialization with non-existent file should give error")
	}
}
