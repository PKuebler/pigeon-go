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
		Change      Change          `json:"change"`
		Result      json.RawMessage `json:"result"`
		Name        string          `json:"name"`
		GoChange    Change          `json:"goChange,omitempty"`
		GoResult    json.RawMessage `json:"goResult"`
	}
	err = json.Unmarshal(file, &testCases)
	assert.Nil(t, err)

	for i, testCase := range testCases {
		fmt.Println("")
		fmt.Println("##### ", testCase.Name)
		oldDocument, err := NewDocument(
			testCase.OldDocument,
			WithIdentifiers([][]string{
				{"attrs", "id"}, {"id"}, {"_id"}, {"uuid"}, {"slug"},
			}),
		)
		assert.Nil(t, err)
		newDocument, err := NewDocument(
			testCase.NewDocument,
			WithIdentifiers([][]string{
				{"attrs", "id"}, {"id"}, {"_id"}, {"uuid"}, {"slug"},
			}),
		)
		assert.Nil(t, err)
		change, err := oldDocument.Diff(newDocument)
		assert.Nil(t, err)

		// format values
		for i, op := range change.Diff {
			change.Diff[i] = formatOperations(op)
		}
		for i, op := range testCase.Change.Diff {
			testCase.Change.Diff[i] = formatOperations(op)
		}

		// set random values, that are static from testfile
		change.TimestampMillis = testCase.Change.TimestampMillis
		change.ClientID = testCase.Change.ClientID
		change.ChangeID = testCase.Change.ChangeID
		change.Seq = testCase.Change.Seq

		testCase.GoChange = testCase.Change

		rawTestChange, _ := json.Marshal(testCase.Change)
		rawDiffChange, _ := json.Marshal(change)

		failed := false

		pigeonJSDoc := oldDocument.Clone()

		// use own change
		fmt.Println("own change!")
		err = oldDocument.ApplyChange(change)
		if !assert.Nil(t, err, "%s - ownChange Warning", testCase.Name) {
			failed = true
		}
		if !assert.Equal(t, string(sortKeys(testCase.NewDocument)), string(sortKeys(oldDocument.JSON())), "testcase %s - ownChange new document", testCase.Name) {
			failed = true
		}
		testCase.GoResult = json.RawMessage(oldDocument.JSON())

		// use pigoenJS change
		fmt.Println("pigoenJS change")
		err = pigeonJSDoc.ApplyChange(testCase.Change)
		if !assert.Nil(t, err, "testcase %d - pigeonChange Warning", i) {
			failed = true
		}
		if !assert.Equal(t, string(sortKeys(testCase.Result)), string(sortKeys(pigeonJSDoc.JSON())), "testcase %d - pigeonChange new document", i) {
			failed = true
		}

		if failed {
			assert.Equal(t, string(rawTestChange), string(rawDiffChange), "testcase %d - change", i)
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
