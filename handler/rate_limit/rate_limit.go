// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package rate_limit

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/AlphaOne1/midgard/defs"
	"github.com/AlphaOne1/midgard/util"
)

// Limiter is the interface a limiter has to implement to be used in the rate
// limiter middleware.
type Limiter interface {
	Limit() bool
}

// Handler holds the internal rate limiter information.
type Handler struct {
	defs.MWBase
	Limit Limiter
}

func (h *Handler) GetMWBase() *defs.MWBase {
	if h == nil {
		return nil
	}

	return &h.MWBase
}

// ServeHTTP limits the requests using the internal Limiter.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !util.IntroCheck(h, w, r) {
		return
	}

	if !h.Limit.Limit() {
		util.WriteState(w, h.Log(), http.StatusTooManyRequests)
		return
	}

	h.Next().ServeHTTP(w, r)
}

// WithLimiter sets the Limiter to use.
func WithLimiter(l Limiter) func(h *Handler) error {
	return func(h *Handler) error {
		if l == nil {
			return errors.New("invalid limiter (nil)")
		}

		h.Limit = l

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

// New creates a new rate limiter middleware.
func New(options ...func(*Handler) error) (defs.Middleware, error) {
	h := Handler{}

	for _, opt := range options {
		if opt == nil {
			return nil, errors.New("options cannot be nil")
		}

		if err := opt(&h); err != nil {
			return nil, err
		}
	}

	if h.Limit == nil {
		return nil, errors.New("invalid limiter (nil)")
	}

	return func(next http.Handler) http.Handler {
		if err := h.SetNext(next); err != nil {
			return nil
		}
		return &h
	}, nil
}
