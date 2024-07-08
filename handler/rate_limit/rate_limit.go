package rate_limit

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/AlphaOne1/midgard/defs"
)

// Limiter is the interface a limiter has to implement to be used in the rate
// limiter middleware.
type Limiter interface {
	Limit() bool
}

// Handler holds the internal rate limiter information.
type Handler struct {
	Limit Limiter
	Next  http.Handler
}

// ServeHTTP limits the requests using the internal Limiter.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h == nil {
		slog.Error("rate limiter not initialized")
		w.WriteHeader(http.StatusServiceUnavailable)
		if _, err := w.Write([]byte("service not available")); err != nil {
			slog.Error("failed to write response", slog.String("error", err.Error()))
		}
		return
	}

	if !h.Limit.Limit() {
		w.WriteHeader(http.StatusTooManyRequests)
		if _, err := w.Write([]byte("too many requests")); err != nil {
			slog.Error("failed to write response", slog.String("error", err.Error()))
		}
		return
	}

	h.Next.ServeHTTP(w, r)
}

// WithLimiter sets the Limiter to use.
func WithLimiter(l Limiter) func(h *Handler) error {
	return func(h *Handler) error {
		if l == nil {
			return errors.New("invalid limiter (nil)")
		}

		h.Limit = l

		return nil
	}
}

// New creates a new rate limiter middleware.
func New(options ...func(*Handler) error) (defs.Middleware, error) {
	h := Handler{}

	for _, opt := range options {
		if err := opt(&h); err != nil {
			return nil, err
		}
	}

	if h.Limit == nil {
		return nil, errors.New("invalid limiter (nil)")
	}

	return func(next http.Handler) http.Handler {
		h.Next = next
		return &h
	}, nil
}
