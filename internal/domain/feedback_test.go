package domain

import (
	"testing"
)

func TestFeedback_Validate_Valid(t *testing.T) {
	f := &Feedback{
		UserID:       "u-001",
		FeedbackType: FeedbackTypeBug,
		Rating:       3,
		Comment:      "This is a valid comment",
	}

	if err := f.Validate(); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestFeedback_Validate_EmptyUserID(t *testing.T) {
	f := &Feedback{
		UserID:       "",
		FeedbackType: FeedbackTypeBug,
		Rating:       3,
		Comment:      "Valid comment",
	}

	err := f.Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	ve, ok := IsValidationError(err)
	if !ok {
		t.Fatal("expected ValidationError type")
	}

	found := false
	for _, fe := range ve.Errors {
		if fe.Field == "user_id" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected validation error for field user_id")
	}
}

func TestFeedback_Validate_InvalidUserIDFormat(t *testing.T) {
	tests := []struct {
		name   string
		userID string
	}{
		{"missing prefix", "001"},
		{"wrong prefix", "usr-001"},
		{"too many digits", "u-0001"},
		{"too few digits", "u-01"},
		{"letters instead of digits", "u-abc"},
		{"uppercase prefix", "U-001"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Feedback{
				UserID:       tt.userID,
				FeedbackType: FeedbackTypeBug,
				Rating:       3,
				Comment:      "Valid comment",
			}

			err := f.Validate()
			if err == nil {
				t.Errorf("expected validation error for user_id %q, got nil", tt.userID)
			}
		})
	}
}

func TestFeedback_Validate_ValidUserIDFormats(t *testing.T) {
	tests := []string{"u-001", "u-015", "u-999", "u-100"}

	for _, uid := range tests {
		t.Run(uid, func(t *testing.T) {
			f := &Feedback{
				UserID:       uid,
				FeedbackType: FeedbackTypeBug,
				Rating:       3,
				Comment:      "Valid comment",
			}

			if err := f.Validate(); err != nil {
				t.Errorf("expected no error for user_id %q, got: %v", uid, err)
			}
		})
	}
}

func TestFeedback_Validate_InvalidFeedbackType(t *testing.T) {
	f := &Feedback{
		UserID:       "u-001",
		FeedbackType: "invalid",
		Rating:       3,
		Comment:      "Valid comment",
	}

	err := f.Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestFeedback_Validate_RatingOutOfRange(t *testing.T) {
	tests := []struct {
		name   string
		rating int
	}{
		{"too low", 0},
		{"too high", 6},
		{"negative", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Feedback{
				UserID:       "u-001",
				FeedbackType: FeedbackTypeBug,
				Rating:       tt.rating,
				Comment:      "Valid comment",
			}

			if err := f.Validate(); err == nil {
				t.Errorf("expected validation error for rating %d, got nil", tt.rating)
			}
		})
	}
}

func TestFeedback_Validate_EmptyComment(t *testing.T) {
	f := &Feedback{
		UserID:       "u-001",
		FeedbackType: FeedbackTypeBug,
		Rating:       3,
		Comment:      "",
	}

	err := f.Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestFeedback_Validate_CommentTooLong(t *testing.T) {
	longComment := make([]byte, 2001)
	for i := range longComment {
		longComment[i] = 'a'
	}

	f := &Feedback{
		UserID:       "u-001",
		FeedbackType: FeedbackTypeBug,
		Rating:       3,
		Comment:      string(longComment),
	}

	err := f.Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestFeedback_Validate_MultipleErrors(t *testing.T) {
	f := &Feedback{
		UserID:       "",
		FeedbackType: "invalid",
		Rating:       0,
		Comment:      "",
	}

	err := f.Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	ve, ok := IsValidationError(err)
	if !ok {
		t.Fatal("expected ValidationError type")
	}

	if len(ve.Errors) < 4 {
		t.Errorf("expected at least 4 validation errors, got %d", len(ve.Errors))
	}
}

func TestIsValidFeedbackType(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"bug", true},
		{"sugerencia", true},
		{"elogio", true},
		{"duda", true},
		{"queja", true},
		{"suggestion", false},
		{"praise", false},
		{"question", false},
		{"invalid", false},
		{"", false},
		{"BUG", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsValidFeedbackType(tt.input); got != tt.valid {
				t.Errorf("IsValidFeedbackType(%q) = %v, want %v", tt.input, got, tt.valid)
			}
		})
	}
}
