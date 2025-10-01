// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

// Package cors provides a middleware for handling CORS (Cross-Origin Resource Sharing) requests.
package cors

import (
	"errors"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/AlphaOne1/midgard/defs"
	"github.com/AlphaOne1/midgard/helper"
)

// ErrNilOption is returned when an option is nil.
var ErrNilOption = errors.New("option cannot be nil")

// ErrNoOrigin is returned when there is no origin in the request.
var ErrNoOrigin = errors.New("no origin in header")

// ErrOriginNotAllowed is returned when the origin is not allowed.
var ErrOriginNotAllowed = errors.New("origin not allowed")

// Handler is a middleware that sets up the cross-site scripting circumvention headers.
type Handler struct {
	defs.MWBase

	// Headers contains the allowed headers
	Headers map[string]bool
	// HeadersReturn contains the comma-concatenated allowed headers
	// as returned in the allow-header header
	HeadersReturn string
	// Methods contains the allowed methods specific for CSS for the given handler.
	Methods map[string]bool
	// MethodsReturn contains the comma-concatenated allowed methods
	// as returned in the allow-methods header
	MethodsReturn string
	// Origins contains the allowed origins
	Origins []string
}

// GetMWBase returns the MWBase instance of the handler.
func (h *Handler) GetMWBase() *defs.MWBase {
	if h == nil {
		return nil
	}

	return &h.MWBase
}

// MinimumAllowHeaders returns a minimal list of headers, that should not do
// harm. It can be used to limit the allowed headers to a reasonable small set.
func MinimumAllowHeaders() []string {
	return []string{
		"Accept",
		"Accept-Encoding",
		"Authorization",
		"Content-Length",
		"Content-Type",
		"Origin",
		"User-Agent",
		"X-CSRF-Token",
	}
}

// relevantOrigin gets the origin that the client matches with the allowed origins.
// If there is no match or there are no origins set, an error is returned.
func relevantOrigin(origin []string, allowed []string) (string, error) {
	if len(allowed) == 1 && allowed[0] == "*" {
		return "*", nil
	}

	if len(origin) == 0 {
		return "", ErrNoOrigin
	}

	for _, orig := range origin {
		if len(orig) == 0 {
			continue
		}

		if slices.Contains(allowed, orig) {
			return orig, nil
		}
	}

	return "", ErrOriginNotAllowed
}

// ServeHTTP sets up the client with the appropriate headers.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !helper.IntroCheck(h, w, r) {
		return
	}

	origin := r.Header["Origin"]

	relevantOrigin, roErr := relevantOrigin(origin, h.Origins)

	// no relevant origin found in the request
	if roErr != nil {
		helper.WriteState(w, h.Log(), http.StatusForbidden)

		return
	}

	w.Header().Set("Access-Control-Allow-Origin", relevantOrigin)

	// on OPTIONS request, just give the possible methods and headers
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", h.MethodsReturn)
		w.Header().Set("Access-Control-Allow-Headers", h.HeadersReturn)

		w.WriteHeader(http.StatusOK)

		return
	}

	// we have methods configured, but the request does not match any of them
	if len(h.Methods) > 0 && !h.Methods[r.Method] {
		helper.WriteState(w, h.Log(), http.StatusMethodNotAllowed)

		return
	}

	// if there are headers configured, check the headers of the request and
	// disallow in case there are non-configured ones
	if len(h.Headers) > 0 {
		for hdr := range r.Header {
			if !h.Headers[strings.ToLower(hdr)] {
				helper.WriteState(w, h.Log(), http.StatusForbidden)

				return
			}
		}
	}

	h.Next().ServeHTTP(w, r)
}

// WithHeaders sets the allowed headers. If later a request contains headers that are not
// contained in this list, it will be denied the service.
func WithHeaders(headers []string) func(handler *Handler) error {
	return func(handler *Handler) error {
		headersMap := make(map[string]bool, len(headers))

		for _, h := range headers {
			headersMap[strings.ToLower(h)] = true
		}

		handler.Headers = headersMap
		handler.HeadersReturn = strings.Join(headers, ", ")

		return nil
	}
}

// WithMethods sets the allowed methods. If later a request uses a method that are not
// contained in this list, it will be denied the service.
func WithMethods(methods []string) func(handler *Handler) error {
	return func(handler *Handler) error {
		methodsMap := make(map[string]bool, len(methods))

		for _, m := range methods {
			methodsMap[m] = true
		}

		handler.Methods = methodsMap
		handler.MethodsReturn = strings.Join(methods, ", ")

		return nil
	}
}

// WithOrigins sets the allowed origins. If later a comes from and origin that are not
// contained in this list, it will be denied the service. A special origin is "*", that
// is the wildcard for "all" origins.
func WithOrigins(origins []string) func(handler *Handler) error {
	return func(handler *Handler) error {
		handler.Origins = origins

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

// New sets up the cross-site scripting circumvention disable headers.
// If no methods are specified, all methods are allowed.
// If no headers are specified, all headers are allowed.
// If origin contains "*" or is empty, the allowed origins are set to *.
func New(options ...func(handler *Handler) error) (defs.Middleware, error) {
	handler := Handler{}

	for _, opt := range options {
		if opt == nil {
			return nil, ErrNilOption
		}

		if err := opt(&handler); err != nil {
			return nil, err
		}
	}

	// if no origins are specified or one of the specified allowed origins is *
	// just set the origins to *
	if len(handler.Origins) == 0 || slices.Contains(handler.Origins, "*") {
		_ = WithOrigins([]string{"*"})(&handler)
	}

	return func(next http.Handler) http.Handler {
		if err := handler.SetNext(next); err != nil {
			return nil
		}

		return &handler
	}, nil
}
