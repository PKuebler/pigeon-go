package pigeongo

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPatch(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		doc   []byte
		patch []byte
		want  []byte
	}{
		{
			doc:   []byte(`{ "count": 1 }`),
			patch: []byte(`[{ "op": "replace", "path": "/count", "value": 2 }]`),
			want:  []byte(`{ "count": 2 }`),
		},
		{
			doc:   []byte(`{ "count": 1, "actor": { "name": "bim" } }`),
			patch: []byte(`[{ "op": "replace", "path": "/actor/name", "value": "bam" }]`),
			want:  []byte(`{ "count": 1, "actor": { "name": "bam" } }`),
		},
		{
			doc:   []byte(`[{"id": 1, "name": "betsy"}, {"id": 2, "name": "hank"}]`),
			patch: []byte(`[{ "op": "replace", "path": "/[2]/name", "value": "henry" }]`),
			want:  []byte(`[{"id": 1, "name": "betsy"}, {"id": 2, "name": "henry"}]`),
		},
		{
			doc:   []byte(`[1,2,3,4,5,6]`),
			patch: []byte(`[{"op":"add","path":"/3","value":33}]`),
			want:  []byte(`[1,2,3,33,4,5,6]`),
		},
		{
			doc:   []byte(`[1,2,3,4,5,6]`),
			patch: []byte(`[{"op":"remove","path":"/3"}]`),
			want:  []byte(`[1,2,3,5,6]`),
		},
		{
			doc:   []byte(`[{"id":"def"},{"id":"abc"},{"id":"ghi"}]`),
			patch: []byte(`[{"op":"move","from":"/[abc]","path":"/0"}]`),
			want:  []byte(`[{"id":"abc"},{"id":"def"},{"id":"ghi"}]`),
		},
		{
			doc:   []byte(`[{"id":"def"},{"id":"abc"},{"id":"ghi"}]`),
			patch: []byte(`[{"op":"move","from":"/1","path":"/0"}]`),
			want:  []byte(`[{"id":"abc"},{"id":"def"},{"id":"ghi"}]`),
		},
	}

	for i, testCase := range testCases {
		operations := []Operation{}

		err := json.Unmarshal(testCase.patch, &operations)
		assert.NoError(t, err, fmt.Sprintf("test %d", i))

		for _, op := range operations {
			fmt.Println(i, op.Op, op.Path, op.From, op.Value, op.Prev)
		}

		result, err := patch(testCase.doc, operations)
		assert.NoError(t, err, fmt.Sprintf("test %d", i))
		assert.JSONEq(t, string(testCase.want), string(result), fmt.Sprintf("test %d", i))
	}
}

func BenchmarkPatch(b *testing.B) {
	doc := []byte(`{ "count": 1 }`)
	rawPatch := []byte(`[{ "op": "replace", "path": "/count", "value": 2 }]`)

	for i := 0; i < b.N; i++ {
		operations := []Operation{}

		err := json.Unmarshal(rawPatch, &operations)
		if err != nil {
			b.Fatal(err)
		}

		result, err := patch(doc, operations)
		if err != nil {
			b.Fatal(err)
		}

		if string(result) != `{"count":2}` {
			b.Fatal("result is not expected")
		}
	}
}
