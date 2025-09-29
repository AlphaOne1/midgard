// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

// Package add_header provides a middleware for adding headers to HTTP responses.
package add_header

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/AlphaOne1/midgard/defs"
	"github.com/AlphaOne1/midgard/util"
)

// Handler holds the information of the added headers.
type Handler struct {
	defs.MWBase

	headers [][2]string
}

// GetMWBase returns the MWBase instance of the handler.
func (h *Handler) GetMWBase() *defs.MWBase {
	if h == nil {
		return nil
	}

	return &h.MWBase
}

// ServeHTTP handles the requests, adding the additionally provided headers to the responses.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !util.IntroCheck(h, w, r) {
		return
	}

	for _, v := range h.headers {
		w.Header().Set(v[0], v[1])
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

// WithHeaders configures the headers to add to responses.
func WithHeaders(headers [][2]string) func(*Handler) error {
	return func(h *Handler) error {
		h.headers = append(h.headers, headers...)

		return nil
	}
}

// New generates a new header adding middleware.
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
