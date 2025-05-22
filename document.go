package pigeongo

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	jsonpatch "gopkg.in/evanphx/json-patch.v5"
)

type DocumentOption func(*Document)

func WithIdentifiers(identifiers [][]string) DocumentOption {
	return func(d *Document) {
		d.identifiers = identifiers
	}
}

func WithInitialTime(t time.Time) DocumentOption {
	return func(d *Document) {
		d.history[0].TimestampMillis = t.UnixMilli()
	}
}

func WithInitialIDs(clientID, changeID string) DocumentOption {
	return func(d *Document) {
		d.history[0].ClientID = clientID
		d.history[0].ChangeID = changeID
	}
}

func WithInitialMessageID(messageID string) DocumentOption {
	return func(d *Document) {
		d.history[0].MessageID = messageID
	}
}

type Document struct {
	raw         []byte
	history     []Change
	changeIDs   map[string]int
	stash       []Change
	identifiers [][]string
}

func NewDocument(raw []byte, opts ...DocumentOption) (*Document, error) {
	doc := &Document{
		raw:         raw,
		changeIDs:   map[string]int{},
		stash:       []Change{},
		identifiers: [][]string{{"id"}},
	}

	doc.history = []Change{
		{
			Diff:            createInitialDiff(raw),
			TimestampMillis: 0,
			ClientID:        "0",
			ChangeID:        "0",
		},
	}

	for _, opt := range opts {
		opt(doc)
	}

	if err := validateDuplicateIdentifiers(doc.raw, doc.identifiers); err != nil {
		return nil, fmt.Errorf("error in identifiers: %s", err.Error())
	}

	return doc, nil
}

func createInitialDiff(raw []byte) []Operation {
	diff := []Operation{}

	_, dataType, _, _ := jsonparser.Get(raw)
	switch dataType {
	case jsonparser.Array:
		i := 0
		_, _ = jsonparser.ArrayEach(raw, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			diff = append(diff, Operation{
				Op:    "add",
				Path:  fmt.Sprintf("/%d", i),
				Value: rawToJSON(value, dataType),
			})
			i++
		})
	case jsonparser.Object:
		_ = jsonparser.ObjectEach(raw, func(key, value []byte, dataType jsonparser.ValueType, offset int) error {
			diff = append(diff, Operation{
				Op:    "add",
				Path:  fmt.Sprintf("/%s", key),
				Value: rawToJSON(value, dataType),
			})
			return nil
		})
	default:
		diff = append(diff, Operation{
			Op:    "add",
			Path:  "/",
			Value: rawToJSON(raw, dataType),
		})
	}

	return diff
}

type Change struct {
	Diff            []Operation `json:"diff"`
	TimestampMillis int64       `json:"timestamp_ms"`
	ClientID        string      `json:"client_id"`
	Seq             int         `json:"seq"`
	ChangeID        string      `json:"change_id"`
	MessageID       string      `json:"msg_id,omitempty"`
}

func NewJsonpatchPatch(diff []Operation) jsonpatch.Patch {
	patch, _ := json.Marshal(diff)
	patchObj, _ := jsonpatch.DecodePatch(patch)
	return patchObj
}

type Operation struct {
	Op    string           `json:"op"`
	Path  string           `json:"path,omitempty"`
	From  string           `json:"from,omitempty"`
	Value *json.RawMessage `json:"value,omitempty"`
	Prev  *json.RawMessage `json:"_prev,omitempty"`
}

func (d *Document) JSON() []byte {
	return d.raw
}

func (d *Document) History() []Change {
	return d.history
}

func (d *Document) Clone() *Document {
	clone := &Document{
		raw:         make([]byte, len(d.raw)),
		history:     make([]Change, len(d.history)),
		stash:       make([]Change, len(d.stash)),
		identifiers: make([][]string, len(d.identifiers)),
		changeIDs:   map[string]int{},
	}

	copy(clone.raw, d.raw)
	copy(clone.history, d.history)
	copy(clone.stash, d.stash)

	for a, b := range d.changeIDs {
		clone.changeIDs[a] = b
	}

	for i, identifiers := range d.identifiers {
		clone.identifiers[i] = make([]string, len(identifiers))
		copy(clone.identifiers[i], identifiers)
	}

	return clone
}

