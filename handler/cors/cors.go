// Copyright the midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package cors

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/AlphaOne1/midgard/defs"
	"github.com/AlphaOne1/midgard/util"
)

// Handler is a middleware that sets up the cross site scripting circumvention headers.
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

func (h *Handler) GetMWBase() *defs.MWBase {
	if h == nil {
		return nil
	}

	return &h.MWBase
}

var minimumAllowHeaders = []string{
	"Accept",
	"Accept-Encoding",
	"Authorization",
	"Content-Length",
	"Content-Type",
	"Origin",
	"User-Agent",
	"X-CSRF-Token",
}

// MinimumAllowHeaders returns a minimal list of headers, that should not do
// harm. It can be used to limit the allowed headers to a reasonable small set.
func MinimumAllowHeaders() []string {
	return append(make([]string, 0, len(minimumAllowHeaders)), minimumAllowHeaders...)
}

// relevantOrigin gets the origin that the client matches with the allowed origins.
// If there is no match or there are no origins set, an error is returned.
func relevantOrigin(origin []string, allowed []string) (string, error) {
	if len(allowed) == 1 && allowed[0] == "*" {
		return "*", nil
	}

	if len(origin) == 0 {
		return "", fmt.Errorf("no origin in header")
	}

	for _, orig := range origin {
		if len(orig) == 0 {
			continue
		}

		if slices.Contains(allowed, orig) {
			return orig, nil
		}
	}

	return "", fmt.Errorf("origin not allowed")
}

// ServeHTTP sets up the client with the appropriate headers.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !util.IntroCheck(h, w, r) {
		return
	}

	origin := r.Header["Origin"]

	relevantOrigin, roErr := relevantOrigin(origin, h.Origins)

	// no relevant origin found in request
	if roErr != nil {
		util.WriteState(w, h.Log(), http.StatusForbidden)
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
		util.WriteState(w, h.Log(), http.StatusMethodNotAllowed)
		return
	}

	// if there are headers configured, check the headers of the request and
	// disallow in case there are non-configured ones
	if len(h.Headers) > 0 {
		for hdr := range r.Header {
			if !h.Headers[strings.ToLower(hdr)] {
				util.WriteState(w, h.Log(), http.StatusForbidden)
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

// New sets up the cross site scripting circumvention disable headers.
// If no methods are specified, all methods are allowed.
// If no headers are specified, all headers are allowed.
// If origin contains "*" or is empty, the allowed origins are set to *.
func New(options ...func(handler *Handler) error) (defs.Middleware, error) {
	h := Handler{}

	for _, opt := range options {
		if opt == nil {
			return nil, errors.New("options cannot be nil")
		}

		if err := opt(&h); err != nil {
			return nil, err
		}
	}

	// if no origins are specified or one of the specified allowed origins is *
	// just set the origins to *
	if len(h.Origins) == 0 || slices.Contains(h.Origins, "*") {
		_ = WithOrigins([]string{"*"})(&h)
	}

	return func(next http.Handler) http.Handler {
		if err := h.SetNext(next); err != nil {
			return nil
		}
		return &h
	}, nil
}
