package access_log

import (
	"log/slog"
	"net/http"

	"github.com/AlphaOne1/midgard/defs"
	"github.com/AlphaOne1/midgard/handler/basic_auth"
)

// New generates a new access logging middleware.
func New() defs.Middleware {
	return accessLogging
}

// accessLogging is the access logging middleware. It logs every request with its
// correlationID, the clients address, http method and accessed path.
func accessLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		entries := []any{
			slog.String("client", r.RemoteAddr),
			slog.String("method", r.Method),
			slog.String("target", r.URL.Path),
		}

		if correlationID := r.Header.Get("X-Correlation-ID"); correlationID != "" {
			entries = append(entries, slog.String("correlation_id", correlationID))
		}

		if authLine := r.Header.Get("Authorization"); authLine != "" {
			username, _, userFound, _ := basic_auth.ExtractUserPass(authLine)

			if userFound {
				entries = append(entries, slog.String("user", username))
			}
		}

		slog.Info("access", entries...)

		next.ServeHTTP(w, r)
	})
}
