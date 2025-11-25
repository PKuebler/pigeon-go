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
		var value any
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
		left            []any
		right           []any
		expectedReplace bool
	}{
		{
			description:     "empty slice",
			left:            []any{},
			right:           []any{},
			expectedReplace: false,
		},
		{
			description:     "different length",
			left:            []any{"1"},
			right:           []any{"1", "2"},
			expectedReplace: true,
		},
		{
			description:     "different values",
			left:            []any{"1"},
			right:           []any{"2"},
			expectedReplace: true,
		},
		{
			description:     "same length and values",
			left:            []any{"1"},
			right:           []any{"1"},
			expectedReplace: false,
		},
		{
			description:     "different types",
			left:            []any{"1"},
			right:           []any{2345},
			expectedReplace: true,
		},
		{
			description:     "new slice is not primitive",
			left:            []any{"1"},
			right:           []any{[]string{"1"}},
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

func TestDiffWithNull(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description string
		payloadA    string
		payloadB    string
		expected    string
	}{
		{
			description: "null values in right should be removed, not added",
			payloadA:    `{"b": {"b1": "value", "c1": "test"}, "c": "test"}`,
			payloadB:    `{"a": null, "b": {"a1": null, "b1": "value", "c1": null}, "c": null}`,
			expected:    `[{"op":"remove","path":"/b/c1","_prev":"test"},{"op":"remove","path":"/c","_prev":"test"}]`,
		},
		{
			description: "new null key should not create add operation",
			payloadA:    `{"a": "value"}`,
			payloadB:    `{"a": "value", "b": null}`,
			expected:    `[]`,
		},
		{
			description: "existing value to null should create remove",
			payloadA:    `{"a": "value"}`,
			payloadB:    `{"a": null}`,
			expected:    `[{"op":"remove","path":"/a","_prev":"value"}]`,
		},
		{
			description: "null to value should create add",
			payloadA:    `{"a": null}`,
			payloadB:    `{"a": "value"}`,
			expected:    `[{"op":"add","path":"/a","value":"value"}]`,
		},
		{
			description: "null to null should not create any operation",
			payloadA:    `{"a": null}`,
			payloadB:    `{"a": null}`,
			expected:    `[]`,
		},
		{
			description: "deeply nested null values",
			payloadA:    `{"a": {"b": {"c": "value"}}}`,
			payloadB:    `{"a": {"b": {"c": null, "d": null}}}`,
			expected:    `[{"op":"remove","path":"/a/b/c","_prev":"value"}]`,
		},
		{
			description: "empty objects remain unchanged",
			payloadA:    `{}`,
			payloadB:    `{"a": null}`,
			expected:    `[]`,
		},
		{
			description: "multiple null keys should all be ignored",
			payloadA:    `{}`,
			payloadB:    `{"a": null, "b": null, "c": null}`,
			expected:    `[]`,
		},
		{
			description: "mix of null and real values",
			payloadA:    `{"existing": "old"}`,
			payloadB:    `{"existing": "new", "null_key": null, "new_key": "value"}`,
			expected:    `[{"op":"replace","path":"/existing","value":"new","_prev":"old"},{"op":"add","path":"/new_key","value":"value"}]`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()
			ops, err := diff([]byte(testCase.payloadA), []byte(testCase.payloadB), [][]string{{"id"}, {"reference", "id"}})
			assert.Nil(t, err)

			b, _ := json.Marshal(ops)
			assert.Equal(t, testCase.expected, string(b), testCase.description)
		})
	}
}
