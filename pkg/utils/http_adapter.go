package utils

import (
	"encoding/json"
	"net/http"
	"strings"
)

// StandardHTTPAdapter adapts standard http.ResponseWriter and http.Request to our HTTPContext interface
type StandardHTTPAdapter struct {
	writer  http.ResponseWriter
	request *http.Request
	params  map[string]string
}

// NewStandardHTTPAdapter creates a new adapter for standard HTTP
func NewStandardHTTPAdapter(w http.ResponseWriter, r *http.Request) HTTPContext {
	return &StandardHTTPAdapter{
		writer:  w,
		request: r,
		params:  make(map[string]string),
	}
}

// NewStandardHTTPAdapterWithParams creates a new adapter with URL parameters
func NewStandardHTTPAdapterWithParams(w http.ResponseWriter, r *http.Request, params map[string]string) HTTPContext {
	return &StandardHTTPAdapter{
		writer:  w,
		request: r,
		params:  params,
	}
}

// BindJSON binds the request body to the given struct
func (s *StandardHTTPAdapter) BindJSON(obj interface{}) error {
	return json.NewDecoder(s.request.Body).Decode(obj)
}

// GetParam gets a URL parameter by key
func (s *StandardHTTPAdapter) GetParam(key string) string {
	if value, exists := s.params[key]; exists {
		return value
	}
	return ""
}

// GetQuery gets a query parameter by key
func (s *StandardHTTPAdapter) GetQuery(key string) string {
	return s.request.URL.Query().Get(key)
}

// GetHeader gets a header value by key
func (s *StandardHTTPAdapter) GetHeader(key string) string {
	return s.request.Header.Get(key)
}

// JSON sends a JSON response
func (s *StandardHTTPAdapter) JSON(statusCode int, obj interface{}) {
	s.writer.Header().Set("Content-Type", "application/json")
	s.writer.WriteHeader(statusCode)
	json.NewEncoder(s.writer).Encode(obj)
}

// SetHeader sets a response header
func (s *StandardHTTPAdapter) SetHeader(key, value string) {
	s.writer.Header().Set(key, value)
}

// Status sets the response status code
func (s *StandardHTTPAdapter) Status(code int) {
	s.writer.WriteHeader(code)
}

// SetParams sets URL parameters (useful for testing or manual routing)
func (s *StandardHTTPAdapter) SetParams(params map[string]string) {
	s.params = params
}

// ParseIDFromPath extracts ID from URL path (simple implementation)
func ParseIDFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 {
		return parts[1] // assumes /issues/123 format
	}
	return ""
}
