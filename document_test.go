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

func TestCreateDocument(t *testing.T) {
	t.Parallel()

	// test with invalid JSON
	doc, err := NewDocument([]byte(`{"id": "123", "name": "card1", "value": 1`))
	assert.NotNil(t, err)
	assert.Nil(t, doc)

	// test with duplicate ids
	doc, err = NewDocument([]byte(`[{"id": "123", "name": "card1", "value": 1}, {"id": "123", "name": "card2", "value": 2}]`))
	assert.NotNil(t, err)
	assert.Nil(t, doc)

	// test with string
	doc, err = NewDocument([]byte(`"123"`))
	assert.Nil(t, err)
	assert.NotNil(t, doc)
}

func TestApplyChange(t *testing.T) {
	t.Parallel()

	doc, err := NewDocument([]byte(`{ "name": "Philipp" }`))
	assert.Nil(t, err)
	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/name",
				Value: rawMessage(`"Phil"`),
				Prev:  rawMessage(`"Philipp"`),
			},
		},
		TimestampMillis: 2,
		ClientID:        "50reifj9hyt",
		ChangeID:        "dva96nqsdd",
		MessageID:       "abcdef",
	})
	assert.Nil(t, err)
	// add old message
	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/name",
				Value: rawMessage(`"Hans"`),
				Prev:  rawMessage(`"Philipp"`),
			},
		},
		TimestampMillis: 1,
		ClientID:        "50reifj9hyt",
		ChangeID:        "fhs52fqgdd",
	})
	assert.Nil(t, err)
	// add old message
	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/name",
				Value: rawMessage(`"Philipp"`),
				Prev:  rawMessage(`"Dieter"`),
			},
		},
		TimestampMillis: 0,
		ClientID:        "50reifj9hyt",
		ChangeID:        "fhs52fqgdd",
		MessageID:       "sdfsdf",
	})
	assert.Nil(t, err)
	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:   "remove",
				Path: "/notexist",
				Prev: rawMessage(`"Dieter"`),
			},
		},
		TimestampMillis: 5,
		ClientID:        "50reifj9hyt",
		ChangeID:        "hhde2ffcgj",
	})

	assert.Equal(t, `{"name":"Phil"}`, string(doc.JSON()))
	assert.Equal(t, "patch error: can't apply changeID hhde2ffcgj: error in remove for path: '/notexist': Unable to remove nonexistent key: notexist: missing value", err.Error())

	_, err = json.Marshal(doc.History())
	assert.NoError(t, err)

	assert.Equal(t, "abcdef", doc.History()[2].MessageID, doc.History()[2].ChangeID)

	assert.Nil(t, doc.ReduceHistory(5))
	assert.Len(t, doc.History(), 1)
}

func BenchmarkApplyChange(b *testing.B) {
	doc, err := NewDocument([]byte(`{ "name": "Philipp" }`))
	assert.Nil(b, err)
	b.ResetTimer()
	lastValue := rawMessage(`"Philipp"`)
	for n := 0; n < b.N; n++ {
		// generate random string
		nextValue := rawMessage(uuid.New().String())
		_ = doc.ApplyChange(Change{
			Diff: []Operation{
				{
					Op:    "replace",
					Path:  "/name",
					Value: nextValue,
					Prev:  lastValue,
				},
			},
			TimestampMillis: int64(n),
			ClientID:        "50reifj9hyt",
			ChangeID:        "dva96nqsdd",
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
			expectedError: "patch error: can't apply changeID dva96nqsdd: id `123` not found",
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
			expectedError: "patch error: can't apply changeID dva96nqsdd: id `notfound` not found",
		},
		{
			identifiers:   [][]string{{"id"}},
			path:          "/[123]/content/[1234]/text",
			value:         []byte(`[{"id": "123", "content": [{"id":"1", "text":"card1"}, {"id":"2", "text": "baa"}], "value": 1}]`),
			expected:      `[{"id": "123", "content": [{"id":"1", "text":"card1"}, {"id":"2", "text": "baa"}], "value": 1}]`,
			expectedError: "patch error: can't apply changeID dva96nqsdd: id `1234` not found",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase %d", i), func(t *testing.T) {
			t.Parallel()

			var opts []DocumentOption
			if testCase.identifiers != nil {
				opts = append(opts, WithIdentifiers(testCase.identifiers))
			}
			doc, err := NewDocument(testCase.value, opts...)
			assert.Nil(t, err)
			if testCase.identifiers != nil {
				assert.Equal(t, testCase.identifiers, doc.identifiers)
			}

			// patch name by identifier
			err = doc.ApplyChange(Change{
				Diff: []Operation{
					{
						Op:    "replace",
						Path:  testCase.path,
						Value: rawMessage(`"card2"`),
						Prev:  rawMessage(`"card1"`),
					},
				},
				TimestampMillis: 1,
				ClientID:        "50reifj9hyt",
				ChangeID:        "dva96nqsdd",
			})

			if testCase.expectedError == "" {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, testCase.expectedError, err.Error())
			}
			assert.JSONEq(t, testCase.expected, string(doc.JSON()))
		})
	}
}

