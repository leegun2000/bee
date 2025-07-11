package service

import (
	"sync"

	"aoroa/internal/domain"
	"aoroa/internal/models"
)

// UserService handles user-related operations
type UserService struct {
	users map[uint]*models.User
	mu    sync.RWMutex
}

// NewUserService creates a new UserService with predefined users
func NewUserService() *UserService {
	service := &UserService{
		users: make(map[uint]*models.User),
	}

	// Initialize with required users
	predefinedUsers := []*models.User{
		{ID: 1, Name: "김개발", Email: "kim@example.com"},
		{ID: 2, Name: "이디자인", Email: "lee@example.com"},
		{ID: 3, Name: "박기획", Email: "park@example.com"},
	}

	for _, user := range predefinedUsers {
		service.users[user.ID] = user
	}

	return service
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id uint) (*models.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	return user, exists
}

// CreateUser creates a new user
func (s *UserService) CreateUser(req domain.CreateUserRequest) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find next available ID
	nextID := uint(1)
	for id := range s.users {
		if id >= nextID {
			nextID = id + 1
		}
	}

	user := &models.User{
		ID:    nextID,
		Name:  req.Name,
		Email: req.Email,
	}

	s.users[nextID] = user
	return user, nil
}

// GetAllUsers returns all users
func (s *UserService) GetAllUsers() []*models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*models.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}
