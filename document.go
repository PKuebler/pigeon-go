package pigeongo

import (
	"encoding/json"
	"fmt"

	jsonpatch "gopkg.in/evanphx/json-patch.v5"
)

type Document struct {
	raw     []byte
	history []Changes
	Warning string
	gids    map[string]int
	stash   []Changes
}

func NewDocument(raw []byte) *Document {
	return &Document{
		raw: raw,
		history: []Changes{
			{
				Diff: []Operation{
					{
						Op:    "add",
						Path:  "/",
						Value: rawMessage(string(raw)),
					},
				},
				Ts:  0,
				Cid: "0",
				Gid: "0",
			},
		},
		gids:  map[string]int{},
		stash: []Changes{},
	}
}

type Changes struct {
	Diff []Operation `json:"diff"`
	Ts   int         `json:"ts"`
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
	d.raw, err = patch(d.raw, changes.Diff)
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
		d.raw, err = patch(d.raw, change.Diff)
		if err != nil {
			return err
		}
		d.gids[change.Gid] = 1
		d.history = append(d.history, change)
	}
	d.stash = []Changes{}

	return nil
}

func (d *Document) RewindChanges(ts int, cid string) error {
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
			docJSON, err = patch(docJSON, reverse(c.Diff))
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

func rawMessage(s string) *json.RawMessage {
	raw := json.RawMessage([]byte(s))
	return &raw
}
