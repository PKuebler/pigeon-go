package pigeongo

import (
	"encoding/json"
	"testing"

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
			Ts:  n,
			Cid: "50reifj9hyt",
			Gid: "dva96nqsdd",
		})
		lastValue = nextValue
	}
}
