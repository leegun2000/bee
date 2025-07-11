package domain

// IssueStatus constants
const (
	StatusPending    = "PENDING"
	StatusInProgress = "IN_PROGRESS"
	StatusCompleted  = "COMPLETED"
	StatusCancelled  = "CANCELLED"
)

// IsValidStatus checks if the given status is valid
func IsValidStatus(status string) bool {
	switch status {
	case StatusPending, StatusInProgress, StatusCompleted, StatusCancelled:
		return true
	default:
		return false
	}
}

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

// CreateIssueRequest represents the request payload for creating an issue
type CreateIssueRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	UserID      *uint  `json:"userId,omitempty"`
}

// UpdateIssueRequest represents the request payload for updating an issue
type UpdateIssueRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
	UserID      *uint   `json:"userId,omitempty"`
	RemoveUser  bool    `json:"-"` // Internal flag for removing user
}

// IssueResponse represents the response for a single issue
type IssueResponse struct {
	ID          uint    `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	UserID      *uint   `json:"userId,omitempty"`
	UserName    *string `json:"userName,omitempty"`
}

// GetIssuesResponse represents the response for listing issues
type GetIssuesResponse struct {
	Issues []IssueResponse `json:"issues"`
}

// IssuesResponse represents the response for listing issues
type IssuesResponse struct {
	Issues []interface{} `json:"issues"` // Will be []models.Issue
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}
