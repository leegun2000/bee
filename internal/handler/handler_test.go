package handler

import (
"bytes"
"encoding/json"
"net/http"
"net/http/httptest"
"testing"

"aoroa/internal/domain"
"aoroa/internal/service"
"aoroa/pkg/utils"
)

// TestNewIssueHandler tests the interface-based handler creation
func TestNewIssueHandler(t *testing.T) {
	userService := service.NewUserService()
	issueService := service.NewIssueService(userService)

	handler := NewIssueHandler(issueService)

	if handler == nil {
		t.Fatal("NewIssueHandler() returned nil")
	}
}

// TestCreateIssueWithStandardHTTP tests CreateIssue using standard HTTP
func TestCreateIssueWithStandardHTTP(t *testing.T) {
	userService := service.NewUserService()
	issueService := service.NewIssueService(userService)
	handler := NewIssueHandler(issueService)

	request := domain.CreateIssueRequest{
		Title:       "Test Issue",
		Description: "Test Description",
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/issues", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Standard HTTP adapter를 사용하여 인터페이스 기반 핸들러 테스트
	ctx := utils.NewStandardHTTPAdapter(rr, req)
	handler.CreateIssue(ctx)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, status)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["title"] != request.Title {
		t.Errorf("Expected title %s, got %s", request.Title, response["title"])
	}
}
