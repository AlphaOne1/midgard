// vim: set ts=8 sw=8 smartindent:

package method_filter

import (
	"log/slog"
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
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusMethodNotAllowed)

		if _, err := w.Write([]byte("method not allowed")); err != nil {
			slog.Error("could not write response",
				slog.String("error", err.Error()))
		}
	}
}

type Config struct {
	methods []string
}

// WithMethods sets the methods_filter configuration to allow the given methods to pass. If used multiple times,
// the allowed methods of the different calls are all enabled.
func WithMethods(methods []string) func(c *Config) error {
	return func(c *Config) error {
		c.methods = append(c.methods, methods...)
		return nil
	}
}

// New sets up the method filter middleware. Its parameters are functions manipulating an internal Config variable.
func New(configs ...func(c *Config) error) (midgard.Middleware, error) {
	cfg := Config{}

	for _, c := range configs {
		if err := c(&cfg); err != nil {
			return nil, err
		}
	}

	return func(next http.Handler) http.Handler {
		return MethodsFilter{
			Methods: func() map[string]bool {
				result := make(map[string]bool, len(cfg.methods))

				for _, v := range cfg.methods {
					result[v] = true
				}

				return result
			}(),
			Next: next,
		}
	}, nil
}
