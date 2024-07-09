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

// Authenticator is an interface the basic auth handler uses to check if the
// given credentials match an allowed entry.
type Authenticator interface {
	// Authenticate checks, if a given username and password are allowed
	// credentials
	Authenticate(username, password string) (bool, error)
}

// Handler holds the internal data of the basic authentication middleware.
type Handler struct {
	auth          Authenticator // auth holds the Authenticator used
	realm         string        // realm to report to the client
	authRealmInfo string        // authRealmInfo holds the response header
	next          http.Handler  // next handler in the middleware stack
}

// sendNoAuth sends the client that his credentials are not allowed
func (h *Handler) sendNoAuth(w http.ResponseWriter) {
	w.Header().Add("WWW-Authenticate", h.authRealmInfo)
	w.WriteHeader(http.StatusUnauthorized)

	if _, err := w.Write([]byte("unauthorized")); err != nil {
		slog.Error("could not write", slog.String("error", err.Error()))
	}
}

// ExtractUserPass extracts the username and the password out of the given header
// value for Authorization. It signalizes if the desired information exists or en
// error, when the auth string is unprocessable.
func ExtractUserPass(auth string) (user, pass string, found bool, err error) {

	authInfo, headerPrefixOK := strings.CutPrefix(auth, "Basic ")

	if !headerPrefixOK || len(authInfo) < 6 {
		return "", "", false, nil
	}

	decode, decodeErr := base64.StdEncoding.DecodeString(authInfo)

	if decodeErr != nil {
		slog.Debug("could not decode auth info",
			slog.String("error", decodeErr.Error()),
			slog.String("authInfo", authInfo))

		return "", "", false, errors.New("could not decode auth info")
	}

	credentials := bytes.Split(decode, []byte(":"))

	if len(credentials) != 2 ||
		len(credentials[0]) == 0 ||
		len(credentials[1]) == 0 {
		return "", "", false, nil
	}

	return string(credentials[0]), string(credentials[1]), true, nil
}

// ServeHTTP implements the basic auth functionality.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h == nil {
		slog.Error("basic auth not initialized")
		w.WriteHeader(http.StatusServiceUnavailable)
		if _, err := w.Write([]byte("service not available")); err != nil {
			slog.Error("failed to write response", slog.String("error", err.Error()))
		}
		return
	}

	authInfo := r.Header.Get("Authorization")
	username, password, authFound, authErr := ExtractUserPass(authInfo)

	if authErr != nil {
		slog.Debug("could not process auth info",
			slog.String("error", authErr.Error()),
			slog.String("authInfo", authInfo))
	}

	if !authFound {
		h.sendNoAuth(w)
		return
	}

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

// WithAuthenticator sets the Authenticator to use.
func WithAuthenticator(auth Authenticator) func(h *Handler) error {
	return func(h *Handler) error {
		h.auth = auth

		return nil
	}
}

// WithRealm sets the realm to use
func WithRealm(realm string) func(h *Handler) error {
	return func(h *Handler) error {
		h.realm = realm

		return nil
	}
}

// New generates a new basic authentication middleware.
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
