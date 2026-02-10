package handler

// CreateFeedbackRequest represents the request body for creating a feedback.
type CreateFeedbackRequest struct {
	UserID       string `json:"user_id"`
	FeedbackType string `json:"feedback_type"`
	Rating       int    `json:"rating"`
	Comment      string `json:"comment"`
}

// UpdateFeedbackRequest represents the request body for updating a feedback.
// All fields are optional (pointers) to support partial updates.
type UpdateFeedbackRequest struct {
	FeedbackType *string `json:"feedback_type,omitempty"`
	Rating       *int    `json:"rating,omitempty"`
	Comment      *string `json:"comment,omitempty"`
}

// FeedbackResponse represents the API response for a single feedback.
type FeedbackResponse struct {
	FeedbackID   string `json:"feedback_id"`
	UserID       string `json:"user_id"`
	FeedbackType string `json:"feedback_type"`
	Rating       int    `json:"rating"`
	Comment      string `json:"comment"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// ListMeta contains pagination metadata for list responses.
type ListMeta struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
