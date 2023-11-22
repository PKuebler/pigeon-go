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
		doc       []byte
		patch     []byte
		want      []byte
		wantError bool
	}{
		{
			doc:       []byte(`{ "count": 1 }`),
			patch:     []byte(`[{ "op": "replace", "path": "/count", "value": 2 }]`),
			want:      []byte(`{ "count": 2 }`),
			wantError: false,
		},
		{
			doc:       []byte(`{ "count": 1, "actor": { "name": "bim" } }`),
			patch:     []byte(`[{ "op": "replace", "path": "/actor/name", "value": "bam" }]`),
			want:      []byte(`{ "count": 1, "actor": { "name": "bam" } }`),
			wantError: false,
		},
		{
			doc:       []byte(`[{"id": 1, "name": "betsy"}, {"id": 2, "name": "hank"}]`),
			patch:     []byte(`[{ "op": "replace", "path": "/[2]/name", "value": "henry" }]`),
			want:      []byte(`[{"id": 1, "name": "betsy"}, {"id": 2, "name": "henry"}]`),
			wantError: false,
		},
		{
			doc:       []byte(`[{"id": 1, "name": "betsy"}, {"id": 2, "name": "hank"}]`),
			patch:     []byte(`[{ "op": "replace", "path": "/[5]/name", "value": "henry" }]`),
			want:      []byte(`[{"id": 1, "name": "betsy"}, {"id": 2, "name": "hank"}]`),
			wantError: true,
		},
		{
			doc:       []byte(`[1,2,3,4,5,6]`),
			patch:     []byte(`[{"op":"add","path":"/3","value":33}]`),
			want:      []byte(`[1,2,3,33,4,5,6]`),
			wantError: false,
		},
		{
			doc:       []byte(`[1,2,3,4,5,6]`),
			patch:     []byte(`[{"op":"remove","path":"/3"}]`),
			want:      []byte(`[1,2,3,5,6]`),
			wantError: false,
		},
		{
			doc:       []byte(`[{"id":"def"},{"id":"abc"},{"id":"ghi"}]`),
			patch:     []byte(`[{"op":"move","from":"/[abc]","path":"/0"}]`),
			want:      []byte(`[{"id":"abc"},{"id":"def"},{"id":"ghi"}]`),
			wantError: false,
		},
		{
			doc:       []byte(`[{"id":"def"},{"id":"abc"},{"id":"ghi"}]`),
			patch:     []byte(`[{"op":"move","from":"/1","path":"/0"}]`),
			want:      []byte(`[{"id":"abc"},{"id":"def"},{"id":"ghi"}]`),
			wantError: false,
		},
		{
			doc:       []byte(`[{"id":"abc"},{"id":"def"},{"id":"ghi"}]`),
			patch:     []byte(`[{"op":"move","from":"/[ghi]","path":"/1"},{"op":"move","from":"/[def]","path":"/2"}]`),
			want:      []byte(`[{"id":"abc"},{"id":"ghi"},{"id":"def"}]`),
			wantError: false,
		},
		{
			doc:       []byte(`[{"id":"abc"},{"id":"def"},{"id":"ghi"}]`),
			patch:     []byte(`[{"op":"move","from":"/2","path":"/1"},{"op":"move","from":"/1","path":"/2"}]`),
			want:      []byte(`[{"id":"abc"},{"id":"def"},{"id":"ghi"}]`),
			wantError: false,
		},
		{
			doc:       []byte(`[{"id":"abc"},{"id":"def"},{"id":"ghi"},{"id":"jkl"}]`),
			patch:     []byte(`[{"op":"move","from":"/2","path":"/1"},{"op":"move","from":"/3","path":"/0"}]`),
			want:      []byte(`[{"id":"jkl"},{"id":"abc"},{"id":"ghi"},{"id":"def"}]`),
			wantError: false,
		},
		{
			doc:       []byte(`[{"id":"abc"},{"id":"def"},{"id":"ghi"},{"id":"jkl"}]`),
			patch:     []byte(`[{"op":"move","from":"/[ghi]","path":"/0"},{"op":"move","from":"/[abc]","path":"/1"},{"op":"move","from":"/[def]","path":"/2"}]`),
			want:      []byte(`[{"id":"ghi"},{"id":"abc"},{"id":"def"},{"id":"jkl"}]`),
			wantError: false,
		},
		{
			doc:       []byte(`{"id":"def"}`),
			patch:     []byte(`[{"op":"remove","path":"/email"}]`),
			want:      []byte(`{"id":"def"}`),
			wantError: true,
		},
		{
			doc:   []byte(`{}`),
			patch: []byte(`[{"op":"add","path":"/","value":{"hello":"world"}}]`),
			want:  []byte(`{"":{"hello":"world"}}`),
		},
		{
			doc:   []byte(`{}`),
			patch: []byte(`[{"op":"add","path":"/hello","value":"world"}]`),
			want:  []byte(`{"hello":"world"}`),
		},
		{
			doc:       []byte(`{ "body": [{"id":"abc1234","content":[{"text":"content 0"},{"text":"content 1"},{"text":"content 2"}]}] }`),
			patch:     []byte(`[{ "op": "replace", "path": "/body/[abc1234]/content/2/text", "value": "hallo" }]`),
			want:      []byte(`{ "body": [{"id":"abc1234","content":[{"text":"content 0"},{"text":"content 1"},{"text":"hallo"}]}] }`),
			wantError: false,
		},
		{
			doc: []byte(`{"body":[{"id":"id-3"},{"id":"id-2"},{"id":"id-5"},{"id":"id-6"},{"id":"id-1"},{"id":"id-7"},{"id":"id-4"}]}`),
			patch: []byte(`[
				{"op":"move","from":"/body/[id-1]","path":"/body/0"},
				{"op":"move","from":"/body/[id-3]","path":"/body/1"},
				{"op":"move","from":"/body/[id-2]","path":"/body/2"},
				{"op":"move","from":"/body/[id-5]","path":"/body/3"},
				{"op":"move","from":"/body/[id-6]","path":"/body/4"}
			]`),
			want:      []byte(`{"body":[{"id":"id-1"},{"id":"id-3"},{"id":"id-2"},{"id":"id-5"},{"id":"id-6"},{"id":"id-7"},{"id":"id-4"}]}`),
			wantError: false,
		},
	}

	for i, testCase := range testCases {
		fmt.Println(i, ">", string(testCase.patch))
		operations := []Operation{}

		err := json.Unmarshal(testCase.patch, &operations)
		assert.NoError(t, err, fmt.Sprintf("test %d", i))

		result, err := patch(testCase.doc, operations, [][]string{{"id"}})
		if testCase.wantError {
			assert.Error(t, err, fmt.Sprintf("test %d", i))
		} else {
			assert.NoError(t, err, fmt.Sprintf("test %d", i))
		}
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

		result, err := patch(doc, operations, [][]string{{"id"}})
		if err != nil {
			b.Fatal(err)
		}

		if string(result) != `{"count":2}` {
			b.Fatal("result is not expected")
		}
	}
}
