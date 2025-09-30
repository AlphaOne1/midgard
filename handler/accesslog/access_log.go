// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

// Package accesslog provides a middleware that logs every request.
package accesslog

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/AlphaOne1/midgard/defs"
	"github.com/AlphaOne1/midgard/handler/basicauth"
	"github.com/AlphaOne1/midgard/helper"
)

// ErrNilOption is returned when an option is nil.
var ErrNilOption = errors.New("option cannot be nil")

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
	if !helper.IntroCheck(h, w, r) {
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
		username, _, userFound, _ := basicauth.ExtractUserPass(authLine)

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
	handler := new(Handler)

	for _, opt := range options {
		if opt == nil {
			return nil, ErrNilOption
		}

		if err := opt(handler); err != nil {
			return nil, err
		}
	}

	return func(next http.Handler) http.Handler {
		if err := handler.SetNext(next); err != nil {
			return nil
		}

		return handler
	}, nil
}
