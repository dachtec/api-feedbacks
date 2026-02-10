package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/dev/api-feedbacks/internal/middleware"
	"github.com/dev/api-feedbacks/pkg/response"
)

// NewRouter creates and configures the HTTP router with all routes and middleware.
func NewRouter(feedbackHandler *FeedbackHandler, apiKey string) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Recovery)
	r.Use(middleware.Logger)
	r.Use(middleware.CORS)
	r.Use(middleware.RateLimit(100))

	// Public routes (no auth required)
	r.Get("/health", healthCheck)
	r.Get("/ready", readinessCheck)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.APIKeyAuth(apiKey))

		r.Route("/api/v1/feedbacks", func(r chi.Router) {
			r.Post("/", feedbackHandler.Create)
			r.Get("/", feedbackHandler.List)
			r.Get("/{id}", feedbackHandler.GetByID)
			r.Patch("/{id}", feedbackHandler.Update)
		})
	})

	return r
}

// healthCheck handles GET /health
func healthCheck(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// readinessCheck handles GET /ready
func readinessCheck(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"status": "ready"})
}
