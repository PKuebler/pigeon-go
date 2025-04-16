package pigeongo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDuplicateIdentifiers(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		data        []byte
		identifiers [][]string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "No identifiers",
			data:        []byte(`[{"id": "1", "name": "Alice"}, {"id": "2", "name": "Bob"}]`),
			identifiers: nil,
			expectError: false,
		},
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
		{
			name: "Bad JSON",
			data: []byte(`{"id": "1", "name": "Alice"`),
			identifiers: [][]string{
				{"id"},
			},
			expectError: true,
			errorMsg:    "unexpected end of JSON input",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := validateDuplicateIdentifiers(testCase.data, testCase.identifiers)
			if testCase.expectError {
				assert.NotNil(t, err)
				assert.EqualError(t, err, testCase.errorMsg)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestWalkValidateDuplicateIdentifiers(t *testing.T) {
	t.Parallel()

	testCases := []struct {
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
			name: "No identifiers",
			data: []any{
				map[string]any{"id": "1", "name": "Alice"},
				map[string]any{"id": "2", "name": "Bob"},
			},
			identifiers: nil,
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
		{
			name: "Duplicate ids in nested object",
			data: map[string]any{
				"children": []any{
					map[string]any{"id": "1"},
					map[string]any{"id": "1"},
				},
			},
			identifiers: [][]string{{"id"}},
			expectError: true,
			errorMsg:    "duplicate identifier found: [1] at index 1",
		},
		{
			name: "Duplicate ids in nested array",
			data: []any{
				[]any{
					map[string]any{"id": "1"},
					map[string]any{"id": "1"},
				},
			},
			identifiers: [][]string{{"id"}},
			expectError: true,
			errorMsg:    "duplicate identifier found: [1] at index 1",
		},
		{
			name: "Empty ids",
			data: []any{
				map[string]any{"name": "Alice"},
				map[string]any{"name": "Bob"},
			},
			identifiers: [][]string{{"id"}},
			expectError: true,
			errorMsg:    "duplicate identifier found: missing id at index 1",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := walkValidateDuplicateIdentifiers(testCase.data, testCase.identifiers)
			if testCase.expectError {
				assert.NotNil(t, err)
				assert.EqualError(t, err, testCase.errorMsg)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
