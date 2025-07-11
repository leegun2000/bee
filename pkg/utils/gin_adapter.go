package utils

import (
	"github.com/gin-gonic/gin"
)

// GinContextAdapter adapts gin.Context to our HTTPContext interface
type GinContextAdapter struct {
	ctx *gin.Context
}

// NewGinContextAdapter creates a new adapter for gin.Context
func NewGinContextAdapter(ctx *gin.Context) HTTPContext {
	return &GinContextAdapter{ctx: ctx}
}

// BindJSON binds the request body to the given struct
func (g *GinContextAdapter) BindJSON(obj interface{}) error {
	return g.ctx.ShouldBindJSON(obj)
}

// GetParam gets a URL parameter by key
func (g *GinContextAdapter) GetParam(key string) string {
	return g.ctx.Param(key)
}

// GetQuery gets a query parameter by key
func (g *GinContextAdapter) GetQuery(key string) string {
	return g.ctx.Query(key)
}

// GetHeader gets a header value by key
func (g *GinContextAdapter) GetHeader(key string) string {
	return g.ctx.GetHeader(key)
}

// JSON sends a JSON response
func (g *GinContextAdapter) JSON(statusCode int, obj interface{}) {
	g.ctx.JSON(statusCode, obj)
}

// SetHeader sets a response header
func (g *GinContextAdapter) SetHeader(key, value string) {
	g.ctx.Header(key, value)
}

// Status sets the response status code
func (g *GinContextAdapter) Status(code int) {
	g.ctx.Status(code)
}
