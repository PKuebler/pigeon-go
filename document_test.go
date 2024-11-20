package pigeongo

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/buger/jsonparser"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestApplyChanges(t *testing.T) {
	t.Parallel()

	doc := NewDocument([]byte(`{ "name": "Philipp" }`))
	err := doc.ApplyChanges(Changes{
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
		Mid: "abcdef",
	})
	assert.Nil(t, err)
	// add old message
	err = doc.ApplyChanges(Changes{
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
	assert.Nil(t, err)
	// add old message
	err = doc.ApplyChanges(Changes{
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
		Mid: "sdfsdf",
	})
	assert.Nil(t, err)
	err = doc.ApplyChanges(Changes{
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
	assert.Equal(t, "patch error: error in remove for path: '/notexist': Unable to remove nonexistent key: notexist: missing value", err.Error())

	_, err = json.Marshal(doc.History())
	assert.NoError(t, err)

	assert.Equal(t, "abcdef", doc.History()[2].Mid, doc.History()[2].Gid)

	assert.Nil(t, doc.ReduceHistory(5))
	assert.Len(t, doc.History(), 1)
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
		identifiers   [][]string
		path          string
		value         []byte
		expected      string
		expectedError string
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
			identifiers:   [][]string{{"attrs", "id"}},
			path:          "/[123]/name",
			value:         []byte(`[{"id": "123", "name": "card1", "value": 1}]`),
			expected:      `[{"id":"123","name":"card1","value": 1}]`,
			expectedError: "patch error: id `123` not found",
		},
		{
			identifiers: [][]string{{"id"}},
			path:        "/[123]/content/[hello]/text",
			value:       []byte(`[{"id": "123", "content": [{"id":"hello", "text":"card1"}, {"id": "foo", "text": "baa"}], "value": 1}]`),
			expected:    `[{"id": "123", "content": [{"id":"hello", "text":"card2"}, {"id": "foo", "text": "baa"}], "value": 1}]`,
		},
		{
			identifiers:   [][]string{{"id"}},
			path:          "/[123]/content/[notfound]/text",
			value:         []byte(`[{"id": "123", "content": [{"id":"hello", "text":"card1"}, {"id": "foo", "text": "baa"}], "value": 1}]`),
			expected:      `[{"id": "123", "content": [{"id":"hello", "text":"card1"}, {"id": "foo", "text": "baa"}], "value": 1}]`,
			expectedError: "patch error: id `notfound` not found",
		},
		{
			identifiers:   [][]string{{"id"}},
			path:          "/[123]/content/[1234]/text",
			value:         []byte(`[{"id": "123", "content": [{"text":"card1"}, {"text": "baa"}], "value": 1}]`),
			expected:      `[{"id": "123", "content": [{"text":"card1"}, {"text": "baa"}], "value": 1}]`,
			expectedError: "patch error: id `1234` not found",
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
		err := doc.ApplyChanges(Changes{
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

		if testCase.expectedError == "" {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, testCase.expectedError, err.Error())
		}
		assert.JSONEq(t, testCase.expected, string(doc.JSON()))
	}
}

func TestOptions(t *testing.T) {
	t.Parallel()

	now := time.Now()
	doc := NewDocument(
		[]byte(`{"id": "123"}`),
		WithInitialIDs("client-id", "g-id"),
		WithInitialTime(now),
	)

	firstEntry := doc.History()[0]
	assert.Equal(t, now.UnixMilli(), firstEntry.Ts)
	assert.Equal(t, "client-id", firstEntry.Cid)
	assert.Equal(t, "g-id", firstEntry.Gid)
}

func TestClone(t *testing.T) {
	t.Parallel()

	doc := NewDocument([]byte(`{"id": "123"}`))

	err := doc.ApplyChanges(Changes{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/id",
				Value: rawMessage(`"card2"`),
				Prev:  rawMessage(`"123"`),
			},
		},
		Ts:  1,
		Cid: "50reifj9hyt",
		Gid: "dva96nqsdd",
	})
	assert.Nil(t, err)
	err = doc.ApplyChanges(Changes{
		Diff: []Operation{
			{
				Op:    "add",
				Path:  "/foo",
				Value: rawMessage(`"baa"`),
			},
		},
		Ts:  50,
		Cid: "50reifj9hyt",
		Gid: "jdva96nqsdd",
	})
	assert.Nil(t, err)

	clone := doc.Clone()
	assert.Equal(t, doc.JSON(), clone.JSON())
	assert.Equal(t, doc.History(), clone.History())

	assert.Nil(t, doc.ReduceHistory(100))
	assert.Len(t, doc.History(), 1)
	assert.Len(t, clone.History(), 3)

	err = doc.ApplyChanges(Changes{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/id",
				Value: rawMessage(`"card3"`),
				Prev:  rawMessage(`"card2"`),
			},
		},
		Ts:  110,
		Cid: "50reifj9hyt",
		Gid: "hre46nqsdd",
	})
	assert.Nil(t, err)
	assert.NotEqual(t, string(doc.JSON()), string(clone.JSON()))

	// don't find doc gid in clone gid list
	_, found := clone.gids["hre46nqsdd"]
	assert.False(t, found)

	assert.Nil(t, doc.ReduceHistory(108))
	assert.Len(t, doc.History(), 2)
	assert.Equal(t, "hre46nqsdd", doc.History()[1].Gid)
}

