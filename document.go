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
		d.history[0].Ts = t.UnixMilli()
	}
}

func WithInitialIDs(cid, gid string) DocumentOption {
	return func(d *Document) {
		d.history[0].Cid = cid
		d.history[0].Gid = gid
	}
}

func WithInitialMid(mid string) DocumentOption {
	return func(d *Document) {
		d.history[0].Mid = mid
	}
}

type Document struct {
	raw         []byte
	history     []Changes
	Warning     string
	gids        map[string]int
	stash       []Changes
	identifiers [][]string
}

func NewDocument(raw []byte, opts ...DocumentOption) *Document {
	doc := &Document{
		raw:         raw,
		gids:        map[string]int{},
		stash:       []Changes{},
		identifiers: [][]string{{"id"}},
	}

	doc.history = []Changes{
		{
			Diff: createInitialDiff(raw),
			Ts:   0,
			Cid:  "0",
			Gid:  "0",
		},
	}

	for _, opt := range opts {
		opt(doc)
	}

	return doc
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

type Changes struct {
	Diff []Operation `json:"diff"`
	Ts   int64       `json:"ts"`
	Cid  string      `json:"cid"`
	Seq  int         `json:"seq"`
	Gid  string      `json:"gid"`
	Mid  string      `json:"mid,omitempty"`
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

func (d *Document) History() []Changes {
	return d.history
}

func (d *Document) Clone() *Document {
	clone := &Document{
		raw:         make([]byte, len(d.raw)),
		history:     make([]Changes, len(d.history)),
		stash:       make([]Changes, len(d.stash)),
		identifiers: make([][]string, len(d.identifiers)),
		gids:        map[string]int{},
	}

	copy(clone.raw, d.raw)
	copy(clone.history, d.history)
	copy(clone.stash, d.stash)

	for a, b := range d.gids {
		clone.gids[a] = b
	}

	for i, identifiers := range d.identifiers {
		clone.identifiers[i] = make([]string, len(identifiers))
		copy(clone.identifiers[i], identifiers)
	}

	return clone
}

func (d *Document) ApplyChanges(changes Changes) {
	// reset warning
	d.Warning = ""

	if _, ok := d.gids[changes.Gid]; ok {
		return
	}

	if err := d.RewindChanges(changes.Ts, changes.Cid); err != nil {
		d.Warning = fmt.Sprintf("rewind error: %s", err.Error())
	}

	// remove external _prev from changes
	// set prev value
	for i := range changes.Diff {
		changes.Diff[i].Prev = d.getValue(changes.Diff[i].Path)
	}

	var err error
	d.raw, err = patch(d.raw, changes.Diff, d.identifiers)
	if err != nil {
		d.Warning = fmt.Sprintf("patch error: %s", err.Error())
	}

	d.gids[changes.Gid] = 1

	if err := d.FastForwardChanges(); err != nil {
		d.Warning = fmt.Sprintf("fast forward error: %s", err.Error())
	}

	idx := len(d.history)
	if idx == 0 {
		d.history = append(d.history, changes)
		return
	}

	// find position to insert
	for idx > 1 && d.history[idx-1].Ts > changes.Ts {
		idx--
	}

	// empty history or after last element
	if len(d.history) == idx {
		d.history = append(d.history, changes)
		return
	}

	d.history = append(d.history[:idx+1], d.history[idx:]...)
	d.history[idx] = changes
}

func (d *Document) FastForwardChanges() error {
	var err error
	for i := len(d.stash) - 1; i >= 0; i-- {
		change := d.stash[i]
		d.raw, err = patch(d.raw, change.Diff, d.identifiers)
		if err != nil {
			return err
		}
		d.gids[change.Gid] = 1
		d.history = append(d.history, change)
	}
	d.stash = []Changes{}

	return nil
}

func (d *Document) RewindChanges(ts int64, cid string) error {
	docJSON := d.raw
	for {
		if len(d.history) <= 1 {
			break
		}

		change := d.history[len(d.history)-1]
		if change.Ts > ts || (change.Ts == ts && change.Cid > cid) {
			// get element and pop from history
			c := d.history[len(d.history)-1]
			d.history = d.history[:len(d.history)-1]

			var err error
			docJSON, err = patch(docJSON, reverse(c.Diff, d.identifiers), d.identifiers)
			if err != nil {
				return err
			}

			delete(d.gids, c.Gid)
			d.stash = append(d.stash, c)
			continue
		}
		break
	}

	d.raw = docJSON

	return nil
}

func (d *Document) ReduceHistory(minTs int64) error {
	if len(d.history) == 1 {
		// only the initial diff, nothing to reduce
		return nil
	}

	// rewind all changes that are newer the minimum timestamp
	if err := d.RewindChanges(minTs, ""); err != nil {
		return err
	}

	// new first diff is the initial diff
	newHistory := []Changes{
		{
			Diff: createInitialDiff(d.raw),
			Ts:   d.history[len(d.history)-1].Ts,
			Cid:  d.history[len(d.history)-1].Cid,
			Gid:  d.history[len(d.history)-1].Gid,
			Mid:  d.history[len(d.history)-1].Mid,
		},
	}

	// reverse append the stash changes
	for i := len(d.stash) - 1; i >= 0; i-- {
		newHistory = append(newHistory, d.stash[i])
	}

	// append all newer changes to history
	if err := d.FastForwardChanges(); err != nil {
		return err
	}

	d.history = newHistory
	return nil
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
