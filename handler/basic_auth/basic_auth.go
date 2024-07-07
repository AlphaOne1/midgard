package basic_auth

import (
	"bytes"
	"encoding/base64"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/AlphaOne1/midgard/defs"
)

type Authenticator interface {
	Authenticate(username, password string) (bool, error)
}

type Handler struct {
	auth          Authenticator
	realm         string
	authRealmInfo string
	next          http.Handler
}

func (h *Handler) sendNoAuth(w http.ResponseWriter) {
	w.Header().Add("WWW-Authenticate", h.authRealmInfo)
	w.WriteHeader(http.StatusUnauthorized)

	if _, err := w.Write([]byte("unauthorized")); err != nil {
		slog.Error("could not write", slog.String("error", err.Error()))
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authInfo, headerPrefixOK := strings.CutPrefix(
		r.Header.Get("Authorization"),
		"Basic ")

	if !headerPrefixOK || len(authInfo) < 6 {
		h.sendNoAuth(w)
		return
	}

	decode, decodeErr := base64.StdEncoding.DecodeString(authInfo)

	if decodeErr != nil {
		slog.Debug("could not decode auth info",
			slog.String("error", decodeErr.Error()),
			slog.String("authInfo", authInfo))

		h.sendNoAuth(w)
		return
	}

	credentials := bytes.Split(decode, []byte(":"))

	if len(credentials) != 2 ||
		len(credentials[0]) == 0 ||
		len(credentials[1]) == 0 {
		h.sendNoAuth(w)
		return
	}

	username := string(credentials[0])
	password := string(credentials[1])

	hasAuth, authErr := h.auth.Authenticate(username, password)

	if authErr != nil {
		slog.Error("authentication error",
			slog.String("error", authErr.Error()),
			slog.String("user", username))
	}

	if !hasAuth {
		h.sendNoAuth(w)
		return
	}

	h.next.ServeHTTP(w, r)
}

func WithAuthenticator(auth Authenticator) func(h *Handler) error {
	return func(h *Handler) error {
		h.auth = auth

		return nil
	}
}

func WithRealm(realm string) func(h *Handler) error {
	return func(h *Handler) error {
		h.realm = realm

		return nil
	}
}

func New(options ...func(handler *Handler) error) (defs.Middleware, error) {
	handler := Handler{}

	for _, opt := range options {
		if err := opt(&handler); err != nil {
			return nil, err
		}
	}

	if handler.auth == nil {
		return nil, errors.New("no authenticator configured")
	}

	if handler.realm == "" {
		handler.realm = "Restricted"
	}

	handler.authRealmInfo = `Basic realm="` + handler.realm + `", charset="UTF-8"`

	return func(next http.Handler) http.Handler {
		handler.next = next
		return &handler
	}, nil
}
