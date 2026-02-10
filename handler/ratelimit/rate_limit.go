// SPDX-FileCopyrightText: 2026 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

// Package ratelimit provides middleware for rate limiting HTTP requests.
package ratelimit

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/AlphaOne1/midgard/defs"
	"github.com/AlphaOne1/midgard/helper"
)

// ErrInvalidLimiter is returned when the given limiter is invalid.
var ErrInvalidLimiter = errors.New("invalid limiter")

// ErrNilOption is returned when an option is nil.
var ErrNilOption = errors.New("option cannot be nil")

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

// GetMWBase returns the MWBase instance of the handler.
func (h *Handler) GetMWBase() *defs.MWBase {
	if h == nil {
		return nil
	}

	return &h.MWBase
}

// ServeHTTP limits the requests using the internal Limiter.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !helper.IntroCheck(h, w, r) {
		return
	}

	if !h.Limit.Limit() {
		helper.WriteState(w, h.Log(), http.StatusTooManyRequests)

		return
	}

	h.Next().ServeHTTP(w, r)
}

// WithLimiter sets the Limiter to use.
func WithLimiter(l Limiter) func(h *Handler) error {
	return func(h *Handler) error {
		if l == nil {
			return ErrInvalidLimiter
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
	handler := Handler{}

	for _, opt := range options {
		if opt == nil {
			return nil, ErrNilOption
		}

		if err := opt(&handler); err != nil {
			return nil, err
		}
	}

	if handler.Limit == nil {
		return nil, ErrInvalidLimiter
	}

	return func(next http.Handler) http.Handler {
		if err := handler.SetNext(next); err != nil {
			return nil
		}

		return &handler
	}, nil
}
