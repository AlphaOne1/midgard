// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

// Package access_log provides a middleware that logs every request.
package access_log

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/AlphaOne1/midgard/defs"
	"github.com/AlphaOne1/midgard/handler/basic_auth"
	"github.com/AlphaOne1/midgard/util"
)

// Handler holds the information necessary for the log.
type Handler struct {
	defs.MWBase
}

// GetMWBase returns the MWBase instance of the handler.
func (h *Handler) GetMWBase() *defs.MWBase {
	if h == nil {
		return nil
	}

	return &h.MWBase
}

// ServeHTTP implements the access logging middleware. It logs every request with its
// correlationID, the client's address, http method and accessed path.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !util.IntroCheck(h, w, r) {
		return
	}

	entries := []any{
		slog.String("client", r.RemoteAddr),
		slog.String("method", r.Method),
	}

	if r.URL != nil {
		entries = append(entries, slog.String("target", r.URL.Path))
	}

	if correlationID := r.Header.Get("X-Correlation-ID"); correlationID != "" {
		entries = append(entries, slog.String("correlation_id", correlationID))
	}

	if authLine := r.Header.Get("Authorization"); authLine != "" {
		username, _, userFound, _ := basic_auth.ExtractUserPass(authLine)

		if userFound {
			entries = append(entries, slog.String("user", username))
		}
	}

	h.Log().Log(r.Context(), h.LogLevel(), "access", entries...)

	h.Next().ServeHTTP(w, r)
}

// WithLogger configures the logger to use.
func WithLogger(log *slog.Logger) func(h *Handler) error {
	return defs.WithLogger[*Handler](log)
}

// WithLogLevel configures the log level to use with the logger.
func WithLogLevel(level slog.Level) func(h *Handler) error {
	return defs.WithLogLevel[*Handler](level)
}

// New generates a new access logging middleware.
func New(options ...func(*Handler) error) (defs.Middleware, error) {
	h := new(Handler)

	for _, opt := range options {
		if opt == nil {
			return nil, errors.New("options cannot be nil")
		}

		if err := opt(h); err != nil {
			return nil, err
		}
	}

	return func(next http.Handler) http.Handler {
		if err := h.SetNext(next); err != nil {
			return nil
		}

		return h
	}, nil
}
