// Copyright the midgard contributors.
// SPDX-License-Identifier: MPL-2.0
package access_log

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/AlphaOne1/midgard/defs"
	"github.com/AlphaOne1/midgard/handler/basic_auth"
)

// Handler holds the information necessary for the log
type Handler struct {
	log   *slog.Logger
	level slog.Level
	next  http.Handler
}

// accessLogging is the access logging middleware. It logs every request with its
// correlationID, the clients address, http method and accessed path.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	entries := []any{
		slog.String("client", r.RemoteAddr),
		slog.String("method", r.Method),
		slog.String("target", r.URL.Path),
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

	slog.Log(context.Background(), h.level, "access", entries...)

	h.next.ServeHTTP(w, r)
}

// WithLogger configures the logger to use.
func WithLogger(log *slog.Logger) func(h *Handler) error {
	return func(h *Handler) error {
		if log == nil {
			return errors.New("cannot configure with nil logger")
		}

		h.log = log

		return nil
	}
}

// WithLogLevel configures the log level to use with the logger.
func WithLogLevel(level slog.Level) func(h *Handler) error {
	return func(h *Handler) error {
		h.level = level

		return nil
	}
}

// New generates a new access logging middleware.
func New(options ...func(*Handler) error) (defs.Middleware, error) {
	h := &Handler{
		log:   slog.Default(),
		level: slog.LevelInfo,
	}

	for _, opt := range options {
		if err := opt(h); err != nil {
			return nil, err
		}
	}

	return func(next http.Handler) http.Handler {
		h.next = next
		return h
	}, nil
}
