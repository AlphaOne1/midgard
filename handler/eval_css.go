// vim: set ts=8 sw=8 smartindent:

package handler

import (
	"net/http"
	"slices"
	"strings"

	"github.com/AlphaOne1/midgard"
)

// EvalCSSHandler is a middleware that sets up the cross site scripting circumvention headers.
type EvalCSSHandler struct {
	// Methods contains the allowed methods specific for CSS for the given handler.
	Methods string
	// Origins contains the allowed origins
	Origins []string
	// Next contains the next handler in the handler chain.
	Next http.Handler
}

// ServeHTTP sets up the client with the appropriate headers.
func (e EvalCSSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, originFound := r.Header["Origin"]

	if r.Method == "OPTIONS" || originFound {
		relevantOrigin := e.Origins[0]

		if relevantOrigin != "*" && len(e.Origins) > 1 {
			requestHost := r.URL.Host

			if portIdx := strings.LastIndex(requestHost, ":"); portIdx != -1 {
				requestHost = requestHost[:portIdx]
			}

			if slices.Contains(e.Origins, requestHost) {
				relevantOrigin = requestHost
			}
		}

		w.Header().Set("Access-Control-Allow-Origin", relevantOrigin)
	}

	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", e.Methods)
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(
			[]string{
				"Accept",
				"Content-Type",
				"Content-Length",
				"Accept-Encoding",
				"X-CSRF-Token",
				"Authorization",
			},
			", "))

		return
	}

	e.Next.ServeHTTP(w, r)
}

// GenerateEvalCSSHandler sets up the cross site scripting circumvention disable headers.
func NewEvalCSSHandler(methods []string, origins []string) midgard.Middleware {
	if len(origins) == 0 || slices.Contains(origins, "*") {
		origins = []string{"*"}
	}

	return func(next http.Handler) http.Handler {
		return EvalCSSHandler{
			Methods: strings.Join(methods, ", "),
			Origins: origins,
			Next:    next,
		}
	}
}
