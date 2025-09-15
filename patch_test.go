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
			doc:       []byte(`["one", "two", "three"]`),
			patch:     []byte(`[{"op":"add","path":"/1","value":"four"}]`),
			want:      []byte(`["one", "four", "two", "three"]`),
			wantError: false,
		},
		{
			doc:       []byte(`["one", "four", "two", "three"]`),
			patch:     []byte(`[{"op":"remove","path":"/1"}]`),
			want:      []byte(`["one", "two", "three"]`),
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
			doc:       []byte(`{ "body": [{"id":"abc1234","content":[{"id":"child1","text":"content 0"},{"id":"child2","text":"content 1"},{"id":"child3","text":"content 2"}]}] }`),
			patch:     []byte(`[{ "op": "replace", "path": "/body/[abc1234]/content/[child2]/text", "value": "hallo" }]`),
			want:      []byte(`{ "body": [{"id":"abc1234","content":[{"id":"child1","text":"content 0"},{"id":"child2","text":"hallo"},{"id":"child3","text":"content 2"}]}] }`),
			wantError: false,
		},
		{
			doc:       []byte(`{ "body": [{"id":"abc1234","content":[{"id":"child1","text":"content 0"},{"id":"child2","text":"content 1"},{"id":"child3","text":"content 2"}]}] }`),
			patch:     []byte(`[{ "op": "replace", "path": "/body/[abc1234]/content/[child3]/text", "value": "hallo" }]`),
			want:      []byte(`{ "body": [{"id":"abc1234","content":[{"id":"child1","text":"content 0"},{"id":"child2","text":"content 1"},{"id":"child3","text":"hallo"}]}] }`),
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
		{
			doc: []byte(`{"abc":"def","123":null}`),
			patch: []byte(`[
				{"op":"add","path":"/hello"},
				{"op":"add","path":"/foo","value":null},
				{"op":"replace","path":"/abc","value":null},
				{"op":"remove","path":"/123"}
			]`),
			want:      []byte(`{"abc":null,"foo":null,"hello":null}`),
			wantError: false,
		},
	}

	for i, testCase := range testCases {
		operations := []Operation{}

		fmt.Println("")
		fmt.Println("# TESTCASE", i)

		err := json.Unmarshal(testCase.patch, &operations)
		assert.NoError(t, err, fmt.Sprintf("test %d", i))

		result, err := patch(testCase.doc, operations, [][]string{{"id"}})
		fmt.Printf("result: %s\n", string(result))
		if testCase.wantError {
			assert.Error(t, err, fmt.Sprintf("test %d", i))
		} else {
			assert.NoError(t, err, fmt.Sprintf("test %d", i))
		}
		assert.JSONEq(t, string(testCase.want), string(result), fmt.Sprintf("test %d", i))
	}
	t.Error("done")
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

func TestFixEndOfArrayPaths(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		doc          []byte
		path         string
		expectedPath string
	}{
		{
			doc:          []byte(`["a","b","c"]`),
			path:         "/3",
			expectedPath: "/-",
		},
		{
			doc:          []byte(`["a"]`),
			path:         "/1",
			expectedPath: "/-",
		},
		{
			doc:          []byte(`{"array":["a","b","c"]}`),
			path:         "/array/3",
			expectedPath: "/array/-",
		},
		{
			doc:          []byte(`{"array":["a","b","c"]}`),
			path:         "/array/2",
			expectedPath: "/array/2",
		},
		{
			doc:          []byte(`{"object":{"1":"def"}}`),
			path:         "/object/1",
			expectedPath: "/object/1",
		},
		{
			doc:          []byte(`{"array":[{"id":"def"},{"id":"abc"},{"id":"ghj"}]}`),
			path:         "/array/1",
			expectedPath: "/array/1",
		},
		{
			doc:          []byte(`{"array":[{"id":"def"},{"id":"abc"},{"id":"ghj"}]}`),
			path:         "/array/3",
			expectedPath: "/array/-",
		},
	}

	for i, testCase := range testCases {
		patchObj := NewJsonpatchPatch([]Operation{{Op: "add", Path: testCase.path}})
		patchObj = fixEndOfArrayPaths(testCase.doc, patchObj)

		path, err := patchObj[0].Path()
		assert.Nil(t, err)
		assert.Equal(t, testCase.expectedPath, path, "Test case %d failed", i)
	}
}
