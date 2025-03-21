// Copyright the midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"log/slog"
	"net/http"
	"os"
	"reflect"

	"github.com/google/uuid"

	"github.com/AlphaOne1/midgard/defs"
)

// exitFunc is used to exit the program. For testing purposes it can be set to another function suitable
// for non-exiting tests. Do not touch, if insecure!
var exitFunc = os.Exit

// Must exits the program if the given pair of function return and error contains an non-nil error value,
// otherwise the function return value val is returned.
func Must[T any](val T, err error) T {
	if err != nil {
		slog.Error("must-condition not met",
			slog.String("error", err.Error()))

		exitFunc(1)
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

// GetOrCreateID generates a new uuid, if the given id is empty, otherwise the given id is returned.
func GetOrCreateID(id string) string {
	if len(id) > 0 {
		return id
	}

	newID := "n/a"

	if newUuid, err := uuid.NewRandom(); err == nil {
		newID = newUuid.String()
	}

	return newID
}

// WriteState sets the specified HTTP response code and writes the code specific text as body.
// If an error occurs on writing to the client, it is logged to the specified logging instance.
// It is intended to give error feedback to clients.
func WriteState(w http.ResponseWriter, log *slog.Logger, httpState int) {
	h := w.Header()

	h.Del("Content-Length")
	h.Set("Content-Type", "text/plain; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(httpState)

	if _, err := w.Write([]byte(http.StatusText(httpState))); err != nil {
		log.Error("failed to write response", slog.String("error", err.Error()))
	}
}

// IntroCheck is used to facilitate the introductory check in each handler for the basic requirements.
// It manages the corresponding logging operations and can be used as follows:
//
//	if !util.IntroCheck(h, w, r) {
//	    return
//	}
func IntroCheck(h defs.MWBaser, w http.ResponseWriter, r *http.Request) bool {
	if reflect.ValueOf(h).IsNil() {
		slog.Error("handler nil")
		WriteState(w, slog.Default(), http.StatusInternalServerError)
		return false
	}

	if r == nil {
		slog.Debug("request nil")
		WriteState(w, h.GetMWBase().Log(), http.StatusBadRequest)
		return false
	}

	return true
}

// DummyHandler is a handler used for internal testing.
// It simply writes the text "dummy" to the given http.ResponseWriter.
func DummyHandler(w http.ResponseWriter, _ /*r*/ *http.Request) {
	_, _ = w.Write([]byte("dummy"))
}
