package pigeongo

import (
	"encoding/json"
	"fmt"

	jsonpatch "gopkg.in/evanphx/json-patch.v5"
)

type Document struct {
	JSON    []byte         `json:"json"`
	History []Changes      `json:"history"`
	Warning string         `json:"warning"`
	Gids    map[string]int `json:"gids"`
	Stash   []Changes      `json:"stash"`
}

func NewDocument(json []byte) *Document {
	return &Document{
		JSON: json,
		History: []Changes{
			{
				Diff: []Operation{
					{
						Op:    "add",
						Path:  "/",
						Value: rawMessage(string(json)),
					},
				},
				Ts:  0,
				Cid: "0",
				Gid: "0",
			},
		},
		Gids:  map[string]int{},
		Stash: []Changes{},
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

func (d *Document) ApplyChanges(changes Changes) {
	if _, ok := d.Gids[changes.Gid]; ok {
		return
	}

	if err := d.RewindChanges(changes.Ts, changes.Cid); err != nil {
		d.Warning = fmt.Sprintf("rewind error: %s", err.Error())
	}

	var err error
	d.JSON, err = patch(d.JSON, changes.Diff)
	if err != nil {
		d.Warning = fmt.Sprintf("patch error: %s", err.Error())
	}

	d.Gids[changes.Gid] = 1

	if err := d.FastForwardChanges(); err != nil {
		d.Warning = fmt.Sprintf("fast forward error: %s", err.Error())
	}

	idx := len(d.History)
	if idx == 0 {
		d.History = append(d.History, changes)
		return
	}

	// find position to insert
	for idx > 1 && d.History[idx-1].Ts > changes.Ts {
		idx--
	}

	// empty history or after last element
	if len(d.History) == idx {
		d.History = append(d.History, changes)
		return
	}

	d.History = append(d.History[:idx+1], d.History[idx:]...)
	d.History[idx] = changes
}

func (d *Document) FastForwardChanges() error {
	var err error
	for _, change := range d.Stash {
		d.JSON, err = patch(d.JSON, change.Diff)
		if err != nil {
			return err
		}
		d.Gids[change.Gid] = 1
		d.History = append(d.History, change)
	}
	d.Stash = []Changes{}

	return nil
}

func (d *Document) RewindChanges(ts int, cid string) error {
	docJSON := d.JSON
	for {
		if len(d.History) <= 1 {
			break
		}

		change := d.History[len(d.History)-1]
		if change.Ts > ts || (change.Ts == ts && change.Cid > cid) {
			// get element and pop from history
			c := d.History[len(d.History)-1]
			d.History = d.History[:len(d.History)-1]

			var err error
			docJSON, err = patch(docJSON, reverse(c.Diff))
			if err != nil {
				return err
			}

			delete(d.Gids, c.Gid)
			d.Stash = append(d.Stash, c)
			continue
		}
		break
	}

	d.JSON = docJSON

	return nil
}

func rawMessage(s string) *json.RawMessage {
	raw := json.RawMessage([]byte(s))
	return &raw
}
