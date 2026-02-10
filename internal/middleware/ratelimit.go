package middleware

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"

	"github.com/dev/api-feedbacks/pkg/response"
)

// visitors tracks rate limiters per IP address.
var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.RWMutex
)

// getVisitor returns the rate limiter for the given IP, creating one if needed.
func getVisitor(ip string, rps rate.Limit) *rate.Limiter {
	mu.RLock()
	limiter, exists := visitors[ip]
	mu.RUnlock()

	if exists {
		return limiter
	}

	mu.Lock()
	defer mu.Unlock()

	// Double-check after acquiring write lock
	if limiter, exists = visitors[ip]; exists {
		return limiter
	}

	limiter = rate.NewLimiter(rps, int(rps)*2)
	visitors[ip] = limiter
	return limiter
}

// RateLimit returns a middleware that limits requests per IP using a token bucket.
func RateLimit(rps int) func(http.Handler) http.Handler {
	limit := rate.Limit(rps)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limiter := getVisitor(r.RemoteAddr, limit)
			if !limiter.Allow() {
				response.Error(w, http.StatusTooManyRequests, "RATE_LIMITED", "too many requests", nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
