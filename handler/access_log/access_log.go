package access_log

import (
	"github.com/AlphaOne1/midgard"
	"log/slog"
	"net/http"
)

func New() midgard.Middleware {
	return accessLogging
}

func accessLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := r.Header.Get("X-Correlation-ID")

		if correlationID == "" {
			correlationID = "unknown"
		}

		slog.Info("access",
			slog.String("correlation_id", correlationID),
			slog.String("client", r.RemoteAddr),
			slog.String("method", r.Method),
			slog.String("target", r.URL.Path),
		)

		next.ServeHTTP(w, r)
	})
}
