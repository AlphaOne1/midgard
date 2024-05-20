// vim: set ts=8 sw=8 smartindent:

package handler

import (
	"net/http"

	"github.com/AlphaOne1/midgard"
)

// MethodsFilter only lets configured HTTP methods pass.
type MethodsFilter struct {
	// Methods contains methods whitelist for the endpoint.
	Methods map[string]bool
	// Next contains the next handler in the handler chain.
	Next http.Handler
}

// ServeHTTP denies access (405) if the method is not in the whitelist.
func (m MethodsFilter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.Methods[r.Method] {
		m.Next.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
	}
}

// NewMethodsFilter sets up the method filter middlware. Its methods parameter contains
// the whitelist of allowed methods.
func NewMethodsFilter(methods []string) midgard.Middleware {
	return func(next http.Handler) http.Handler {
		return MethodsFilter{
			Methods: func() map[string]bool {
				result := make(map[string]bool, len(methods))

				for _, v := range methods {
					result[v] = true
				}

				return result
			}(),
			Next: next,
		}
	}
}
