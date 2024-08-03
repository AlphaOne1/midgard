// Copyright the midgard contributors.
// SPDX-License-Identifier: MPL-2.0
package util

import (
	"log/slog"
	"net/http"
	"os"
	"testing"
)

// Must exits the program if the given pair of function return and error contains an non-nil error value,
// otherwise the function return value val is returned.
func Must[T any](val T, err error) T {
	if err != nil {
		slog.Error("must-condition not met",
			slog.String("error", err.Error()))

		if !testing.Testing() {
			os.Exit(1)
		} else {
			slog.Info("not exiting due to test-mode")
		}
	}

	return val
}

// MapKeys gets the keys of the given map and saves them into a slice. There is not to expect a particular
// order of the elements of that slice. A nil map will produce a nil result, an empty map produces an empty
// slice.
func MapKeys[T comparable, S any](m map[T]S) []T {
	if m == nil {
		return nil
	}

	if len(m) == 0 {
		return make([]T, 0)
	}

	result := make([]T, 0, len(m))

	for k := range m {
		result = append(result, k)
	}

	return result
}

// DummyHandler is a handler used for internal testing.
// It simply writes the text "dummy" to the given http.ResponseWriter.
func DummyHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("dummy")); err != nil {
		slog.Error("could not write dummy", slog.String("error", err.Error()))
	}
	if r.Body != nil {
		if err := r.Body.Close(); err != nil {
			slog.Error("could not close request body", slog.String("error", err.Error()))
		}
	}
}
