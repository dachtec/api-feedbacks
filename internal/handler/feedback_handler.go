package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/dev/api-feedbacks/internal/domain"
	"github.com/dev/api-feedbacks/internal/repository"
	"github.com/dev/api-feedbacks/internal/service"
	"github.com/dev/api-feedbacks/pkg/response"
)

// timeFormat defines the output format for timestamps (RFC3339 without milliseconds).
const timeFormat = time.RFC3339

// FeedbackHandler handles HTTP requests for feedback operations.
type FeedbackHandler struct {
	svc service.FeedbackService
}

// NewFeedbackHandler creates a new feedback handler.
func NewFeedbackHandler(svc service.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{svc: svc}
}

// Create handles POST /api/v1/feedbacks
func (h *FeedbackHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateFeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "invalid request body", nil)
		return
	}

	feedback := &domain.Feedback{
		UserID:       req.UserID,
		FeedbackType: domain.FeedbackType(req.FeedbackType),
		Rating:       req.Rating,
		Comment:      req.Comment,
	}

	created, err := h.svc.Create(r.Context(), feedback)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, toFeedbackResponse(created), nil)
}

// GetByID handles GET /api/v1/feedbacks/{id}
func (h *FeedbackHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "MISSING_ID", "feedback id is required", nil)
		return
	}

	feedback, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, toFeedbackResponse(feedback), nil)
}

// Update handles PATCH /api/v1/feedbacks/{id}
func (h *FeedbackHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "MISSING_ID", "feedback id is required", nil)
		return
	}

	var req UpdateFeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "invalid request body", nil)
		return
	}

	if req.FeedbackType == nil && req.Rating == nil && req.Comment == nil {
		response.Error(w, http.StatusBadRequest, "EMPTY_UPDATE", "at least one field must be provided for update", nil)
		return
	}

	updates := &service.FeedbackUpdateInput{
		FeedbackType: req.FeedbackType,
		Rating:       req.Rating,
		Comment:      req.Comment,
	}

	updated, err := h.svc.Update(r.Context(), id, updates)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, toFeedbackResponse(updated), nil)
}

// List handles GET /api/v1/feedbacks
func (h *FeedbackHandler) List(w http.ResponseWriter, r *http.Request) {
	filter, err := parseFilter(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_FILTER", err.Error(), nil)
		return
	}

	feedbacks, total, err := h.svc.List(r.Context(), filter)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	items := make([]FeedbackResponse, len(feedbacks))
	for i, f := range feedbacks {
		items[i] = toFeedbackResponse(f)
	}

	meta := ListMeta{
		Total:  total,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	}

	response.Success(w, http.StatusOK, items, meta)
}

// parseFilter extracts filter parameters from query string.
func parseFilter(r *http.Request) (repository.FeedbackFilter, error) {
	q := r.URL.Query()
	filter := repository.FeedbackFilter{
		Limit:  repository.DefaultLimit,
		Offset: 0,
	}

	if v := q.Get("user_id"); v != "" {
		filter.UserID = &v
	}

	if v := q.Get("feedback_type"); v != "" {
		if !domain.IsValidFeedbackType(v) {
			return filter, errors.New("invalid feedback_type filter")
		}
		filter.FeedbackType = &v
	}

	if v := q.Get("min_rating"); v != "" {
		val, err := strconv.Atoi(v)
		if err != nil || val < 1 || val > 5 {
			return filter, errors.New("min_rating must be an integer between 1 and 5")
		}
		filter.MinRating = &val
	}

	if v := q.Get("max_rating"); v != "" {
		val, err := strconv.Atoi(v)
		if err != nil || val < 1 || val > 5 {
			return filter, errors.New("max_rating must be an integer between 1 and 5")
		}
		filter.MaxRating = &val
	}

	if v := q.Get("created_from"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return filter, errors.New("created_from must be in RFC3339 format (e.g., 2026-01-01T00:00:00Z)")
		}
		filter.CreatedFrom = &t
	}

	if v := q.Get("created_to"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return filter, errors.New("created_to must be in RFC3339 format (e.g., 2026-02-01T00:00:00Z)")
		}
		filter.CreatedTo = &t
	}

	if v := q.Get("limit"); v != "" {
		val, err := strconv.Atoi(v)
		if err != nil || val < 1 {
			return filter, errors.New("limit must be a positive integer")
		}
		filter.Limit = val
	}

	if v := q.Get("offset"); v != "" {
		val, err := strconv.Atoi(v)
		if err != nil || val < 0 {
			return filter, errors.New("offset must be a non-negative integer")
		}
		filter.Offset = val
	}

	return filter, nil
}

// handleServiceError maps domain/service errors to HTTP responses.
func handleServiceError(w http.ResponseWriter, err error) {
	if ve, ok := domain.IsValidationError(err); ok {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "validation failed", ve.Errors)
		return
	}

	if errors.Is(err, domain.ErrNotFound) {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "resource not found", nil)
		return
	}

	if errors.Is(err, domain.ErrConflict) {
		response.Error(w, http.StatusConflict, "CONFLICT", "resource conflict", nil)
		return
	}

	response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "an internal error occurred", nil)
}

// toFeedbackResponse converts a domain Feedback to a FeedbackResponse DTO.
func toFeedbackResponse(f *domain.Feedback) FeedbackResponse {
	return FeedbackResponse{
		FeedbackID:   f.FeedbackID,
		UserID:       f.UserID,
		FeedbackType: string(f.FeedbackType),
		Rating:       f.Rating,
		Comment:      f.Comment,
		CreatedAt:    f.CreatedAt.UTC().Format(timeFormat),
		UpdatedAt:    f.UpdatedAt.UTC().Format(timeFormat),
	}
}
