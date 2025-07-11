package handlers

// IssueHandlerInterface defines the interface for issue operations
type IssueHandlerInterface interface {
	CreateIssue(ctx HTTPContext)
	GetIssue(ctx HTTPContext)
	GetIssues(ctx HTTPContext)
	UpdateIssue(ctx HTTPContext)
}

// HTTPContext defines an interface for HTTP request/response operations
type HTTPContext interface {
	// Request parsing
	BindJSON(obj interface{}) error
	GetParam(key string) string
	GetQuery(key string) string
	GetHeader(key string) string

	// Response methods
	JSON(statusCode int, obj interface{})
	SetHeader(key, value string)
	Status(code int)
}
