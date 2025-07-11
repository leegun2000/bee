package service

import (
	"testing"

	"aoroa/internal/domain"
	"aoroa/internal/models"
)

// Test constants to avoid duplication
const (
	testTitle         = "Test Issue"
	testDescription   = "Test Description"
	errorUnexpected   = "Unexpected error: %v"
	errorExpectedNone = "Expected error but got none"
	userNotFound      = "user not found"
	testUserName      = "Test User"
	testUserEmail     = "test@example.com"
)

func TestNewIssueService(t *testing.T) {
	userService := NewUserService()
	issueService := NewIssueService(userService)

	if issueService == nil {
		t.Fatal("NewIssueService() returned nil")
	}
}

func TestIssueServiceCreateIssueWithoutUser(t *testing.T) {
	userService := NewUserService()
	issueService := NewIssueService(userService)

	request := domain.CreateIssueRequest{
		Title:       testTitle,
		Description: testDescription,
	}

	issue, err := issueService.CreateIssue(request)

	if err != nil {
		t.Fatalf(errorUnexpected, err)
	}

	assertIssueCreated(t, issue, request, domain.StatusPending)
}

func TestIssueServiceCreateIssueWithValidUser(t *testing.T) {
	userService := NewUserService()
	issueService := NewIssueService(userService)

	// First create a user
	userReq := domain.CreateUserRequest{
		Name:  testUserName,
		Email: testUserEmail,
	}
	_, err := userService.CreateUser(userReq)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	request := domain.CreateIssueRequest{
		Title:       testTitle + " with User",
		Description: testDescription,
		UserID:      uintPtr(1),
	}

	issue, err := issueService.CreateIssue(request)

	if err != nil {
		t.Fatalf(errorUnexpected, err)
	}

	assertIssueCreated(t, issue, request, domain.StatusInProgress)
}

func TestIssueServiceCreateIssueWithInvalidUser(t *testing.T) {
	userService := NewUserService()
	issueService := NewIssueService(userService)

	request := domain.CreateIssueRequest{
		Title:       testTitle,
		Description: testDescription,
		UserID:      uintPtr(999),
	}

	_, err := issueService.CreateIssue(request)

	if err == nil {
		t.Fatal(errorExpectedNone)
	}

	if err.Error() != userNotFound {
		t.Errorf("Expected error message '%s', got '%s'", userNotFound, err.Error())
	}
}

func TestIssueServiceGetExistingIssue(t *testing.T) {
	userService := NewUserService()
	issueService := NewIssueService(userService)

	// Create a test issue first
	req := domain.CreateIssueRequest{
		Title:       testTitle,
		Description: testDescription,
	}
	createdIssue, err := issueService.CreateIssue(req)
	if err != nil {
		t.Fatalf("Failed to create test issue: %v", err)
	}

	// Get the issue
	issue, err := issueService.GetIssue(createdIssue.ID)

	if err != nil {
		t.Fatalf(errorUnexpected, err)
	}

	if issue == nil {
		t.Fatal("GetIssue() returned nil issue")
	}

	if issue.ID != createdIssue.ID {
		t.Errorf("Expected issue ID %d, got %d", createdIssue.ID, issue.ID)
	}
}

func TestIssueServiceGetNonExistingIssue(t *testing.T) {
	userService := NewUserService()
	issueService := NewIssueService(userService)

	_, err := issueService.GetIssue(999)

	if err == nil {
		t.Fatal(errorExpectedNone)
	}
}

func TestIssueServiceGetAllIssues(t *testing.T) {
	userService := NewUserService()
	issueService := NewIssueService(userService)

	// Create test issues
	req1 := domain.CreateIssueRequest{Title: "Pending Issue", Description: testDescription}
	req2 := domain.CreateIssueRequest{Title: "In Progress Issue", Description: testDescription, UserID: uintPtr(1)}

	_, err := issueService.CreateIssue(req1)
	if err != nil {
		t.Fatalf("Failed to create test issue 1: %v", err)
	}

	_, err = issueService.CreateIssue(req2)
	if err != nil {
		t.Fatalf("Failed to create test issue 2: %v", err)
	}

	issues, err := issueService.GetIssues("")

	if err != nil {
		t.Fatalf(errorUnexpected, err)
	}

	expectedLen := 2
	if len(issues) != expectedLen {
		t.Errorf("Expected %d issues, got %d", expectedLen, len(issues))
	}
}

func TestIssueServiceGetPendingIssues(t *testing.T) {
	userService := NewUserService()
	issueService := NewIssueService(userService)

	// Create a pending issue
	req := domain.CreateIssueRequest{Title: "Pending Issue", Description: testDescription}
	_, err := issueService.CreateIssue(req)
	if err != nil {
		t.Fatalf("Failed to create test issue: %v", err)
	}

	issues, err := issueService.GetIssues(domain.StatusPending)

	if err != nil {
		t.Fatalf(errorUnexpected, err)
	}

	expectedLen := 1
	if len(issues) != expectedLen {
		t.Errorf("Expected %d issues, got %d", expectedLen, len(issues))
	}
}

func TestIssueServiceGetIssuesWithInvalidStatus(t *testing.T) {
	userService := NewUserService()
	issueService := NewIssueService(userService)

	_, err := issueService.GetIssues("INVALID")

	if err == nil {
		t.Fatal(errorExpectedNone)
	}
}

// Helper function to create a pointer to uint
func uintPtr(u uint) *uint {
	return &u
}

// Helper function to assert issue creation
func assertIssueCreated(t *testing.T, issue *models.Issue, request domain.CreateIssueRequest, expectedStatus string) {
	t.Helper()

	if issue == nil {
		t.Fatal("CreateIssue() returned nil issue")
	}

	if issue.Title != request.Title {
		t.Errorf("Expected title '%s', got '%s'", request.Title, issue.Title)
	}

	if issue.Description != request.Description {
		t.Errorf("Expected description '%s', got '%s'", request.Description, issue.Description)
	}

	if issue.Status != expectedStatus {
		t.Errorf("Expected status '%s', got '%s'", expectedStatus, issue.Status)
	}

	if issue.ID == 0 {
		t.Error("Expected issue ID to be set")
	}

	if issue.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if issue.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}