func TestOptions(t *testing.T) {
	t.Parallel()

	now := time.Now()
	doc, err := NewDocument(
		[]byte(`{"id": "123"}`),
		WithInitialIDs("client-id", "g-id"),
		WithInitialTime(now),
	)
	assert.Nil(t, err)

	firstEntry := doc.History()[0]
	assert.Equal(t, now.UnixMilli(), firstEntry.TimestampMillis)
	assert.Equal(t, "client-id", firstEntry.ClientID)
	assert.Equal(t, "g-id", firstEntry.ChangeID)
}

func TestClone(t *testing.T) {
	t.Parallel()

	doc, err := NewDocument([]byte(`{"id": "123"}`))
	assert.Nil(t, err)

	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/id",
				Value: rawMessage(`"card2"`),
				Prev:  rawMessage(`"123"`),
			},
		},
		TimestampMillis: 1,
		ClientID:        "50reifj9hyt",
		ChangeID:        "dva96nqsdd",
	})
	assert.Nil(t, err)
	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:    "add",
				Path:  "/foo",
				Value: rawMessage(`"baa"`),
			},
		},
		TimestampMillis: 50,
		ClientID:        "50reifj9hyt",
		ChangeID:        "jdva96nqsdd",
	})
	assert.Nil(t, err)

	clone := doc.Clone()
	assert.Equal(t, doc.JSON(), clone.JSON())
	assert.Equal(t, doc.History(), clone.History())

	assert.Nil(t, doc.ReduceHistory(100))
	assert.Len(t, doc.History(), 1)
	assert.Len(t, clone.History(), 3)

	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/id",
				Value: rawMessage(`"card3"`),
				Prev:  rawMessage(`"card2"`),
			},
		},
		TimestampMillis: 110,
		ClientID:        "50reifj9hyt",
		ChangeID:        "hre46nqsdd",
	})
	assert.Nil(t, err)
	assert.NotEqual(t, string(doc.JSON()), string(clone.JSON()))

	// don't find doc changeID in clone changeID list
	_, found := clone.changeIDs["hre46nqsdd"]
	assert.False(t, found)

	assert.Nil(t, doc.ReduceHistory(108))
	assert.Len(t, doc.History(), 2)
	assert.Equal(t, "hre46nqsdd", doc.History()[1].ChangeID)
}

func TestReduceHistory(t *testing.T) {
	t.Parallel()

	now := time.Now().Add(-1 * time.Hour)
	doc, err := NewDocument([]byte(`{"id": "123"}`), WithInitialIDs("client-id", "change-id"), WithInitialMessageID("msg-id"), WithInitialTime(now))
	assert.Nil(t, err)

	assert.Nil(t, doc.ReduceHistory(2000))
	assert.Len(t, doc.History(), 1)
	assert.Equal(t, "change-id", doc.History()[0].ChangeID)
	assert.Equal(t, "client-id", doc.History()[0].ClientID)
	assert.Equal(t, "msg-id", doc.History()[0].MessageID)

	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/id",
				Value: rawMessage(`"card2"`),
				Prev:  rawMessage(`"123"`),
			},
		},
		TimestampMillis: now.Add(10 * time.Second).UnixMilli(),
		ClientID:        "50reifj9hyt",
		ChangeID:        "dva96nqsdd",
		MessageID:       "50reifj9hyt-dva96nqsdd",
	})
	assert.Nil(t, err)
	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:    "add",
				Path:  "/foo",
				Value: rawMessage(`"baa"`),
			},
		},
		TimestampMillis: now.Add(60 * time.Second).UnixMilli(),
		ClientID:        "50reifj9hyt",
		ChangeID:        "jdva96nqsdd",
		MessageID:       "50reifj9hyt-jdva96nqsdd",
	})
	assert.Nil(t, err)

	clone := doc.Clone()
	assert.Nil(t, clone.ReduceHistory(now.Add(20*time.Second).UnixMilli()))
	assert.Len(t, clone.History(), 2)
	assert.Equal(t, "dva96nqsdd", clone.History()[0].ChangeID)
	assert.Equal(t, "50reifj9hyt-dva96nqsdd", clone.History()[0].MessageID)

	// reduce with only the initial diff
	clone = doc.Clone()
	assert.Nil(t, clone.ReduceHistory(now.Add(2*time.Hour).UnixMilli()))
	assert.Len(t, clone.History(), 1)
	assert.Equal(t, "jdva96nqsdd", clone.History()[0].ChangeID)
	assert.Equal(t, "50reifj9hyt-jdva96nqsdd", clone.History()[0].MessageID)

	assert.Nil(t, clone.ReduceHistory(now.Add(4*time.Hour).UnixMilli()))
	assert.Len(t, clone.History(), 1)
	assert.Equal(t, "jdva96nqsdd", clone.History()[0].ChangeID)
	assert.Equal(t, "50reifj9hyt-jdva96nqsdd", clone.History()[0].MessageID)

	// reduce with history, but nothing to reduce
	clone = doc.Clone()
	assert.Nil(t, clone.ReduceHistory(0))
	assert.Len(t, clone.History(), 3)
	assert.Equal(t, "change-id", clone.History()[0].ChangeID)
	assert.Equal(t, "client-id", clone.History()[0].ClientID)
	assert.Equal(t, "msg-id", clone.History()[0].MessageID)
	assert.Equal(t, now.UnixMilli(), clone.History()[0].TimestampMillis)
}

