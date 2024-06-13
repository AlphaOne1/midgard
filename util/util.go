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
		}
	}

	return val
}

// DummyHandler is a handler used for internal testing.
// It simply writes the text "dummy" to the given http.ResponseWriter.
func DummyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("dummy"))
}
