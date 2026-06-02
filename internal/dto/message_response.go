package dto

// MessageResponse represents a simple message payload used for error or status replies.
type MessageResponse struct {
	Message string `json:"message"`
}
