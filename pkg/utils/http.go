package utils

import (
	"encoding/json"
	"net/http"
	"strconv"

	"aoroa/pkg/handlers"
)

// HTTPContext는 handlers 패키지의 HTTPContext를 재사용합니다
type HTTPContext = handlers.HTTPContext

// HTTPRequest represents an HTTP request abstraction
type HTTPRequest interface {
	GetParam(key string) string
	GetQuery(key string) string
	GetHeader(key string) string
	BindJSON(obj interface{}) error
}

// HTTPResponse represents an HTTP response abstraction
type HTTPResponse interface {
	JSON(statusCode int, obj interface{})
	SetHeader(key, value string)
	Status(code int)
}

// ParseUintParam parses a string parameter to uint
func ParseUintParam(param string) (uint, error) {
	id, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// WriteJSONResponse writes a JSON response with the given status code
func WriteJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

// WriteJSONError writes a JSON error response
func WriteJSONError(w http.ResponseWriter, message string, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errorResponse := map[string]interface{}{
		"error": message,
		"code":  statusCode,
	}
	return json.NewEncoder(w).Encode(errorResponse)
}

// DecodeJSONRequest decodes a JSON request body into the given interface
func DecodeJSONRequest(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// GetHTTPStatusForError returns appropriate HTTP status code for common errors
func GetHTTPStatusForError(errMsg string) int {
	switch {
	case errMsg == "user not found" || errMsg == "issue not found":
		return http.StatusNotFound
	case errMsg == "cannot update completed or cancelled issue":
		return http.StatusConflict
	case errMsg == "invalid status":
		return http.StatusBadRequest
	default:
		return http.StatusBadRequest
	}
}
