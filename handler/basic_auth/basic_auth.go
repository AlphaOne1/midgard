// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package basic_auth

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/AlphaOne1/midgard/defs"
	"github.com/AlphaOne1/midgard/util"
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
	defs.MWBase

	auth          Authenticator // auth holds the Authenticator used
	realm         string        // realm to report to the client
	authRealmInfo string        // authRealmInfo holds the response header
	redirect      string        // redirect address to authenticate
}

func (h *Handler) GetMWBase() *defs.MWBase {
	if h == nil {
		return nil
	}

	return &h.MWBase
}

// sendNoAuth sends the client that his credentials are not allowed.
func (h *Handler) sendNoAuth(w http.ResponseWriter, r *http.Request) {
	if len(h.redirect) > 0 {
		http.Redirect(w, r, h.redirect, http.StatusFound)
	} else {
		w.Header().Add("WWW-Authenticate", h.authRealmInfo)
		util.WriteState(w, h.Log(), http.StatusUnauthorized)
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
		return "", "", false, fmt.Errorf("could not decode auth info: %w", decodeErr)
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
	if !util.IntroCheck(h, w, r) {
		return
	}

	authInfo := r.Header.Get("Authorization")
	username, password, authFound, authErr := ExtractUserPass(authInfo)

	if authErr != nil {
		h.Log().Debug("could not process auth info",
			slog.String("error", authErr.Error()),
			slog.String("authInfo", authInfo))
	}

	if !authFound {
		h.sendNoAuth(w, r)

		return
	}

	hasAuth, authErr := h.auth.Authenticate(username, password)

	if authErr != nil {
		h.Log().Error("authentication error",
			slog.String("error", authErr.Error()),
			slog.String("user", username))
	}

	if !hasAuth {
		h.sendNoAuth(w, r)

		return
	}

	h.Next().ServeHTTP(w, r)
}

// WithAuthenticator sets the Authenticator to use.
func WithAuthenticator(auth Authenticator) func(h *Handler) error {
	return func(h *Handler) error {
		h.auth = auth

		return nil
	}
}

// WithRealm sets the realm to use.
func WithRealm(realm string) func(h *Handler) error {
	return func(h *Handler) error {
		h.realm = realm

		return nil
	}
}

func WithRedirect(redirect string) func(h *Handler) error {
	return func(h *Handler) error {
		h.redirect = redirect

		return nil
	}
}

// WithLogger configures the logger to use.
func WithLogger(log *slog.Logger) func(h *Handler) error {
	return defs.WithLogger[*Handler](log)
}

// WithLogLevel configures the log level to use with the logger.
func WithLogLevel(level slog.Level) func(h *Handler) error {
	return defs.WithLogLevel[*Handler](level)
}

// New generates a new basic authentication middleware.
func New(options ...func(handler *Handler) error) (defs.Middleware, error) {
	h := Handler{}

	for _, opt := range options {
		if opt == nil {
			return nil, errors.New("options cannot be nil")
		}

		if err := opt(&h); err != nil {
			return nil, err
		}
	}

	if h.auth == nil {
		return nil, errors.New("no authenticator configured")
	}

	if h.realm == "" {
		h.realm = "Restricted"
	}

	h.authRealmInfo = `Basic realm="` + h.realm + `", charset="UTF-8"`

	return func(next http.Handler) http.Handler {
		if err := h.SetNext(next); err != nil {
			return nil
		}

		return &h
	}, nil
}