func TestReduceHistoryWithIDs(t *testing.T) {
	t.Parallel()

	now := time.Now().Add(-1 * time.Hour)
	doc, err := NewDocument([]byte(`{"id": "123", "body":[{"id":"card1","value":"test"},{"id":"card2","value":"baa"}]}`), WithInitialIDs("client-id", "change-id"), WithInitialMessageID("msg-id"), WithInitialTime(now))
	assert.Nil(t, err)

	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:    "add",
				Path:  "/body/[card2]",
				Value: rawMessage(`{"id":"card3","value":"tttt"}`),
			},
		},
		TimestampMillis: now.Add(10 * time.Second).UnixMilli(),
		ClientID:        "50reifj9hyt",
		ChangeID:        "dva96nqsdd",
		MessageID:       "50reifj9hyt-dva96nqsdd",
	})
	assert.Nil(t, err)
	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/body/[card3]/value",
				Value: rawMessage(`"bar"`),
			},
		},
		TimestampMillis: now.Add(40 * time.Second).UnixMilli(),
		ClientID:        "50reifj9hyt",
		ChangeID:        "hreifj9hyt",
		MessageID:       "50reifj9hyt-hreifj9hyt",
	})
	assert.Nil(t, err)
	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/body/1/value",
				Value: rawMessage(`"retert"`),
			},
		},
		TimestampMillis: now.Add(40 * time.Second).UnixMilli(),
		ClientID:        "50reifj9hyt",
		ChangeID:        "75reifj9hyt",
		MessageID:       "50reifj9hyt-75reifj9hyt",
	})
	assert.Nil(t, err)
	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:   "remove",
				Path: "/body/[card3]",
			},
		},
		TimestampMillis: now.Add(60 * time.Second).UnixMilli(),
		ClientID:        "50reifj9hyt",
		ChangeID:        "jdva96nqsdd",
		MessageID:       "50reifj9hyt-jdva96nqsdd",
	})
	assert.Nil(t, err)
	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:   "remove",
				Path: "/body/[card2]",
			},
		},
		TimestampMillis: now.Add(65 * time.Second).UnixMilli(),
		ClientID:        "50reifj9hyt",
		ChangeID:        "5gdva96nqsdd",
		MessageID:       "50reifj9hyt-5gdva96nqsdd",
	})
	assert.Nil(t, err)
	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:    "add",
				Path:  "/body/[card1]",
				Value: rawMessage(`{"id":"card4","value":"ooooo"}`),
			},
		},
		TimestampMillis: now.Add(75 * time.Second).UnixMilli(),
		ClientID:        "50reifj9hyt",
		ChangeID:        "6gdva96nqsdd",
		MessageID:       "50reifj9hyt-6gdva96nqsdd",
	})
	assert.Nil(t, err)
	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:    "add",
				Path:  "/body/0",
				Value: rawMessage(`{"id":"card5","value":"gggggg"}`),
			},
		},
		TimestampMillis: now.Add(80 * time.Second).UnixMilli(),
		ClientID:        "50reifj9hyt",
		ChangeID:        "7gdva96nqsdd",
		MessageID:       "50reifj9hyt-7gdva96nqsdd",
	})
	assert.Nil(t, err)
	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:   "remove",
				Path: "/body/2",
			},
		},
		TimestampMillis: now.Add(85 * time.Second).UnixMilli(),
		ClientID:        "50reifj9hyt",
		ChangeID:        "8gdva96nqsdd",
		MessageID:       "50reifj9hyt-8gdva96nqsdd",
	})
	assert.Nil(t, err)

	assert.Nil(t, doc.ReduceHistory(2000))
}

