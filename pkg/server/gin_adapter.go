package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GinFrameworkAdapter는 Gin을 WebFramework 인터페이스에 맞게 어댑터하는 구조체입니다
type GinFrameworkAdapter struct {
	engine *gin.Engine
}

// NewGinFrameworkAdapter는 새로운 Gin 어댑터를 생성합니다
func NewGinFrameworkAdapter() WebFramework {
	engine := gin.Default()

	// Health check endpoint
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return &GinFrameworkAdapter{
		engine: engine,
	}
}

// GET은 GET 라우트를 등록합니다
func (g *GinFrameworkAdapter) GET(path string, handlerFunc interface{}) {
	if handler, ok := handlerFunc.(gin.HandlerFunc); ok {
		g.engine.GET(path, handler)
	}
}

// POST는 POST 라우트를 등록합니다
func (g *GinFrameworkAdapter) POST(path string, handlerFunc interface{}) {
	if handler, ok := handlerFunc.(gin.HandlerFunc); ok {
		g.engine.POST(path, handler)
	}
}

// PUT은 PUT 라우트를 등록합니다
func (g *GinFrameworkAdapter) PUT(path string, handlerFunc interface{}) {
	if handler, ok := handlerFunc.(gin.HandlerFunc); ok {
		g.engine.PUT(path, handler)
	}
}

// DELETE는 DELETE 라우트를 등록합니다
func (g *GinFrameworkAdapter) DELETE(path string, handlerFunc interface{}) {
	if handler, ok := handlerFunc.(gin.HandlerFunc); ok {
		g.engine.DELETE(path, handler)
	}
}

// Run은 서버를 시작합니다
func (g *GinFrameworkAdapter) Run(addr string) error {
	return g.engine.Run(addr)
}

// GetHTTPHandler는 HTTP 핸들러를 반환합니다
func (g *GinFrameworkAdapter) GetHTTPHandler() http.Handler {
	return g.engine
}