func TestReduceHistory(t *testing.T) {
	t.Parallel()

	now := time.Now().Add(-1 * time.Hour)
	doc := NewDocument([]byte(`{"id": "123"}`), WithInitialIDs("client-id", "change-id"), WithInitialMid("msg-id"), WithInitialTime(now))

	assert.Nil(t, doc.ReduceHistory(2000))
	assert.Len(t, doc.History(), 1)
	assert.Equal(t, "change-id", doc.History()[0].Gid)
	assert.Equal(t, "client-id", doc.History()[0].Cid)
	assert.Equal(t, "msg-id", doc.History()[0].Mid)

	err := doc.ApplyChanges(Changes{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/id",
				Value: rawMessage(`"card2"`),
				Prev:  rawMessage(`"123"`),
			},
		},
		Ts:  now.Add(10 * time.Second).UnixMilli(),
		Cid: "50reifj9hyt",
		Gid: "dva96nqsdd",
		Mid: "50reifj9hyt-dva96nqsdd",
	})
	assert.Nil(t, err)
	err = doc.ApplyChanges(Changes{
		Diff: []Operation{
			{
				Op:    "add",
				Path:  "/foo",
				Value: rawMessage(`"baa"`),
			},
		},
		Ts:  now.Add(60 * time.Second).UnixMilli(),
		Cid: "50reifj9hyt",
		Gid: "jdva96nqsdd",
		Mid: "50reifj9hyt-jdva96nqsdd",
	})
	assert.Nil(t, err)

	clone := doc.Clone()
	assert.Nil(t, clone.ReduceHistory(now.Add(20*time.Second).UnixMilli()))
	assert.Len(t, clone.History(), 2)
	assert.Equal(t, "dva96nqsdd", clone.History()[0].Gid)
	assert.Equal(t, "50reifj9hyt-dva96nqsdd", clone.History()[0].Mid)

	// reduce with only the initial diff
	clone = doc.Clone()
	assert.Nil(t, clone.ReduceHistory(now.Add(2*time.Hour).UnixMilli()))
	assert.Len(t, clone.History(), 1)
	assert.Equal(t, "jdva96nqsdd", clone.History()[0].Gid)
	assert.Equal(t, "50reifj9hyt-jdva96nqsdd", clone.History()[0].Mid)

	assert.Nil(t, clone.ReduceHistory(now.Add(4*time.Hour).UnixMilli()))
	assert.Len(t, clone.History(), 1)
	assert.Equal(t, "jdva96nqsdd", clone.History()[0].Gid)
	assert.Equal(t, "50reifj9hyt-jdva96nqsdd", clone.History()[0].Mid)

	// reduce with history, but nothing to reduce
	clone = doc.Clone()
	assert.Nil(t, clone.ReduceHistory(0))
	assert.Len(t, clone.History(), 3)
	assert.Equal(t, "change-id", clone.History()[0].Gid)
	assert.Equal(t, "client-id", clone.History()[0].Cid)
	assert.Equal(t, "msg-id", clone.History()[0].Mid)
	assert.Equal(t, now.UnixMilli(), clone.History()[0].Ts)
}

