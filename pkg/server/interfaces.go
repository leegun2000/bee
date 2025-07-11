package server

import (
	"net/http"
)

// WebFramework 인터페이스는 웹 프레임워크의 공통 기능을 추상화합니다
type WebFramework interface {
	// Route registration
	GET(path string, handlerFunc interface{})
	POST(path string, handlerFunc interface{})
	PUT(path string, handlerFunc interface{})
	DELETE(path string, handlerFunc interface{})

	// Server control
	Run(addr string) error
	GetHTTPHandler() http.Handler
}

// ServerInterface는 서버의 공통 기능을 추상화합니다
type ServerInterface interface {
	Initialize() error
	Start(addr string) error
	Stop() error
}

// HandlerRegistrar는 라우트 등록을 담당하는 인터페이스입니다
type HandlerRegistrar interface {
	RegisterRoutes(framework WebFramework) error
}
