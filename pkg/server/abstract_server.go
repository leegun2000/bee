package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// AbstractServer는 웹 프레임워크에 독립적인 서버 구현입니다
type AbstractServer struct {
	framework        WebFramework
	handlerRegistrar HandlerRegistrar
	srv              *http.Server
}

// NewAbstractServer는 새로운 추상 서버를 생성합니다
func NewAbstractServer(framework WebFramework, registrar HandlerRegistrar) ServerInterface {
	return &AbstractServer{
		framework:        framework,
		handlerRegistrar: registrar,
	}
}

// Initialize는 서버를 초기화합니다
func (s *AbstractServer) Initialize() error {
	if s.handlerRegistrar != nil {
		return s.handlerRegistrar.RegisterRoutes(s.framework)
	}
	return nil
}

// Start는 서버를 시작합니다
func (s *AbstractServer) Start(addr string) error {
	s.srv = &http.Server{
		Addr:    addr,
		Handler: s.framework.GetHTTPHandler(),
	}

	// 그레이스풀 셧다운을 위한 고루틴
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.srv.Shutdown(ctx); err != nil {
			log.Printf("Server forced to shutdown: %v", err)
		}
	}()

	log.Printf("Server starting on %s", addr)
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop은 서버를 중지합니다
func (s *AbstractServer) Stop() error {
	if s.srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.srv.Shutdown(ctx)
	}
	return nil
}
