package server

import (
	"log"

	serverPkg "aoroa/pkg/server"
)

// Server represents the HTTP server
type Server struct {
	abstractServer serverPkg.ServerInterface
}

// New creates a new server instance
func New() *Server {
	// Gin 프레임워크 어댑터 생성
	ginFramework := serverPkg.NewGinFrameworkAdapter()

	// 핸들러 등록자 생성
	handlerRegistrar := NewIssueHandlerRegistrar()

	// 추상화된 서버 생성
	abstractServer := serverPkg.NewAbstractServer(ginFramework, handlerRegistrar)

	return &Server{
		abstractServer: abstractServer,
	}
}

// Initialize sets up the server with all dependencies
func (s *Server) Initialize() {
	if err := s.abstractServer.Initialize(); err != nil {
		log.Printf("Failed to initialize server: %v", err)
	}
}

// Run starts the server and blocks until shutdown
func (s *Server) Run() {
	s.Initialize()

	if err := s.abstractServer.Start(":8080"); err != nil {
		log.Printf("Server error: %v", err)
	}
}
