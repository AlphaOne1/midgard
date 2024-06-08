package correlation

import (
	"github.com/AlphaOne1/midgard"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func New() midgard.Middleware {
	return correlation
}

func getOrCreateID(id string) string {
	if len(id) > 0 {
		return id
	}

	newID := "n/a"

	if newUuid, err := uuid.NewRandom(); err == nil {
		newID = newUuid.String()
	}

	return newID
}

func correlation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := r.Header.Get("X-Correlation-ID")

		if correlationID == "" {
			tmp := getOrCreateID("")

			r.Header.Set("X-Correlation-ID", tmp)
			w.Header().Set("X-Correlation-ID", tmp)

			slog.Debug("created new correlation id", slog.String("correlation_id", tmp))
		} else {
			w.Header().Set("X-Correlation-ID", correlationID)
		}

		next.ServeHTTP(w, r)
	})
}
