package service

import (
	"testing"

	"aoroa/internal/models"
)

func TestNewUserService(t *testing.T) {
	service := NewUserService()

	if service == nil {
		t.Fatal("NewUserService() returned nil")
	}

	// Check if predefined users are initialized
	expectedUsers := []struct {
		id   uint
		name string
	}{
		{1, "김개발"},
		{2, "이디자인"},
		{3, "박기획"},
	}

	for _, expected := range expectedUsers {
		user, exists := service.GetUser(expected.id)
		if !exists {
			t.Errorf("Expected user with ID %d to exist", expected.id)
			continue
		}
		if user.ID != expected.id {
			t.Errorf("Expected user ID %d, got %d", expected.id, user.ID)
		}
		if user.Name != expected.name {
			t.Errorf("Expected user name %s, got %s", expected.name, user.Name)
		}
	}
}

func TestUserServiceGetUser(t *testing.T) {
	service := NewUserService()

	tests := []struct {
		name     string
		userID   uint
		expected bool
	}{
		{
			name:     "Get existing user",
			userID:   1,
			expected: true,
		},
		{
			name:     "Get non-existing user",
			userID:   999,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, exists := service.GetUser(tt.userID)
			if exists != tt.expected {
				t.Errorf("GetUser(%d) exists = %v, want %v", tt.userID, exists, tt.expected)
			}
			if exists && user == nil {
				t.Errorf("GetUser(%d) returned nil user when exists = true", tt.userID)
			}
			if exists && user.ID != tt.userID {
				t.Errorf("GetUser(%d) returned user with ID %d", tt.userID, user.ID)
			}
		})
	}
}

func TestUserServiceGetAllUsers(t *testing.T) {
	service := NewUserService()
	users := service.GetAllUsers()

	expectedCount := 3
	if len(users) != expectedCount {
		t.Errorf("GetAllUsers() returned %d users, want %d", len(users), expectedCount)
	}

	// Check if all expected users are present
	userMap := make(map[uint]*models.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	expectedUsers := []uint{1, 2, 3}
	for _, expectedID := range expectedUsers {
		if _, exists := userMap[expectedID]; !exists {
			t.Errorf("Expected user with ID %d not found in GetAllUsers() result", expectedID)
		}
	}
}
