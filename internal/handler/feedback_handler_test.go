package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/dev/api-feedbacks/internal/domain"
	"github.com/dev/api-feedbacks/internal/repository"
	"github.com/dev/api-feedbacks/internal/service"
)

// mockService implements service.FeedbackService for handler tests.
type mockService struct {
	feedbacks map[string]*domain.Feedback
	idCounter int
}

func newMockService() *mockService {
	return &mockService{feedbacks: make(map[string]*domain.Feedback)}
}

func (m *mockService) Create(_ context.Context, f *domain.Feedback) (*domain.Feedback, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	m.idCounter++
	f.FeedbackID = "f-0001"
	m.feedbacks[f.FeedbackID] = f
	return f, nil
}

func (m *mockService) GetByID(_ context.Context, id string) (*domain.Feedback, error) {
	f, ok := m.feedbacks[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return f, nil
}

func (m *mockService) Update(_ context.Context, id string, updates *service.FeedbackUpdateInput) (*domain.Feedback, error) {
	f, ok := m.feedbacks[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	if updates.Rating != nil {
		f.Rating = *updates.Rating
	}
	if updates.Comment != nil {
		f.Comment = *updates.Comment
	}
	if updates.FeedbackType != nil {
		f.FeedbackType = domain.FeedbackType(*updates.FeedbackType)
	}
	return f, nil
}

func (m *mockService) List(_ context.Context, _ repository.FeedbackFilter) ([]*domain.Feedback, int, error) {
	var result []*domain.Feedback
	for _, f := range m.feedbacks {
		result = append(result, f)
	}
	return result, len(result), nil
}

func TestHandler_Create_Success(t *testing.T) {
	svc := newMockService()
	h := NewFeedbackHandler(svc)

	body := `{"user_id":"u-001","feedback_type":"bug","rating":3,"comment":"Test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/feedbacks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["success"] != true {
		t.Error("expected success to be true")
	}
}

func TestHandler_Create_InvalidJSON(t *testing.T) {
	svc := newMockService()
	h := NewFeedbackHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/feedbacks", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandler_Create_ValidationError(t *testing.T) {
	svc := newMockService()
	h := NewFeedbackHandler(svc)

	body := `{"user_id":"","feedback_type":"invalid","rating":0,"comment":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/feedbacks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandler_GetByID_NotFound(t *testing.T) {
	svc := newMockService()
	h := NewFeedbackHandler(svc)

	r := chi.NewRouter()
	r.Get("/api/v1/feedbacks/{id}", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/feedbacks/f-9999", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestHandler_Update_EmptyBody(t *testing.T) {
	svc := newMockService()
	h := NewFeedbackHandler(svc)

	r := chi.NewRouter()
	r.Patch("/api/v1/feedbacks/{id}", h.Update)

	body := `{}`
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/feedbacks/f-0001", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandler_List_Success(t *testing.T) {
	svc := newMockService()
	h := NewFeedbackHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/feedbacks", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["success"] != true {
		t.Error("expected success to be true")
	}
}

func TestHandler_List_InvalidFilter(t *testing.T) {
	svc := newMockService()
	h := NewFeedbackHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/feedbacks?min_rating=abc", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}
