// Copyright the midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package correlation

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/AlphaOne1/midgard/defs"
	"github.com/AlphaOne1/midgard/util"
)

type Handler struct {
	defs.MWBase
}

func (h *Handler) GetMWBase() *defs.MWBase {
	if h == nil {
		return nil
	}

	return &h.MWBase
}

// ServeHTTP is implements the correlation id enriching middleware.
// It adds an X-Correlation-ID header if none was present.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !util.IntroCheck(h, w, r) {
		return
	}

	correlationID := r.Header.Get("X-Correlation-ID")

	if correlationID == "" {
		tmp := util.GetOrCreateID("")

		r.Header.Set("X-Correlation-ID", tmp)
		w.Header().Set("X-Correlation-ID", tmp)

		h.Log().Debug("created new correlation id", slog.String("correlation_id", tmp))
	} else {
		w.Header().Set("X-Correlation-ID", correlationID)
	}

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

// New generates a new correlation id enriching middleware.
func New(options ...func(*Handler) error) (defs.Middleware, error) {
	h := &Handler{}

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
