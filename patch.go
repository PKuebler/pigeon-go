package pigeongo

import (
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
	jsonpatch "gopkg.in/evanphx/json-patch.v5"
)

func patch(doc []byte, operations []Operation) ([]byte, error) {
	patchObj := NewJsonpatchPatch(operations)
	patchObj = replacePaths(doc, patchObj)
	newDoc, err := patchObj.Apply(doc)
	if err != nil {
		return doc, err
	}
	return newDoc, nil
}

func replacePaths(doc []byte, patchObj jsonpatch.Patch) jsonpatch.Patch {
	for _, patch := range patchObj {
		path, errPath := patch.Path()
		from, errFrom := patch.From()
		if errPath != nil && errFrom != nil {
			continue
		}

		if path != "" {
			path = replacePath(doc, path)
			patch["path"] = rawMessage(fmt.Sprintf("\"%s\"", path))
		}

		if from != "" {
			from = replacePath(doc, from)
			patch["from"] = rawMessage(fmt.Sprintf("\"%s\"", from))
		}
	}

	return patchObj
}

func replacePath(doc []byte, path string) string {
	parts := strings.Split(path, "/")
	keys := make([]string, len(parts))
	for i, part := range parts {
		if strings.HasPrefix(part, "[") && strings.HasSuffix(part, "]") {
			i := 0
			if _, err := jsonparser.ArrayEach(doc, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				if err != nil {
					panic(err)
				}

				id, _, _, _ := jsonparser.Get(value, "id")
				if string(id) == part[1:len(part)-1] {
					keys[i] = fmt.Sprintf("%d", i)
				}
				i++
			}); err != nil {
				keys[i] = part
			}
		} else {
			keys[i] = part
		}
	}

	return strings.Join(keys, "/")
}

func findID(payload []byte) string {
	id, dataType, _, _ := jsonparser.Get(payload, "id")
	switch dataType {
	case jsonparser.Number, jsonparser.String:
		return string(id)
	case jsonparser.Null, jsonparser.NotExist, jsonparser.Boolean, jsonparser.Object, jsonparser.Array, jsonparser.Unknown:
		return ""
	}

	return ""
}
