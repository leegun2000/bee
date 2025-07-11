package domain

import "testing"

func TestIsValidStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{
			name:     "Valid PENDING status",
			status:   StatusPending,
			expected: true,
		},
		{
			name:     "Valid IN_PROGRESS status",
			status:   StatusInProgress,
			expected: true,
		},
		{
			name:     "Valid COMPLETED status",
			status:   StatusCompleted,
			expected: true,
		},
		{
			name:     "Valid CANCELLED status",
			status:   StatusCancelled,
			expected: true,
		},
		{
			name:     "Invalid status",
			status:   "INVALID_STATUS",
			expected: false,
		},
		{
			name:     "Empty status",
			status:   "",
			expected: false,
		},
		{
			name:     "Lowercase status",
			status:   "pending",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidStatus(tt.status)
			if result != tt.expected {
				t.Errorf("IsValidStatus(%s) = %v, want %v", tt.status, result, tt.expected)
			}
		})
	}
}
