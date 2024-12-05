package pigeongo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		operations []Operation
		expected   []Operation
	}{
		{
			operations: []Operation{{
				Op:    "add",
				Path:  "/name",
				Value: rawMessage(`"Philipp"`),
			}},
			expected: []Operation{{
				Op:   "remove",
				Path: "/name",
				Prev: rawMessage(`"Philipp"`),
			}},
		},
		{
			operations: []Operation{{
				Op:   "remove",
				Path: "/name",
				Prev: rawMessage(`"Philipp"`),
			}},
			expected: []Operation{{
				Op:    "add",
				Path:  "/name",
				Value: rawMessage(`"Philipp"`),
			}},
		},
		{
			operations: []Operation{{
				Op:    "replace",
				Path:  "/name",
				Value: rawMessage(`"Hans"`),
				Prev:  rawMessage(`"Philipp"`),
			}},
			expected: []Operation{{
				Op:    "replace",
				Path:  "/name",
				Value: rawMessage(`"Philipp"`),
				Prev:  rawMessage(`"Hans"`),
			}},
		},
		{
			operations: []Operation{{
				Op:    "add",
				Path:  "/cards/2",
				Value: rawMessage(`{"id": 345, "name": "card2", "value": 2}`),
			}},
			expected: []Operation{{
				Op:   "remove",
				Path: "/cards/[345]",
				Prev: rawMessage(`{"id": 345, "name": "card2", "value": 2}`),
			}},
		},
		{
			operations: []Operation{{
				Op:    "add",
				Path:  "/item/subitem",
				Value: rawMessage(`{"id": 345, "name": "card2", "value": 2}`),
			}},
			expected: []Operation{{
				Op:   "remove",
				Path: "/item/subitem",
				Prev: rawMessage(`{"id": 345, "name": "card2", "value": 2}`),
			}},
		},
		{
			operations: []Operation{{
				Op:   "remove",
				Path: "/cards/[345]",
				Prev: rawMessage(`{"id": 345, "name": "card2", "value": 2}`),
			}},
			expected: []Operation{{
				Op:    "add",
				Path:  "/cards/0",
				Value: rawMessage(`{"id": 345, "name": "card2", "value": 2}`),
			}},
		},
		{
			operations: []Operation{{
				Op:    "add",
				Path:  "/cards/[463]",
				Value: rawMessage(`{"id": 345, "name": "card2", "value": 2}`),
				Prev:  rawMessage(`"bad prev"`),
			}},
			expected: []Operation{{
				Op:   "remove",
				Path: "/cards/[345]",
				Prev: rawMessage(`{"id": 345, "name": "card2", "value": 2}`),
			}},
		},
		{
			operations: []Operation{{
				Op:    "add",
				Path:  "/dashboard/[534345-3435345]/cards/[463]",
				Value: rawMessage(`{"id": 345, "name": "card2", "value": 2}`),
				Prev:  rawMessage(`"bad prev"`),
			}},
			expected: []Operation{{
				Op:   "remove",
				Path: "/dashboard/[534345-3435345]/cards/[345]",
				Prev: rawMessage(`{"id": 345, "name": "card2", "value": 2}`),
			}},
		},
		{
			operations: []Operation{
				{
					Op:    "replace",
					Path:  "/name",
					Value: rawMessage(`"Philipp"`),
					Prev:  rawMessage(`"Dieter"`),
				},
				{
					Op:    "replace",
					Path:  "/name",
					Value: rawMessage(`"Hans"`),
					Prev:  rawMessage(`"Philipp"`),
				},
			},
			expected: []Operation{
				{
					Op:    "replace",
					Path:  "/name",
					Value: rawMessage(`"Philipp"`),
					Prev:  rawMessage(`"Hans"`),
				},
				{
					Op:    "replace",
					Path:  "/name",
					Value: rawMessage(`"Dieter"`),
					Prev:  rawMessage(`"Philipp"`),
				},
			},
		},
	}

	for i, testCase := range testCases {
		reversedOperations := reverse(testCase.operations, [][]string{{"id"}})

		assert.Equal(t, testCase.expected, reversedOperations, fmt.Sprintf("test %d", i))
	}
}

func BenchmarkReverse(b *testing.B) {
	operations := []Operation{
		{
			Op:    "replace",
			Path:  "/name",
			Value: rawMessage(`"Philipp"`),
			Prev:  rawMessage(`"Dieter"`),
		},
		{
			Op:    "replace",
			Path:  "/name",
			Value: rawMessage(`"Hans"`),
			Prev:  rawMessage(`"Philipp"`),
		},
	}

	for i := 0; i < b.N; i++ {
		reverse(operations, [][]string{{"id"}})
	}
}
