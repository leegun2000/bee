package handler

import (
	"aoroa/internal/service"
	"aoroa/pkg/handlers"
	"aoroa/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GinIssueHandler wraps IssueHandler for Gin compatibility
type GinIssueHandler struct {
	handler handlers.IssueHandlerInterface
}

// NewGinIssueHandler creates a new Gin-compatible handler
func NewGinIssueHandler(issueService *service.IssueService) *GinIssueHandler {
	return &GinIssueHandler{
		handler: NewIssueHandler(issueService),
	}
}

// CreateIssue handles POST /issues for Gin
func (g *GinIssueHandler) CreateIssue(c *gin.Context) {
	ctx := utils.NewGinContextAdapter(c)
	g.handler.CreateIssue(ctx)
}

// GetIssue handles GET /issues/:id for Gin
func (g *GinIssueHandler) GetIssue(c *gin.Context) {
	ctx := utils.NewGinContextAdapter(c)
	g.handler.GetIssue(ctx)
}

// GetIssues handles GET /issues for Gin
func (g *GinIssueHandler) GetIssues(c *gin.Context) {
	ctx := utils.NewGinContextAdapter(c)
	g.handler.GetIssues(ctx)
}

// UpdateIssue handles PUT /issues/:id for Gin
func (g *GinIssueHandler) UpdateIssue(c *gin.Context) {
	ctx := utils.NewGinContextAdapter(c)
	g.handler.UpdateIssue(ctx)
}