func TestFastForwardChanges(t *testing.T) {
	t.Parallel()

	now := time.Now()
	doc, err := NewDocument([]byte(`{"id":"card1"}`))
	assert.Nil(t, err)

	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/id",
				Value: rawMessage(`"card2"`),
				Prev:  rawMessage(`"card1"`),
			},
		},
		TimestampMillis: now.Add(10 * time.Second).UnixMilli(),
		ClientID:        "50reifj9hyt",
		ChangeID:        "dva96nqsdd",
	})
	assert.Nil(t, err)

	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/id",
				Value: rawMessage(`"card3"`),
				Prev:  rawMessage(`"card2"`),
			},
		},
		TimestampMillis: now.Add(20 * time.Second).UnixMilli(),
		ClientID:        "50reifj9hyt",
		ChangeID:        "gdva96nqsdd",
	})
	assert.Nil(t, err)

	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/id",
				Value: rawMessage(`"card4"`),
				Prev:  rawMessage(`"card3"`),
			},
		},
		TimestampMillis: now.Add(30 * time.Second).UnixMilli(),
		ClientID:        "50reifj9hyt",
		ChangeID:        "hdva96nqsdd",
	})
	assert.Nil(t, err)

	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/id",
				Value: rawMessage(`"card16"`),
				Prev:  rawMessage(`"card1"`),
			},
		},
		TimestampMillis: now.UnixMilli(),
		ClientID:        "50reifj9hyt",
		ChangeID:        "udva96nqsdd",
	})
	assert.Nil(t, err)

	ids := []string{}
	for _, change := range doc.History() {
		ids = append(ids, change.ChangeID)
	}
	assert.Equal(t, []string{"0", "udva96nqsdd", "dva96nqsdd", "gdva96nqsdd", "hdva96nqsdd"}, ids)
	assert.Equal(t, `{"id":"card4"}`, string(doc.JSON()))
}

func TestWrongPrev(t *testing.T) {
	t.Parallel()

	now := time.Now()
	doc, err := NewDocument([]byte(`{"id":"card1"}`))
	assert.Nil(t, err)

	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/id",
				Value: rawMessage(`"card2"`),
				Prev:  rawMessage(`"baa"`),
			},
		},
		TimestampMillis: now.Add(10 * time.Minute).UnixMilli(),
		ClientID:        "50reifj9hyt",
		ChangeID:        "dva96nqsdd",
	})
	assert.Nil(t, err)
	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/id",
				Value: rawMessage(`"card3"`),
				Prev:  rawMessage(`"foo"`),
			},
		},
		TimestampMillis: now.Add(20 * time.Minute).UnixMilli(),
		ClientID:        "50reifj9hyt",
		ChangeID:        "gdva96nqsdd",
	})
	assert.Nil(t, err)
	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/id",
				Value: rawMessage(`"card4"`),
				Prev:  rawMessage(`"other value"`),
			},
		},
		TimestampMillis: now.Add(15 * time.Minute).UnixMilli(),
		ClientID:        "50reifj9hyt",
		ChangeID:        "hdva96nqsdd",
	})
	assert.Nil(t, err)
	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:    "replace",
				Path:  "/id",
				Value: rawMessage(`"card5"`),
				Prev:  rawMessage(`"oh"`),
			},
		},
		TimestampMillis: now.Add(18 * time.Minute).UnixMilli(),
		ClientID:        "50reifj9hyt",
		ChangeID:        "ba96nqsdd",
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
	doc, err := NewDocument([]byte(`{"cards":[{"id":"card1"}]}`))
	assert.Nil(t, err)

	err = doc.ApplyChange(Change{
		Diff: []Operation{
			{
				Op:   "remove",
				Path: "/cards/[card132]",
			},
		},
		TimestampMillis: now.Add(5 * time.Minute).UnixMilli(),
		ClientID:        "50reifj9hyt",
		ChangeID:        "dva96nqsdd",
	})

	assert.Len(t, doc.history, 1)
	assert.Equal(t, "patch error: can't apply changeID dva96nqsdd: id `card132` not found", err.Error())
}

func TestGetValue(t *testing.T) {
	t.Parallel()

	doc, err := NewDocument([]byte(`{"id":"card1", "number": 1234, "boolean": true, "object": {"foo": "bar"},  "array": ["one", "two"], "complex": [{"id":"1234","foo":"baa"}]}`))
	assert.Nil(t, err)

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
	docA, errA := NewDocument(rawA)
	docB, errB := NewDocument(rawB)
	assert.Nil(t, errA)
	assert.Nil(t, errB)

	change, err := docA.Diff(docB)
	change.TimestampMillis = time.Now().UnixMilli()
	change.ClientID = "50reifj9hyt"
	change.ChangeID = "dva96nqsdd"
	assert.Nil(t, err)

	err = docA.ApplyChange(change)
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