// replaceByWorkingCopy overwrite all fields in this document (without identifiers)!
func (d *Document) replaceByWorkingCopy(workingCopy *Document) {
	d.raw = workingCopy.raw
	d.history = workingCopy.history
	d.stash = workingCopy.stash
	d.changeIDs = workingCopy.changeIDs
}

// FastForwardChanges apply all changes from the stash. It change nothing, if one change failed.
func (d *Document) FastForwardChanges() error {
	workingCopy := d.Clone()

	if err := workingCopy.fastForwardChanges(); err != nil {
		return err
	}

	d.replaceByWorkingCopy(workingCopy)
	return nil
}

// RewindChanges rewind to a specific change. It change nothing, if one change failed.
func (d *Document) RewindChanges(timestampMillis int64, clientID string) error {
	workingCopy := d.Clone()

	if err := workingCopy.rewindChanges(timestampMillis, clientID); err != nil {
		return err
	}

	d.replaceByWorkingCopy(workingCopy)
	return nil
}

// ApplyChange to the document. It change nothing, if one operation failed.
func (d *Document) ApplyChange(change Change) error {
	// skip change if changeID is processed
	if _, ok := d.changeIDs[change.ChangeID]; ok {
		return nil
	}

	workingCopy := d.Clone()

	if err := workingCopy.rewindChanges(change.TimestampMillis, change.ClientID); err != nil {
		return fmt.Errorf("patch error for changeID %s: %s", change.ChangeID, err)
	}

	// remove external _prev from change
	// set prev value
	for i := range change.Diff {
		if change.Diff[i].Op == "add" {
			change.Diff[i].Prev = nil
		} else {
			change.Diff[i].Prev = workingCopy.getValue(change.Diff[i].Path)
		}

		if change.Diff[i].Op == "remove" {
			change.Diff[i].Value = nil
		}
	}

	// apply
	var err error
	workingCopy.raw, err = patch(workingCopy.raw, change.Diff, workingCopy.identifiers)
	if err != nil {
		return fmt.Errorf("patch error: can't apply changeID %s: %s", change.ChangeID, err.Error())
	}

	if err := validateDuplicateIdentifiers(workingCopy.raw, workingCopy.identifiers); err != nil {
		return fmt.Errorf("patch error: can't apply changeID %s: %s", change.ChangeID, err.Error())
	}

	workingCopy.changeIDs[change.ChangeID] = 1

	if err := workingCopy.fastForwardChanges(); err != nil {
		return fmt.Errorf("patch error for changeID %s: %s", change.ChangeID, err)
	}

	idx := len(workingCopy.history)
	if idx == 0 {
		workingCopy.history = append(workingCopy.history, change)
		d.replaceByWorkingCopy(workingCopy)
		return nil
	}

	// find position to insert
	for idx > 1 && workingCopy.history[idx-1].TimestampMillis > change.TimestampMillis {
		idx--
	}

	// empty history or after last element
	if len(workingCopy.history) == idx {
		workingCopy.history = append(workingCopy.history, change)
		d.replaceByWorkingCopy(workingCopy)
		return nil
	}

	workingCopy.history = append(workingCopy.history[:idx+1], workingCopy.history[idx:]...)
	workingCopy.history[idx] = change

	d.replaceByWorkingCopy(workingCopy)
	return nil
}

func (d *Document) ReduceHistory(minTimestampMillis int64) error {
	if len(d.history) == 1 {
		// only the initial diff, nothing to reduce
		return nil
	}

	workingCopy := d.Clone()

	// rewind all changes that are newer the minimum timestamp
	if err := workingCopy.rewindChanges(minTimestampMillis, ""); err != nil {
		return err
	}

	// new first diff is the initial diff
	newHistory := []Change{
		{
			Diff:            createInitialDiff(workingCopy.raw),
			TimestampMillis: workingCopy.history[len(workingCopy.history)-1].TimestampMillis,
			ClientID:        workingCopy.history[len(workingCopy.history)-1].ClientID,
			ChangeID:        workingCopy.history[len(workingCopy.history)-1].ChangeID,
			MessageID:       workingCopy.history[len(workingCopy.history)-1].MessageID,
		},
	}

	// reverse append the stash changes
	for i := len(workingCopy.stash) - 1; i >= 0; i-- {
		newHistory = append(newHistory, workingCopy.stash[i])
	}

	// append all newer changes to history
	if err := workingCopy.fastForwardChanges(); err != nil {
		return err
	}

	workingCopy.history = newHistory

	d.replaceByWorkingCopy(workingCopy)
	return nil
}