func TestFastForwardChanges(t *testing.T) {
	t.Parallel()

	now := time.Now()
	doc := NewDocument([]byte(`{"id":"card1"}`))

	err := doc.ApplyChanges(Changes{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/id",
				Value: rawMessage(`"card2"`),
				Prev:  rawMessage(`"card1"`),
			},
		},
		Ts:  now.Add(10 * time.Second).UnixMilli(),
		Cid: "50reifj9hyt",
		Gid: "dva96nqsdd",
	})
	assert.Nil(t, err)

	err = doc.ApplyChanges(Changes{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/id",
				Value: rawMessage(`"card3"`),
				Prev:  rawMessage(`"card2"`),
			},
		},
		Ts:  now.Add(20 * time.Second).UnixMilli(),
		Cid: "50reifj9hyt",
		Gid: "gdva96nqsdd",
	})
	assert.Nil(t, err)

	err = doc.ApplyChanges(Changes{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/id",
				Value: rawMessage(`"card4"`),
				Prev:  rawMessage(`"card3"`),
			},
		},
		Ts:  now.Add(30 * time.Second).UnixMilli(),
		Cid: "50reifj9hyt",
		Gid: "hdva96nqsdd",
	})
	assert.Nil(t, err)

	err = doc.ApplyChanges(Changes{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/id",
				Value: rawMessage(`"card16"`),
				Prev:  rawMessage(`"card1"`),
			},
		},
		Ts:  now.UnixMilli(),
		Cid: "50reifj9hyt",
		Gid: "udva96nqsdd",
	})
	assert.Nil(t, err)

	ids := []string{}
	for _, change := range doc.History() {
		ids = append(ids, change.Gid)
	}
	assert.Equal(t, []string{"0", "udva96nqsdd", "dva96nqsdd", "gdva96nqsdd", "hdva96nqsdd"}, ids)
	assert.Equal(t, `{"id":"card4"}`, string(doc.JSON()))
}

func TestWrongPrev(t *testing.T) {
	t.Parallel()

	now := time.Now()
	doc := NewDocument([]byte(`{"id":"card1"}`))

	err := doc.ApplyChanges(Changes{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/id",
				Value: rawMessage(`"card2"`),
				Prev:  rawMessage(`"baa"`),
			},
		},
		Ts:  now.Add(10 * time.Minute).UnixMilli(),
		Cid: "50reifj9hyt",
		Gid: "dva96nqsdd",
	})
	assert.Nil(t, err)
	err = doc.ApplyChanges(Changes{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/id",
				Value: rawMessage(`"card3"`),
				Prev:  rawMessage(`"foo"`),
			},
		},
		Ts:  now.Add(20 * time.Minute).UnixMilli(),
		Cid: "50reifj9hyt",
		Gid: "gdva96nqsdd",
	})
	assert.Nil(t, err)
	err = doc.ApplyChanges(Changes{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/id",
				Value: rawMessage(`"card4"`),
				Prev:  rawMessage(`"other value"`),
			},
		},
		Ts:  now.Add(15 * time.Minute).UnixMilli(),
		Cid: "50reifj9hyt",
		Gid: "hdva96nqsdd",
	})
	assert.Nil(t, err)
	err = doc.ApplyChanges(Changes{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/id",
				Value: rawMessage(`"card5"`),
				Prev:  rawMessage(`"oh"`),
			},
		},
		Ts:  now.Add(18 * time.Minute).UnixMilli(),
		Cid: "50reifj9hyt",
		Gid: "ba96nqsdd",
	})
	assert.Nil(t, err)
	fmt.Println("------")
	assert.Equal(t, `{"id":"card3"}`, string(doc.JSON()))
	assert.Nil(t, doc.RewindChanges(now.Add(16*time.Minute).UnixMilli(), ""))
	assert.Equal(t, `{"id":"card4"}`, string(doc.JSON()))
	assert.Nil(t, doc.FastForwardChanges())
	assert.Equal(t, `{"id":"card3"}`, string(doc.JSON()))
}

