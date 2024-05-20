// vim: set ts=8 sw=8 smartindent:

package handler

import (
	"net/http"
	"strings"

	"github.com/AlphaOne1/midgard"
)

// EvalCSSHandler is a middleware that sets up the cross site scripting circumvention headers.
type EvalCSSHandler struct {
	// Methods contains the allowed methods specific for CSS for the given handler.
	Methods string
	// Next contains the next handler in the handler chain.
	Next http.Handler
}

// ServeHTTP sets up the client with the appropriate headers.
func (e EvalCSSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, originFound := r.Header["Origin"]

	if r.Method == "OPTIONS" || originFound {
		w.Header().Set("Access-Control-Allow-Origin", "*")
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
func NewEvalCSSHandler(methods []string) midgard.Middleware {
	return func(next http.Handler) http.Handler {
		return EvalCSSHandler{
			Methods: strings.Join(methods, ", "),
			Next:    next,
		}
	}
}
