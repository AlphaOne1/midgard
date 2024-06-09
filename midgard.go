// vim: set ts=8 sw=8 smartindent:

package midgard

import (
	"net/http"
)

// Middleware represents the common type of http middleware.
// The idea is to have a common interface for all types of middlewares, that is, they get an
// input handler and return an output handler, that is extended by the middlewares functionality.
// Customization is done in generator functions, that take parameters to modify the behaviour of
// the final http.Handler, e.g. methods to allow.
type Middleware func(http.Handler) http.Handler

// ToStandard converts a Middleware to the standard interface.
// As the Middleware type just exists to improve readability, this conversion is without losses.
func ToStandard(mw Middleware) func(http.Handler) http.Handler {
	return (func(http.Handler) http.Handler)(mw)
}

// FromStandard converts a standard middleware to the Middleware type.
// As the Middleware type just exists to improve readability, this conversion is without losses.
func FromStandard(s func(http.Handler) http.Handler) Middleware {
	return Middleware(s)
}

// StackMiddleware stacks the given middleware slice to generate a single combined middleware.
// The middleware at index 0 is the outermost, going step by step to the innermost,
// e. g. mw[0](mw[1](mw[2]())).
func StackMiddleware(mw []Middleware) Middleware {
	switch l := len(mw); {
	case l < 1:
		return nil
	case l == 1:
		return mw[0]
	default:
		return func(next http.Handler) http.Handler {
			return mw[0](StackMiddleware(mw[1:])(next))
		}
	}
}

// StackMiddlewareHandler calls StackMiddleware on mw and applies it to the handler final.
func StackMiddlewareHandler(mw []Middleware, final http.Handler) http.Handler {
	if len(mw) == 0 {
		return final
	}

	return StackMiddleware(mw)(final)
}