func TestRemove(t *testing.T) {
	t.Parallel()

	now := time.Now()
	doc := NewDocument([]byte(`{"cards":[{"id":"card1"}]}`))
	err := doc.ApplyChanges(Changes{
		Diff: []Operation{
			{
				Op:   "remove",
				Path: "/cards/[card132]",
			},
		},
		Ts:  now.Add(5 * time.Minute).UnixMilli(),
		Cid: "50reifj9hyt",
		Gid: "dva96nqsdd",
	})

	assert.Len(t, doc.history, 1)
	assert.Equal(t, "patch error: id `card132` not found", err.Error())
}

func TestGetValue(t *testing.T) {
	t.Parallel()

	doc := NewDocument([]byte(`{"id":"card1", "number": 1234, "boolean": true, "object": {"foo": "bar"},  "array": ["one", "two"], "complex": [{"id":"1234","foo":"baa"}]}`))

	assert.Equal(t, `"card1"`, string(*doc.getValue("/id")))
	assert.Equal(t, `1234`, string(*doc.getValue("/number")))
	assert.Equal(t, `true`, string(*doc.getValue("/boolean")))
	assert.Equal(t, `{"foo": "bar"}`, string(*doc.getValue("/object")))
	assert.Equal(t, `["one", "two"]`, string(*doc.getValue("/array")))
	assert.Equal(t, `"bar"`, string(*doc.getValue("/object/foo")))
	assert.Equal(t, `"baa"`, string(*doc.getValue("/complex/[1234]/foo")))
}

func TestDiff(t *testing.T) {
	t.Parallel()

	rawA := []byte(`{"id":"card1", "number": 5678, "movearray": [{"id":"1234","value":"car"},{"id":"4734","value":"people"},{"id":"6523","value":"wood"}], "objectarray": [{"id":"1234","value":"car"},{"id":"4734","value":"people"}], "deepobject":{"id":"1423","text":"foo","childs":[{"id":"64334","value":"baa"}]}}`)
	rawB := []byte(`{"id":"card1", "number": 1234, "movearray": [{"id":"6523","value":"wood"},{"id":"1234","value":"car"},{"id":"4734","value":"people"}], "objectarray": [{"id":"1234","value":"radio"},{"id":"5321","value":"home"},{"id":"4734","value":"people"}], "deepobject":{"id":"1423","text":"foo","childs":[{"id":"64334","value":"foo"}]},"boolean":true,"object":{"foo":"bar"},"array":["one","two"],"complex":[{"id":"1234","foo":"baa"}]}`)
	docA := NewDocument(rawA)
	docB := NewDocument(rawB)

	changes, err := docA.Diff(docB)
	changes.Ts = time.Now().UnixMilli()
	changes.Cid = "50reifj9hyt"
	changes.Gid = "dva96nqsdd"
	assert.Nil(t, err)

	err = docA.ApplyChanges(changes)
	assert.Nil(t, err)

	assert.Equal(t, string(sortKeys(rawB)), string(sortKeys(docA.JSON())))
}

// sortKeys for the assert functions
func sortKeys(doc []byte) []byte {
	obj := map[string]interface{}{}
	_ = json.Unmarshal(doc, &obj)
	sorted, _ := json.Marshal(obj)
	return sorted
}
