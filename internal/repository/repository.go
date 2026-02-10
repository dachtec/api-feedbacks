package repository

import (
	"context"
	"time"

	"github.com/dev/api-feedbacks/internal/domain"
)

// FeedbackFilter contains optional filters for listing feedbacks.
type FeedbackFilter struct {
	UserID       *string
	FeedbackType *string
	MinRating    *int
	MaxRating    *int
	CreatedFrom  *time.Time
	CreatedTo    *time.Time
	Limit        int
	Offset       int
}

// DefaultLimit is the default number of results per page.
const DefaultLimit = 20

// MaxLimit is the maximum number of results per page.
const MaxLimit = 100

// FeedbackRepository defines the contract for feedback persistence.
type FeedbackRepository interface {
	Create(ctx context.Context, f *domain.Feedback) error
	GetByID(ctx context.Context, id string) (*domain.Feedback, error)
	Update(ctx context.Context, f *domain.Feedback) error
	List(ctx context.Context, filter FeedbackFilter) ([]*domain.Feedback, int, error)
}
