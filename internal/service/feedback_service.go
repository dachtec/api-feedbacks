package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/dev/api-feedbacks/internal/domain"
	"github.com/dev/api-feedbacks/internal/repository"
)

// FeedbackService defines the business logic contract for feedbacks.
type FeedbackService interface {
	Create(ctx context.Context, f *domain.Feedback) (*domain.Feedback, error)
	GetByID(ctx context.Context, id string) (*domain.Feedback, error)
	Update(ctx context.Context, id string, updates *FeedbackUpdateInput) (*domain.Feedback, error)
	List(ctx context.Context, filter repository.FeedbackFilter) ([]*domain.Feedback, int, error)
}

// FeedbackUpdateInput represents the optional fields for a partial update.
type FeedbackUpdateInput struct {
	FeedbackType *string
	Rating       *int
	Comment      *string
}

type feedbackService struct {
	repo repository.FeedbackRepository
}

// NewFeedbackService creates a new feedback service with the given repository.
func NewFeedbackService(repo repository.FeedbackRepository) FeedbackService {
	return &feedbackService{repo: repo}
}

// Create validates and creates a new feedback.
func (s *feedbackService) Create(ctx context.Context, f *domain.Feedback) (*domain.Feedback, error) {
	f.ID = uuid.New().String()
	now := time.Now().UTC()
	f.CreatedAt = now
	f.UpdatedAt = now

	if err := f.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, f); err != nil {
		return nil, err
	}

	return f, nil
}

// GetByID retrieves a feedback by its unique ID.
func (s *feedbackService) GetByID(ctx context.Context, id string) (*domain.Feedback, error) {
	return s.repo.GetByID(ctx, id)
}

// Update applies partial updates to an existing feedback.
func (s *feedbackService) Update(ctx context.Context, id string, updates *FeedbackUpdateInput) (*domain.Feedback, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if updates.FeedbackType != nil {
		existing.FeedbackType = domain.FeedbackType(*updates.FeedbackType)
	}
	if updates.Rating != nil {
		existing.Rating = *updates.Rating
	}
	if updates.Comment != nil {
		existing.Comment = *updates.Comment
	}
	existing.UpdatedAt = time.Now().UTC()

	if err := existing.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// List returns feedbacks matching the given filter, plus a total count.
func (s *feedbackService) List(ctx context.Context, filter repository.FeedbackFilter) ([]*domain.Feedback, int, error) {
	if filter.Limit <= 0 {
		filter.Limit = repository.DefaultLimit
	}
	if filter.Limit > repository.MaxLimit {
		filter.Limit = repository.MaxLimit
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	return s.repo.List(ctx, filter)
}
