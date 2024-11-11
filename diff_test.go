package pigeongo

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetID(t *testing.T) {
	t.Parallel()

	identifiers := [][]string{
		{"id"}, {"reference", "id"},
	}

	testCases := []struct {
		description string
		value       string
		expected    string
	}{
		{
			description: "object with id",
			value:       `{"id":"1234567890"}`,
			expected:    "[1234567890]",
		},
		{
			description: "object with number id",
			value:       `{"id": 1234567890}`,
			expected:    "[1234567890]",
		},
		{
			description: "object with float id",
			value:       `{"id": 1234567.890}`,
			expected:    "[1234567.890000]",
		},
		{
			description: "object without id",
			value:       `{"name": "<NAME>"}`,
			expected:    "",
		},
		{
			description: "object child with id",
			value:       `{"reference":{"id": "1234567890"}}`,
			expected:    "[1234567890]",
		},
		{
			description: "object child without id",
			value:       `{"reference":{"name":"<NAME>"}}`,
			expected:    "",
		},
		{
			description: "string",
			value:       "1234567890",
			expected:    "",
		},
		{
			description: "float",
			value:       "1234",
			expected:    "",
		},
	}

	for _, testCase := range testCases {
		var value interface{}
		err := json.Unmarshal([]byte(testCase.value), &value)
		assert.Nil(t, err)

		id := getID(value, identifiers)
		assert.Equal(t, testCase.expected, id, testCase.description)
	}
}

func TestComparePrimitiveSlices(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description     string
		left            []interface{}
		right           []interface{}
		expectedReplace bool
	}{
		{
			description:     "empty slice",
			left:            []interface{}{},
			right:           []interface{}{},
			expectedReplace: false,
		},
		{
			description:     "different length",
			left:            []interface{}{"1"},
			right:           []interface{}{"1", "2"},
			expectedReplace: true,
		},
		{
			description:     "different values",
			left:            []interface{}{"1"},
			right:           []interface{}{"2"},
			expectedReplace: true,
		},
		{
			description:     "same length and values",
			left:            []interface{}{"1"},
			right:           []interface{}{"1"},
			expectedReplace: false,
		},
		{
			description:     "differnet types",
			left:            []interface{}{"1"},
			right:           []interface{}{2345},
			expectedReplace: true,
		},
		{
			description:     "new slice is not primitive",
			left:            []interface{}{"1"},
			right:           []interface{}{[]string{"1"}},
			expectedReplace: true,
		},
	}

	for _, testCase := range testCases {
		ops := []Operation{}
		ops = comparePrimitiveSlices(ops, "/array", testCase.left, testCase.right)

		if testCase.expectedReplace {
			rightBytes, _ := json.Marshal(testCase.right)
			newValue := (*json.RawMessage)(&rightBytes)

			assert.Len(t, ops, 1)
			assert.Equal(t, "/array", ops[0].Path)
			assert.Equal(t, "", ops[0].From)
			assert.Equal(t, "replace", ops[0].Op)
			assert.Equal(t, newValue, ops[0].Value)
		} else {
			assert.Len(t, ops, 0)
		}
	}
}

func TestDiffWithBadPayloads(t *testing.T) {
	t.Parallel()

	ops, err := diff([]byte("bad"), []byte("[]"), [][]string{})
	assert.Nil(t, ops)
	assert.Error(t, err)

	ops, err = diff([]byte("[]"), []byte("bad"), [][]string{})
	assert.Nil(t, ops)
	assert.Error(t, err)
}
