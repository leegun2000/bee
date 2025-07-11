package service

import (
	"errors"
	"sync"
	"time"

	"aoroa/internal/domain"
	"aoroa/internal/models"
)

// IssueService handles issue-related operations
type IssueService struct {
	issues      map[uint]*models.Issue
	userService *UserService
	nextID      uint
	mu          sync.RWMutex
}

// NewIssueService creates a new IssueService
func NewIssueService(userService *UserService) *IssueService {
	return &IssueService{
		issues:      make(map[uint]*models.Issue),
		userService: userService,
		nextID:      1,
	}
}

// CreateIssue creates a new issue
func (s *IssueService) CreateIssue(req domain.CreateIssueRequest) (*models.Issue, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate user if provided
	var user *models.User
	if req.UserID != nil {
		if u, exists := s.userService.GetUser(*req.UserID); exists {
			user = u
		} else {
			return nil, errors.New("user not found")
		}
	}

	// Determine initial status
	status := domain.StatusPending
	if user != nil {
		status = domain.StatusInProgress
	}

	now := time.Now()
	issue := &models.Issue{
		ID:          s.nextID,
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		User:        user,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	s.issues[s.nextID] = issue
	s.nextID++

	return issue, nil
}

// GetIssue retrieves an issue by ID
func (s *IssueService) GetIssue(id uint) (*models.Issue, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	issue, exists := s.issues[id]
	if !exists {
		return nil, errors.New("issue not found")
	}

	return issue, nil
}

// GetIssues retrieves all issues, optionally filtered by status
func (s *IssueService) GetIssues(status string) ([]models.Issue, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Validate status if provided
	if status != "" && !domain.IsValidStatus(status) {
		return nil, errors.New("invalid status")
	}

	var result []models.Issue
	for _, issue := range s.issues {
		if status == "" || issue.Status == status {
			result = append(result, *issue)
		}
	}

	return result, nil
}

// UpdateIssue updates an existing issue
func (s *IssueService) UpdateIssue(id uint, req domain.UpdateIssueRequest) (*models.Issue, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	issue, exists := s.issues[id]
	if !exists {
		return nil, errors.New("issue not found")
	}

	// Check if issue is in final state
	if issue.Status == domain.StatusCompleted || issue.Status == domain.StatusCancelled {
		return nil, errors.New("cannot update completed or cancelled issue")
	}

	// Validate status if provided
	if req.Status != nil && !domain.IsValidStatus(*req.Status) {
		return nil, errors.New("invalid status")
	}

	// Handle user assignment/removal
	newUser, userChanged, err := s.handleUserChange(issue, req)
	if err != nil {
		return nil, err
	}

	// Determine new status based on business rules
	newStatus := s.determineNewStatus(issue, req, newUser, userChanged)

	// Validate that non-PENDING/CANCELLED statuses have assignees
	if newUser == nil && newStatus != domain.StatusPending && newStatus != domain.StatusCancelled {
		return nil, errors.New("cannot set status to " + newStatus + " without assignee")
	}

	// Update issue fields
	s.updateIssueFields(issue, req, newStatus, newUser)

	return issue, nil
}

// handleUserChange handles user assignment/removal logic
func (s *IssueService) handleUserChange(issue *models.Issue, req domain.UpdateIssueRequest) (*models.User, bool, error) {
	var newUser *models.User
	var userChanged bool

	if req.UserID != nil {
		// User assignment
		if u, exists := s.userService.GetUser(*req.UserID); exists {
			newUser = u
			userChanged = true
		} else {
			return nil, false, errors.New("user not found")
		}
	} else if req.RemoveUser {
		// User removal
		newUser = nil
		userChanged = true
	} else {
		// Keep existing user
		newUser = issue.User
	}

	return newUser, userChanged, nil
}

// determineNewStatus determines the new status based on business rules
func (s *IssueService) determineNewStatus(issue *models.Issue, req domain.UpdateIssueRequest, newUser *models.User, userChanged bool) string {
	newStatus := issue.Status
	if req.Status != nil {
		newStatus = *req.Status
	}

	// Apply business rules
	if userChanged {
		if newUser == nil {
			// User removed -> status becomes PENDING
			newStatus = domain.StatusPending
		} else if issue.Status == domain.StatusPending && req.Status == nil {
			// User assigned to PENDING issue without explicit status -> IN_PROGRESS
			newStatus = domain.StatusInProgress
		}
	}

	return newStatus
}

// updateIssueFields updates the issue with new values
func (s *IssueService) updateIssueFields(issue *models.Issue, req domain.UpdateIssueRequest, newStatus string, newUser *models.User) {
	if req.Title != nil {
		issue.Title = *req.Title
	}
	if req.Description != nil {
		issue.Description = *req.Description
	}
	issue.Status = newStatus
	issue.User = newUser
	issue.UpdatedAt = time.Now()
}
