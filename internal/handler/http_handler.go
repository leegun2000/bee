package handler

import (
	"net/http"

	"aoroa/internal/service"
	"aoroa/pkg/handlers"
	"aoroa/pkg/utils"
)

const (
	contentTypeHeader = "Content-Type"
	jsonContentType   = "application/json"
)

// HTTPHandler handles HTTP requests for issues using standard HTTP (for testing)
type HTTPHandler struct {
	handler handlers.IssueHandlerInterface
}

// NewHTTPHandler creates a new HTTPHandler for testing
func NewHTTPHandler(issueService *service.IssueService) *HTTPHandler {
	return &HTTPHandler{
		handler: NewIssueHandler(issueService),
	}
}

// CreateIssue handles POST requests for creating issues (standard HTTP)
func (h *HTTPHandler) CreateIssue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := utils.NewStandardHTTPAdapter(w, r)
	h.handler.CreateIssue(ctx)
}

// GetIssue handles GET requests for a single issue (standard HTTP)
func (h *HTTPHandler) GetIssue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from path and set it in context
	idParam := utils.ParseIDFromPath(r.URL.Path)
	if idParam == "" {
		idParam = "1" // Default for testing
	}

	ctx := utils.NewStandardHTTPAdapterWithParams(w, r, map[string]string{"id": idParam})
	h.handler.GetIssue(ctx)
}

// GetIssues handles GET requests for listing issues (standard HTTP)
func (h *HTTPHandler) GetIssues(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := utils.NewStandardHTTPAdapter(w, r)
	h.handler.GetIssues(ctx)
}
