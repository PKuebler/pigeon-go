package pigeongo

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPigeonJS(t *testing.T) {
	t.Parallel()

	// load testfile
	file, err := os.ReadFile("testfiles/test-cases.json")
	assert.Nil(t, err)

	var testCases []*struct {
		OldDocument json.RawMessage `json:"oldDocument"`
		NewDocument json.RawMessage `json:"newDocument"`
		Changes     Changes         `json:"changes"`
		Result      json.RawMessage `json:"result"`
		Name        string          `json:"name"`
		GoChanges   Changes         `json:"goChanges,omitempty"`
		GoResult    json.RawMessage `json:"goResult"`
	}
	err = json.Unmarshal(file, &testCases)
	assert.Nil(t, err)

	for i, testCase := range testCases {
		fmt.Println("")
		fmt.Println("##### ", testCase.Name)
		oldDocument := NewDocument(
			testCase.OldDocument,
			WithIdentifiers([][]string{
				{"attrs", "id"}, {"id"}, {"_id"}, {"uuid"}, {"slug"},
			}),
		)
		newDocument := NewDocument(
			testCase.NewDocument,
			WithIdentifiers([][]string{
				{"attrs", "id"}, {"id"}, {"_id"}, {"uuid"}, {"slug"},
			}),
		)
		changes, err := oldDocument.Diff(newDocument)
		assert.Nil(t, err)

		// format values
		for i, op := range changes.Diff {
			changes.Diff[i] = formatOperations(op)
		}
		for i, op := range testCase.Changes.Diff {
			testCase.Changes.Diff[i] = formatOperations(op)
		}

		// set random values, that are static from testfile
		changes.Ts = testCase.Changes.Ts
		changes.Cid = testCase.Changes.Cid
		changes.Gid = testCase.Changes.Gid
		changes.Seq = testCase.Changes.Seq

		testCase.GoChanges = testCase.Changes

		rawTestChanges, _ := json.Marshal(testCase.Changes)
		rawDiffChanges, _ := json.Marshal(changes)

		failed := false

		pigeonJSDoc := oldDocument.Clone()

		// use own changes
		fmt.Println("own changes!")
		err = oldDocument.ApplyChanges(changes)
		if !assert.Nil(t, err, "%s - ownChanges Warning", testCase.Name) {
			failed = true
		}
		if !assert.Equal(t, string(sortKeys(testCase.NewDocument)), string(sortKeys(oldDocument.JSON())), "testcase %s - ownChanges new document", testCase.Name) {
			failed = true
		}
		testCase.GoResult = json.RawMessage(oldDocument.JSON())

		// use pigoenJS changes
		fmt.Println("pigoenJS changes")
		err = pigeonJSDoc.ApplyChanges(testCase.Changes)
		if !assert.Nil(t, err, "testcase %d - pigeonChanges Warning", i) {
			failed = true
		}
		if !assert.Equal(t, string(sortKeys(testCase.Result)), string(sortKeys(pigeonJSDoc.JSON())), "testcase %d - pigeonChanges new document", i) {
			failed = true
		}

		if failed {
			assert.Equal(t, string(rawTestChanges), string(rawDiffChanges), "testcase %d - changes", i)
		}
	}

	// write to file
	raw, _ := json.MarshalIndent(testCases, "", "  ")
	_ = os.WriteFile("./testfiles/go-results.json", raw, 0600)
}

func formatOperations(op Operation) Operation {
	if op.Value != nil {
		var data interface{}
		_ = json.Unmarshal(*op.Value, &data)
		var b json.RawMessage
		b, _ = json.MarshalIndent(data, "", "   ")
		op.Value = &b
	}
	if op.Prev != nil {
		var data interface{}
		_ = json.Unmarshal(*op.Prev, &data)
		var b json.RawMessage
		b, _ = json.MarshalIndent(data, "", "   ")
		op.Prev = &b
	}

	return op
}
