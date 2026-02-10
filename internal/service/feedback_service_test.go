package service

import (
	"context"
	"testing"

	"github.com/dev/api-feedbacks/internal/domain"
	"github.com/dev/api-feedbacks/internal/repository"
)

// mockRepo implements repository.FeedbackRepository for testing.
type mockRepo struct {
	feedbacks map[string]*domain.Feedback
	createErr error
	updateErr error
}

func newMockRepo() *mockRepo {
	return &mockRepo{feedbacks: make(map[string]*domain.Feedback)}
}

func (m *mockRepo) Create(_ context.Context, f *domain.Feedback) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.feedbacks[f.ID] = f
	return nil
}

func (m *mockRepo) GetByID(_ context.Context, id string) (*domain.Feedback, error) {
	f, ok := m.feedbacks[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return f, nil
}

func (m *mockRepo) Update(_ context.Context, f *domain.Feedback) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, ok := m.feedbacks[f.ID]; !ok {
		return domain.ErrNotFound
	}
	m.feedbacks[f.ID] = f
	return nil
}

func (m *mockRepo) List(_ context.Context, _ repository.FeedbackFilter) ([]*domain.Feedback, int, error) {
	var result []*domain.Feedback
	for _, f := range m.feedbacks {
		result = append(result, f)
	}
	return result, len(result), nil
}

func TestFeedbackService_Create_Success(t *testing.T) {
	repo := newMockRepo()
	svc := NewFeedbackService(repo)

	feedback := &domain.Feedback{
		UserID:       "usr-001",
		FeedbackType: domain.FeedbackTypeBug,
		Rating:       3,
		Comment:      "Something is broken",
	}

	created, err := svc.Create(context.Background(), feedback)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if created.ID == "" {
		t.Error("expected ID to be generated")
	}

	if created.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}

	if created.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}

	if len(repo.feedbacks) != 1 {
		t.Errorf("expected 1 feedback in repo, got %d", len(repo.feedbacks))
	}
}

func TestFeedbackService_Create_ValidationError(t *testing.T) {
	repo := newMockRepo()
	svc := NewFeedbackService(repo)

	feedback := &domain.Feedback{
		UserID:       "",
		FeedbackType: "invalid",
		Rating:       0,
		Comment:      "",
	}

	_, err := svc.Create(context.Background(), feedback)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	if _, ok := domain.IsValidationError(err); !ok {
		t.Error("expected error to be ValidationError")
	}

	if len(repo.feedbacks) != 0 {
		t.Error("expected no feedback to be created on validation error")
	}
}

func TestFeedbackService_GetByID_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := NewFeedbackService(repo)

	_, err := svc.GetByID(context.Background(), "non-existent")
	if err != domain.ErrNotFound {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestFeedbackService_Update_Success(t *testing.T) {
	repo := newMockRepo()
	svc := NewFeedbackService(repo)

	// Create a feedback first
	feedback := &domain.Feedback{
		UserID:       "usr-001",
		FeedbackType: domain.FeedbackTypeBug,
		Rating:       3,
		Comment:      "Original comment",
	}
	created, err := svc.Create(context.Background(), feedback)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// Update it
	newRating := 5
	newComment := "Updated comment"
	updated, err := svc.Update(context.Background(), created.ID, &FeedbackUpdateInput{
		Rating:  &newRating,
		Comment: &newComment,
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if updated.Rating != 5 {
		t.Errorf("expected rating 5, got %d", updated.Rating)
	}

	if updated.Comment != "Updated comment" {
		t.Errorf("expected updated comment, got %q", updated.Comment)
	}

	if !updated.UpdatedAt.After(created.CreatedAt) {
		t.Error("expected UpdatedAt to be after CreatedAt")
	}
}

func TestFeedbackService_Update_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := NewFeedbackService(repo)

	newRating := 5
	_, err := svc.Update(context.Background(), "non-existent", &FeedbackUpdateInput{
		Rating: &newRating,
	})

	if err != domain.ErrNotFound {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestFeedbackService_List_DefaultPagination(t *testing.T) {
	repo := newMockRepo()
	svc := NewFeedbackService(repo)

	// Create some feedbacks
	for i := 0; i < 3; i++ {
		f := &domain.Feedback{
			UserID:       "usr-001",
			FeedbackType: domain.FeedbackTypeBug,
			Rating:       3,
			Comment:      "Test comment",
		}
		if _, err := svc.Create(context.Background(), f); err != nil {
			t.Fatalf("setup failed: %v", err)
		}
	}

	filter := repository.FeedbackFilter{}
	feedbacks, total, err := svc.List(context.Background(), filter)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}

	if len(feedbacks) != 3 {
		t.Errorf("expected 3 feedbacks, got %d", len(feedbacks))
	}
}
