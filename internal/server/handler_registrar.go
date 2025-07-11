package server

import (
	"aoroa/internal/handler"
	"aoroa/internal/service"
	serverPkg "aoroa/pkg/server"

	"github.com/gin-gonic/gin"
)

// IssueHandlerRegistrar는 이슈 관련 라우트를 등록하는 구조체입니다
type IssueHandlerRegistrar struct {
	userService  *service.UserService
	issueService *service.IssueService
}

// NewIssueHandlerRegistrar는 새로운 핸들러 등록자를 생성합니다
func NewIssueHandlerRegistrar() *IssueHandlerRegistrar {
	userService := service.NewUserService()
	issueService := service.NewIssueService(userService)

	return &IssueHandlerRegistrar{
		userService:  userService,
		issueService: issueService,
	}
}

// RegisterRoutes는 이슈 관련 라우트들을 등록합니다
func (r *IssueHandlerRegistrar) RegisterRoutes(framework serverPkg.WebFramework) error {
	// Gin 핸들러 래퍼 생성
	ginHandler := handler.NewGinIssueHandler(r.issueService)

	// API 라우트 등록 - 메서드를 함수로 변환
	framework.POST("/issue", gin.HandlerFunc(ginHandler.CreateIssue))
	framework.GET("/issues", gin.HandlerFunc(ginHandler.GetIssues))
	framework.GET("/issue/:id", gin.HandlerFunc(ginHandler.GetIssue))
	framework.PUT("/issue/:id", gin.HandlerFunc(ginHandler.UpdateIssue))

	return nil
}
