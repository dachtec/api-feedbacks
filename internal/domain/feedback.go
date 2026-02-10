package domain

import (
	"strings"
	"time"
)

// FeedbackType represents the type of feedback.
type FeedbackType string

const (
	FeedbackTypeBug        FeedbackType = "bug"
	FeedbackTypeSuggestion FeedbackType = "suggestion"
	FeedbackTypePraise     FeedbackType = "praise"
	FeedbackTypeQuestion   FeedbackType = "question"
)

// validFeedbackTypes contains all valid feedback types for validation.
var validFeedbackTypes = map[FeedbackType]bool{
	FeedbackTypeBug:        true,
	FeedbackTypeSuggestion: true,
	FeedbackTypePraise:     true,
	FeedbackTypeQuestion:   true,
}

// IsValidFeedbackType checks if the given string is a valid feedback type.
func IsValidFeedbackType(ft string) bool {
	return validFeedbackTypes[FeedbackType(ft)]
}

// Feedback represents a user feedback entity.
type Feedback struct {
	ID           string       `json:"id"`
	UserID       string       `json:"user_id"`
	FeedbackType FeedbackType `json:"feedback_type"`
	Rating       int          `json:"rating"`
	Comment      string       `json:"comment"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

const maxCommentLength = 2000

// Validate checks the feedback entity for domain rule violations.
func (f *Feedback) Validate() error {
	var errs []ValidationFieldError

	if strings.TrimSpace(f.UserID) == "" {
		errs = append(errs, ValidationFieldError{Field: "user_id", Message: "user_id is required"})
	}

	if !IsValidFeedbackType(string(f.FeedbackType)) {
		errs = append(errs, ValidationFieldError{
			Field:   "feedback_type",
			Message: "feedback_type must be one of: bug, suggestion, praise, question",
		})
	}

	if f.Rating < 1 || f.Rating > 5 {
		errs = append(errs, ValidationFieldError{Field: "rating", Message: "rating must be between 1 and 5"})
	}

	if strings.TrimSpace(f.Comment) == "" {
		errs = append(errs, ValidationFieldError{Field: "comment", Message: "comment is required"})
	} else if len(f.Comment) > maxCommentLength {
		errs = append(errs, ValidationFieldError{Field: "comment", Message: "comment must not exceed 2000 characters"})
	}

	if len(errs) > 0 {
		return &ValidationError{Errors: errs}
	}

	return nil
}
