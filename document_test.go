package pigeongo

import (
	"encoding/json"
	"testing"

	"github.com/buger/jsonparser"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestApplyChanges(t *testing.T) {
	t.Parallel()

	doc := NewDocument([]byte(`{ "name": "Philipp" }`))
	doc.ApplyChanges(Changes{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/name",
				Value: rawMessage(`"Phil"`),
				Prev:  rawMessage(`"Philipp"`),
			},
		},
		Ts:  2,
		Cid: "50reifj9hyt",
		Gid: "dva96nqsdd",
	})
	// add old message
	doc.ApplyChanges(Changes{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/name",
				Value: rawMessage(`"Hans"`),
				Prev:  rawMessage(`"Philipp"`),
			},
		},
		Ts:  1,
		Cid: "50reifj9hyt",
		Gid: "fhs52fqgdd",
	})
	// add old message
	doc.ApplyChanges(Changes{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/name",
				Value: rawMessage(`"Philipp"`),
				Prev:  rawMessage(`"Dieter"`),
			},
		},
		Ts:  0,
		Cid: "50reifj9hyt",
		Gid: "fhs52fqgdd",
	})
	doc.ApplyChanges(Changes{
		Diff: []Operation{
			{
				Op:   "remove",
				Path: "/notexist",
				Prev: rawMessage(`"Dieter"`),
			},
		},
		Ts:  5,
		Cid: "50reifj9hyt",
		Gid: "hhde2ffcgj",
	})

	assert.Equal(t, `{"name":"Phil"}`, string(doc.JSON()))
	assert.Equal(t, "patch error: error in remove for path: '/notexist': Unable to remove nonexistent key: notexist: missing value", doc.Warning)

	_, err := json.Marshal(doc.History())
	assert.NoError(t, err)

	doc.ReduceHistory(5)
	assert.Len(t, doc.History(), 2)
}

func BenchmarkApplyChanges(b *testing.B) {
	doc := NewDocument([]byte(`{ "name": "Philipp" }`))
	b.ResetTimer()
	lastValue := rawMessage(`"Philipp"`)
	for n := 0; n < b.N; n++ {
		// generate random string
		nextValue := rawMessage(uuid.New().String())
		doc.ApplyChanges(Changes{
			Diff: []Operation{
				{
					Op:    "replace",
					Path:  "/name",
					Value: nextValue,
					Prev:  lastValue,
				},
			},
			Ts:  int64(n),
			Cid: "50reifj9hyt",
			Gid: "dva96nqsdd",
		})
		lastValue = nextValue
	}
}

func TestRawToJSON(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		value []byte
	}{
		{
			value: []byte(`1`),
		},
		{
			value: []byte(`"string"`),
		},
		{
			value: []byte(`{"a": "b"}`),
		},
		{
			value: []byte(`[1,2,3,4]`),
		},
	}

	for _, testCase := range testCases {
		value, dataType, _, _ := jsonparser.Get(testCase.value)
		result := rawToJSON(value, dataType)
		assert.Equal(t, string(testCase.value), string(*result))
	}
}

func TestIdentifiers(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		identifiers [][]string
		path        string
		value       []byte
		expected    string
	}{
		{
			identifiers: nil,
			path:        "/[123]/name",
			value:       []byte(`[{"id": 123, "name": "card1", "value": 1}]`),
			expected:    `[{"id":123,"name":"card2","value":1}]`,
		},
		{
			identifiers: [][]string{{"id"}},
			path:        "/[123]/name",
			value:       []byte(`[{"id": 123, "name": "card1", "value": 1}]`),
			expected:    `[{"id":123,"name":"card2","value":1}]`,
		},
		{
			identifiers: [][]string{{"id"}},
			path:        "/[123]/name",
			value:       []byte(`[{"id": "123", "name": "card1", "value": 1}]`),
			expected:    `[{"id":"123","name":"card2","value":1}]`,
		},
		{
			identifiers: [][]string{{"attrs", "id"}},
			path:        "/[123]/name",
			value:       []byte(`[{"attrs": {"id": 123}, "name": "card1", "value": 1}]`),
			expected:    `[{"attrs":{"id":123},"name":"card2","value": 1}]`,
		},
		{
			identifiers: [][]string{{"attrs", "id"}},
			path:        "/[123]/name",
			value:       []byte(`[{"attrs": {"id": "123"}, "name": "card1", "value": 1}]`),
			expected:    `[{"attrs":{"id":"123"},"name":"card2","value": 1}]`,
		},
		{
			identifiers: [][]string{{"attrs", "id"}},
			path:        "/[123]/name",
			value:       []byte(`[{"id": "123", "name": "card1", "value": 1}]`),
			expected:    `[{"id":"123","name":"card1","value": 1}]`,
		},
		{
			identifiers: [][]string{{"id"}},
			path:        "/[123]/content/[hello]/text",
			value:       []byte(`[{"id": "123", "content": [{"id":"hello", "text":"card1"}, {"id": "foo", "text": "baa"}], "value": 1}]`),
			expected:    `[{"id": "123", "content": [{"id":"hello", "text":"card2"}, {"id": "foo", "text": "baa"}], "value": 1}]`,
		},
		{
			identifiers: [][]string{{"id"}},
			path:        "/[123]/content/[notfound]/text",
			value:       []byte(`[{"id": "123", "content": [{"id":"hello", "text":"card1"}, {"id": "foo", "text": "baa"}], "value": 1}]`),
			expected:    `[{"id": "123", "content": [{"id":"hello", "text":"card1"}, {"id": "foo", "text": "baa"}], "value": 1}]`,
		},
		{
			identifiers: [][]string{{"id"}},
			path:        "/[123]/content/[1234]/text",
			value:       []byte(`[{"id": "123", "content": [{"text":"card1"}, {"text": "baa"}], "value": 1}]`),
			expected:    `[{"id": "123", "content": [{"text":"card1"}, {"text": "baa"}], "value": 1}]`,
		},
		{
			identifiers: [][]string{{"id"}},
			path:        "/[123]/content/1/text",
			value:       []byte(`[{"id": "123", "content": [{"text":"card1"}, {"text": "card1"}], "value": 1}]`),
			expected:    `[{"id": "123", "content": [{"text":"card1"}, {"text": "card2"}], "value": 1}]`,
		},
	}

	for _, testCase := range testCases {
		var opts []DocumentOption
		if testCase.identifiers != nil {
			opts = append(opts, WithIdentifiers(testCase.identifiers))
		}
		doc := NewDocument(testCase.value, opts...)
		if testCase.identifiers != nil {
			assert.Equal(t, testCase.identifiers, doc.identifiers)
		}

		// patch name by identifier
		doc.ApplyChanges(Changes{
			Diff: []Operation{
				{
					Op:    "replace",
					Path:  testCase.path,
					Value: rawMessage(`"card2"`),
					Prev:  rawMessage(`"card1"`),
				},
			},
			Ts:  1,
			Cid: "50reifj9hyt",
			Gid: "dva96nqsdd",
		})

		assert.JSONEq(t, testCase.expected, string(doc.JSON()))
	}
}
