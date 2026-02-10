package middleware

import (
	"net/http"

	"github.com/dev/api-feedbacks/pkg/response"
)

// APIKeyAuth returns a middleware that validates the X-API-Key header.
func APIKeyAuth(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-API-Key")
			if key == "" {
				response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing API key", nil)
				return
			}
			if key != apiKey {
				response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid API key", nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
