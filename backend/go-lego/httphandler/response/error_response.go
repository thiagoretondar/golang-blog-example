package response

// ErrorResponse contains information about request that went wrong somehow
type ErrorResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message,omitempty"`
}
