package access_log

import (
	"log/slog"
	"net/http"

	"github.com/AlphaOne1/midgard/defs"
)

// New generates a new access logging middleware.
func New() defs.Middleware {
	return accessLogging
}

// accessLogging is the access logging middleware. It logs every request with its
// correlationID, the clients address, http method and accessed path.
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
