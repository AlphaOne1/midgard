// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package midgard

import (
	"net/http"

	"github.com/AlphaOne1/midgard/defs"
)

// StackMiddleware stacks the given middleware slice to generate a single combined middleware.
// The middleware at index 0 is the outermost, going step by step to the innermost,
// e. g. mw[0](mw[1](mw[2]())).
func StackMiddleware(mw []defs.Middleware) defs.Middleware {
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
func StackMiddlewareHandler(mw []defs.Middleware, final http.Handler) http.Handler {
	if len(mw) == 0 {
		return final
	}

	return StackMiddleware(mw)(final)
}
