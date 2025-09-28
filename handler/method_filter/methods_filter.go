// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package method_filter

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/AlphaOne1/midgard/defs"
	"github.com/AlphaOne1/midgard/util"
)

// Handler only lets configured HTTP methods pass.
type Handler struct {
	defs.MWBase
	// Methods contains methods whitelist for the endpoint.
	Methods map[string]bool
}

func (h *Handler) GetMWBase() *defs.MWBase {
	if h == nil {
		return nil
	}

	return &h.MWBase
}

// ServeHTTP denies access (405) if the method is not in the whitelist.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !util.IntroCheck(h, w, r) {
		return
	}

	if h.Methods == nil {
		h.Log().Error("method filter not initialized")
		util.WriteState(w, h.Log(), http.StatusServiceUnavailable)
		return
	}

	if h.Methods[r.Method] {
		h.Next().ServeHTTP(w, r)
	}

	util.WriteState(w, h.Log(), http.StatusMethodNotAllowed)
}

// WithMethods sets the methods_filter configuration to allow the given methods to pass. If used multiple times,
// the allowed methods of the different calls are all enabled.
func WithMethods(methods []string) func(m *Handler) error {
	return func(m *Handler) error {
		if m.Methods == nil {
			m.Methods = make(map[string]bool, len(methods))
		}

		for _, v := range methods {
			m.Methods[v] = true
		}

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

// New sets up the method filter middleware. Its parameters are functions manipulating an internal Config variable.
func New(options ...func(m *Handler) error) (defs.Middleware, error) {
	h := Handler{}

	for _, opt := range options {
		if opt == nil {
			return nil, errors.New("options cannot be nil")
		}

		if err := opt(&h); err != nil {
			return nil, err
		}
	}

	return func(next http.Handler) http.Handler {
		if err := h.SetNext(next); err != nil {
			return nil
		}
		return &h
	}, nil
}
