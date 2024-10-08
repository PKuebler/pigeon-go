package pigeongo

import (
	"encoding/json"
	"fmt"
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
		if len(d.history) == 0 {
			return
		}

		d.history[0].Ts = t.UnixMilli()
	}
}

func WithInitialIDs(cid, gid string) DocumentOption {
	return func(d *Document) {
		if len(d.history) == 0 {
			return
		}

		d.history[0].Cid = cid
		d.history[0].Gid = gid
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

func (d *Document) ApplyChanges(changes Changes) {
	// reset warning
	d.Warning = ""

	if _, ok := d.gids[changes.Gid]; ok {
		return
	}

	if err := d.RewindChanges(changes.Ts, changes.Cid); err != nil {
		d.Warning = fmt.Sprintf("rewind error: %s", err.Error())
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
	for _, change := range d.stash {
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

func (d *Document) ReduceHistory(minTs int64) {
	history := []Changes{
		{
			Diff: createInitialDiff(d.raw),
			Ts:   0,
			Cid:  "0",
			Gid:  "0",
		},
	}

	for i, change := range d.history {
		if i == 0 {
			// skip init
			continue
		}

		if change.Ts >= minTs {
			history = append(history, change)
		} else {
			history[0].Ts = change.Ts
			history[0].Cid = change.Cid
			history[0].Gid = change.Gid
		}
	}

	d.history = history
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
