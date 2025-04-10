package pigeongo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDuplicateIdentifiers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		data        []byte
		identifiers [][]string
		expectError bool
		errorMsg    string
	}{
		{
			name: "No duplicates in array of objects",
			data: []byte(`[{"id": "1", "name": "Alice"}, {"id": "2", "name": "Bob"}]`),
			identifiers: [][]string{
				{"id"},
			},
			expectError: false,
		},
		{
			name: "Duplicate identifiers in array of objects",
			data: []byte(`[{"id": "1", "name": "Alice"}, {"id": "1", "name": "Bob"}]`),
			identifiers: [][]string{
				{"id"},
			},
			expectError: true,
			errorMsg:    "duplicate identifier found: [1] at index 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDuplicateIdentifiers(tt.data, tt.identifiers)
			if tt.expectError {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errorMsg)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestWalkValidateDuplicateIdentifiers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		data        any
		identifiers [][]string
		expectError bool
		errorMsg    string
	}{
		{
			name: "No duplicates in array of objects",
			data: []any{
				map[string]any{"id": "1", "name": "Alice"},
				map[string]any{"id": "2", "name": "Bob"},
			},
			identifiers: [][]string{{"id"}},
			expectError: false,
		},
		{
			name: "Duplicate identifiers in array of objects",
			data: []any{
				map[string]any{"id": "1", "name": "Alice"},
				map[string]any{"id": "1", "name": "Bob"},
			},
			identifiers: [][]string{{"id"}},
			expectError: true,
			errorMsg:    "duplicate identifier found: [1] at index 1",
		},
		{
			name: "Nested objects with no duplicates",
			data: []any{
				map[string]any{
					"id": "1",
					"children": []any{
						map[string]any{"id": "2"},
						map[string]any{"id": "3"},
					},
				},
			},
			identifiers: [][]string{{"id"}},
			expectError: false,
		},
		{
			name: "Nested objects with duplicates",
			data: []any{
				map[string]any{
					"id": "1",
					"children": []any{
						map[string]any{"id": "2"},
						map[string]any{"id": "2"},
					},
				},
			},
			identifiers: [][]string{{"id"}},
			expectError: true,
			errorMsg:    "duplicate identifier found: [2] at index 1",
		},
		{
			name: "Empty identifiers",
			data: []any{
				map[string]any{"id": "1", "name": "Alice"},
				map[string]any{"id": "2", "name": "Bob"},
			},
			identifiers: [][]string{},
			expectError: false,
		},
		{
			name:        "Empty data",
			data:        []any{},
			identifiers: [][]string{{"id"}},
			expectError: false,
		},
		{
			name: "Invalid data type",
			data: map[string]any{
				"id":       "1",
				"children": "invalid",
			},
			identifiers: [][]string{{"id"}},
			expectError: false,
		},
		{
			name: "Nested identifiers",
			data: []any{
				map[string]any{"reference": map[string]any{"id": "1"}},
				map[string]any{"reference": map[string]any{"id": "2"}},
			},
			identifiers: [][]string{{"reference", "id"}},
			expectError: false,
		},
		{
			name: "Duplicate nested identifiers",
			data: []any{
				map[string]any{"reference": map[string]any{"id": "1"}},
				map[string]any{"reference": map[string]any{"id": "1"}},
			},
			identifiers: [][]string{{"reference", "id"}},
			expectError: true,
			errorMsg:    "duplicate identifier found: [1] at index 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := walkValidateDuplicateIdentifiers(tt.data, tt.identifiers)
			if tt.expectError {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.errorMsg)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
