// Copyright the midgard contributors.
// SPDX-License-Identifier: MPL-2.0
package method_filter

import (
	"log/slog"
	"net/http"

	"github.com/AlphaOne1/midgard/defs"
)

// Handler only lets configured HTTP methods pass.
type Handler struct {
	// Methods contains methods whitelist for the endpoint.
	Methods map[string]bool
	// Next contains the next handler in the handler chain.
	Next http.Handler
}

// ServeHTTP denies access (405) if the method is not in the whitelist.
func (m *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m == nil || m.Methods == nil {
		slog.Error("method filter not initialized")
		w.WriteHeader(http.StatusServiceUnavailable)
		if _, err := w.Write([]byte("service not available")); err != nil {
			slog.Error("failed to write response", slog.String("error", err.Error()))
		}
		return
	}

	if r != nil && m.Methods[r.Method] {
		m.Next.ServeHTTP(w, r)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusMethodNotAllowed)

		if _, err := w.Write([]byte("method not allowed")); err != nil {
			slog.Error("could not write response",
				slog.String("error", err.Error()))
		}
	}
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

// New sets up the method filter middleware. Its parameters are functions manipulating an internal Config variable.
func New(options ...func(m *Handler) error) (defs.Middleware, error) {
	h := Handler{}

	for _, opt := range options {
		if err := opt(&h); err != nil {
			return nil, err
		}
	}

	return func(next http.Handler) http.Handler {
		h.Next = next
		return &h
	}, nil
}
