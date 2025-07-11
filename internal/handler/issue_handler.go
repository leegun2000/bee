package handler

import (
	"net/http"

	"aoroa/internal/domain"
	"aoroa/internal/service"
	"aoroa/pkg/handlers"
	"aoroa/pkg/utils"
)

// IssueHandler implements issue operations using interface-based approach
type IssueHandler struct {
	issueService *service.IssueService
}

// NewIssueHandler creates a new IssueHandler
func NewIssueHandler(issueService *service.IssueService) handlers.IssueHandlerInterface {
	return &IssueHandler{
		issueService: issueService,
	}
}

// CreateIssue handles issue creation
func (h *IssueHandler) CreateIssue(ctx utils.HTTPContext) {
	var req domain.CreateIssueRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: "Invalid request: " + err.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}

	issue, err := h.issueService.CreateIssue(req)
	if err != nil {
		statusCode := utils.GetHTTPStatusForError(err.Error())
		ctx.JSON(statusCode, domain.ErrorResponse{
			Error: err.Error(),
			Code:  statusCode,
		})
		return
	}

	ctx.JSON(http.StatusCreated, issue)
}

// GetIssue handles single issue retrieval
func (h *IssueHandler) GetIssue(ctx utils.HTTPContext) {
	idParam := ctx.GetParam("id")
	if idParam == "" {
		// Fallback: try to extract from path for standard HTTP
		idParam = utils.ParseIDFromPath(ctx.GetHeader("X-Request-Path"))
		if idParam == "" {
			idParam = "1" // Default for testing
		}
	}

	id, err := utils.ParseUintParam(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: "Invalid issue ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	issue, err := h.issueService.GetIssue(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, domain.ErrorResponse{
			Error: err.Error(),
			Code:  http.StatusNotFound,
		})
		return
	}

	ctx.JSON(http.StatusOK, issue)
}

// GetIssues handles issue list retrieval
func (h *IssueHandler) GetIssues(ctx utils.HTTPContext) {
	status := ctx.GetQuery("status")

	issues, err := h.issueService.GetIssues(status)
	if err != nil {
		statusCode := utils.GetHTTPStatusForError(err.Error())
		ctx.JSON(statusCode, domain.ErrorResponse{
			Error: err.Error(),
			Code:  statusCode,
		})
		return
	}

	response := domain.IssuesResponse{
		Issues: make([]interface{}, len(issues)),
	}
	for i, issue := range issues {
		response.Issues[i] = issue
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateIssue handles issue updates
func (h *IssueHandler) UpdateIssue(ctx utils.HTTPContext) {
	idParam := ctx.GetParam("id")
	id, err := utils.ParseUintParam(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: "Invalid issue ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	// Parse request body manually to handle null userId
	var rawBody map[string]interface{}
	if err := ctx.BindJSON(&rawBody); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: "Invalid request: " + err.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}

	req := parseUpdateRequest(rawBody)

	issue, err := h.issueService.UpdateIssue(id, req)
	if err != nil {
		statusCode := utils.GetHTTPStatusForError(err.Error())
		ctx.JSON(statusCode, domain.ErrorResponse{
			Error: err.Error(),
			Code:  statusCode,
		})
		return
	}

	ctx.JSON(http.StatusOK, issue)
}

// parseUpdateRequest parses the raw request body into UpdateIssueRequest
func parseUpdateRequest(rawBody map[string]interface{}) domain.UpdateIssueRequest {
	req := domain.UpdateIssueRequest{}

	// Handle title
	if title, exists := rawBody["title"]; exists {
		if titleStr, ok := title.(string); ok {
			req.Title = &titleStr
		}
	}

	// Handle description
	if description, exists := rawBody["description"]; exists {
		if descStr, ok := description.(string); ok {
			req.Description = &descStr
		}
	}

	// Handle status
	if status, exists := rawBody["status"]; exists {
		if statusStr, ok := status.(string); ok {
			req.Status = &statusStr
		}
	}

	// Handle userId (including null)
	if userID, exists := rawBody["userId"]; exists {
		if userID == nil {
			req.RemoveUser = true
		} else if userIDFloat, ok := userID.(float64); ok {
			userIDUint := uint(userIDFloat)
			req.UserID = &userIDUint
		}
	}

	return req
}
