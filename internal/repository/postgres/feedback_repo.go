package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dev/api-feedbacks/internal/domain"
	"github.com/dev/api-feedbacks/internal/repository"
)

// FeedbackRepo implements repository.FeedbackRepository using PostgreSQL.
type FeedbackRepo struct {
	pool *pgxpool.Pool
}

// NewFeedbackRepo creates a new PostgreSQL-backed feedback repository.
func NewFeedbackRepo(pool *pgxpool.Pool) *FeedbackRepo {
	return &FeedbackRepo{pool: pool}
}

// Create inserts a new feedback record into the database.
func (r *FeedbackRepo) Create(ctx context.Context, f *domain.Feedback) error {
	query := `
		INSERT INTO feedbacks (id, user_id, feedback_type, rating, comment, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.pool.Exec(ctx, query,
		f.ID, f.UserID, f.FeedbackType, f.Rating, f.Comment, f.CreatedAt, f.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create feedback: %w", err)
	}

	return nil
}

// GetByID retrieves a single feedback by its ID.
func (r *FeedbackRepo) GetByID(ctx context.Context, id string) (*domain.Feedback, error) {
	query := `
		SELECT id, user_id, feedback_type, rating, comment, created_at, updated_at
		FROM feedbacks
		WHERE id = $1`

	f := &domain.Feedback{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&f.ID, &f.UserID, &f.FeedbackType, &f.Rating, &f.Comment, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get feedback: %w", err)
	}

	return f, nil
}

// Update modifies an existing feedback record.
func (r *FeedbackRepo) Update(ctx context.Context, f *domain.Feedback) error {
	query := `
		UPDATE feedbacks
		SET user_id = $2, feedback_type = $3, rating = $4, comment = $5, updated_at = $6
		WHERE id = $1`

	result, err := r.pool.Exec(ctx, query,
		f.ID, f.UserID, f.FeedbackType, f.Rating, f.Comment, f.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update feedback: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// List retrieves feedbacks matching the given filter with pagination.
func (r *FeedbackRepo) List(ctx context.Context, filter repository.FeedbackFilter) ([]*domain.Feedback, int, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if filter.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIdx))
		args = append(args, *filter.UserID)
		argIdx++
	}

	if filter.FeedbackType != nil {
		conditions = append(conditions, fmt.Sprintf("feedback_type = $%d", argIdx))
		args = append(args, *filter.FeedbackType)
		argIdx++
	}

	if filter.MinRating != nil {
		conditions = append(conditions, fmt.Sprintf("rating >= $%d", argIdx))
		args = append(args, *filter.MinRating)
		argIdx++
	}

	if filter.MaxRating != nil {
		conditions = append(conditions, fmt.Sprintf("rating <= $%d", argIdx))
		args = append(args, *filter.MaxRating)
		argIdx++
	}

	if filter.CreatedFrom != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIdx))
		args = append(args, *filter.CreatedFrom)
		argIdx++
	}

	if filter.CreatedTo != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIdx))
		args = append(args, *filter.CreatedTo)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total matching records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM feedbacks %s", whereClause)
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count feedbacks: %w", err)
	}

	// Fetch paginated results
	dataQuery := fmt.Sprintf(
		`SELECT id, user_id, feedback_type, rating, comment, created_at, updated_at
		 FROM feedbacks %s
		 ORDER BY created_at DESC
		 LIMIT $%d OFFSET $%d`,
		whereClause, argIdx, argIdx+1,
	)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list feedbacks: %w", err)
	}
	defer rows.Close()

	var feedbacks []*domain.Feedback
	for rows.Next() {
		f := &domain.Feedback{}
		if err := rows.Scan(&f.ID, &f.UserID, &f.FeedbackType, &f.Rating, &f.Comment, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan feedback: %w", err)
		}
		feedbacks = append(feedbacks, f)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating feedbacks: %w", err)
	}

	return feedbacks, total, nil
}