// fastForwardChanges will apply all changes in the stash. It will stop if a patch fails and reset nothing!
func (d *Document) fastForwardChanges() error {
	var err error
	for i := len(d.stash) - 1; i >= 0; i-- {
		change := d.stash[i]

		// set prev value, maybe changed by patch before!
		for i := range change.Diff {
			if change.Diff[i].Op == "add" {
				change.Diff[i].Prev = nil
			} else {
				change.Diff[i].Prev = d.getValue(change.Diff[i].Path)
			}

			if change.Diff[i].Op == "remove" {
				change.Diff[i].Value = nil
			}
		}

		d.raw, err = patch(d.raw, change.Diff, d.identifiers)
		if err != nil {
			return fmt.Errorf("fast forward error: can't patch changeID %s from stash: %s", change.ChangeID, err.Error())
		}

		d.changeIDs[change.ChangeID] = 1
		d.history = append(d.history, change)
	}
	d.stash = []Change{}

	return nil
}

// rewindChanges will rewind all changes in the history. It will stop if a patch fails and reset nothing!
func (d *Document) rewindChanges(timestampMillis int64, clientID string) error {
	for len(d.history) > 1 {
		change := d.history[len(d.history)-1]
		if change.TimestampMillis > timestampMillis || (change.TimestampMillis == timestampMillis && change.ClientID > clientID) {
			// get element and pop from history
			c := d.history[len(d.history)-1]
			d.history = d.history[:len(d.history)-1]

			var err error
			d.raw, err = patch(d.raw, reverse(c.Diff, d.identifiers), d.identifiers)
			if err != nil {
				return fmt.Errorf("rewind error: can't reverse patch changeID %s from history: %s", change.ChangeID, err.Error())
			}

			delete(d.changeIDs, c.ChangeID)
			d.stash = append(d.stash, c)
			continue
		}
		break
	}

	return nil
}

func (d *Document) Diff(right *Document) (Change, error) {
	operations, err := diff(d.JSON(), right.JSON(), d.identifiers)
	if err != nil {
		return Change{}, err
	}

	return Change{
		Diff: operations,
	}, nil
}

func rawToJSON(value []byte, dataType jsonparser.ValueType) *json.RawMessage {
	switch dataType {
	case jsonparser.String:
		return rawMessage(fmt.Sprintf(`"%s"`, string(value)))
	case jsonparser.Number:
		return rawMessage(string(value))
	case jsonparser.Object:
		return rawMessage(string(value))
	case jsonparser.Array:
		return rawMessage(string(value))
	case jsonparser.Boolean:
		return rawMessage(string(value))
	case jsonparser.Null:
		return rawMessage(string("null"))
	case jsonparser.Unknown:
		return rawMessage(string(value))
	case jsonparser.NotExist:
		return rawMessage(string("null"))
	}
	return nil
}

func rawMessage(s string) *json.RawMessage {
	raw := json.RawMessage([]byte(s))
	return &raw
}

func (d *Document) getValue(path string) *json.RawMessage {
	if d.raw == nil {
		return nil
	}

	jsonpatchPath, err := replacePath(d.raw, path, d.identifiers)
	if err != nil {
		return nil
	}

	jsonparserPath := []string{}

	pathParts := strings.Split(jsonpatchPath, "/")
	if len(pathParts) < 2 || pathParts[0] != "" {
		return nil
	}

	for _, part := range pathParts {
		if part == "" {
			continue
		}

		if _, err := strconv.Atoi(part); err == nil {
			jsonparserPath = append(jsonparserPath, fmt.Sprintf(`[%s]`, part))
			continue
		}

		jsonparserPath = append(jsonparserPath, part)
	}

	value, dataType, _, err := jsonparser.Get(d.raw, jsonparserPath...)
	if err != nil {
		return nil
	}

	return rawToJSON(value, dataType)
}
