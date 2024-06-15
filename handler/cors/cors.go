// vim: set ts=8 sw=8 smartindent:

package cors

import (
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/AlphaOne1/midgard/defs"
)

// Handler is a middleware that sets up the cross site scripting circumvention headers.
type Handler struct {
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
	// Next contains the next handler in the handler chain.
	Next http.Handler
}

var minimumAllowHeaders = []string{
	"Accept",
	"Accept-Encoding",
	"Authorization",
	"Content-Length",
	"Content-Type",
	"Origin",
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
func (e Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origin, _ := r.Header["Origin"]

	relevantOrigin, roErr := relevantOrigin(origin, e.Origins)

	if roErr != nil {
		w.WriteHeader(http.StatusForbidden)
		if _, err := fmt.Fprintf(w, "origin %v not allowed", origin); err != nil {
			slog.Error("could not write", slog.String("error", err.Error()))
		}
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", relevantOrigin)

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", e.MethodsReturn)
		w.Header().Set("Access-Control-Allow-Headers", e.HeadersReturn)

		w.WriteHeader(http.StatusOK)
		return
	}

	if len(e.Methods) > 0 && !e.Methods[r.Method] {
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := fmt.Fprintf(w, "method %s not allowed", r.Method); err != nil {
			slog.Error("could not write", slog.String("error", err.Error()))
		}
		return
	}

	if len(e.Headers) > 0 {
		for h := range r.Header {
			if !e.Headers[h] {
				w.WriteHeader(http.StatusForbidden)
				if _, err := fmt.Fprintf(w, "header %s not allowed", h); err != nil {
					slog.Error("could not write", slog.String("error", err.Error()))
				}
				return
			}
		}
	}

	e.Next.ServeHTTP(w, r)
}

// WithHeaders sets the allowed headers. If later a request contains headers that are not
// contained in this list, it will be denied the service.
func WithHeaders(headers []string) func(handler *Handler) error {
	return func(handler *Handler) error {
		headersMap := make(map[string]bool, len(headers))

		for _, h := range headers {
			headersMap[h] = true
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

// New sets up the cross site scripting circumvention disable headers.
// If no methods are specified, all methods are allowed.
// If no headers are specified, all headers are allowed.
// If origin contains "*" or is empty, the allowed origins are set to *.
func New(options ...func(handler *Handler) error) (defs.Middleware, error) {
	handler := Handler{}

	for _, opt := range options {
		if err := opt(&handler); err != nil {
			return nil, err
		}
	}

	// if no origins are specified or one of the specified allowed origins is *
	// just set the origins to *
	if len(handler.Origins) == 0 || slices.Contains(handler.Origins, "*") {
		if err := WithOrigins([]string{"*"})(&handler); err != nil {
			return nil, err
		}
	}

	return func(next http.Handler) http.Handler {
		handler.Next = next
		return handler
	}, nil
}
