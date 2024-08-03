// Copyright the midgard contributors.
// SPDX-License-Identifier: MPL-2.0
package defs

import "net/http"

// Middleware represents the common type of http middleware.
// The idea is to have a common interface for all types of middlewares, that is, they get an
// input handler and return an output handler, that is extended by the middlewares functionality.
// Customization is done in generator functions, that take parameters to modify the behaviour of
// the final http.Handler, e.g. methods to allow.
type Middleware func(http.Handler) http.Handler
