package pigeongo

import (
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
	jsonpatch "gopkg.in/evanphx/json-patch.v5"
)

func patch(doc []byte, operations []Operation) ([]byte, error) {
	var newDoc []byte
	for _, operation := range operations {
		patchObj := NewJsonpatchPatch([]Operation{operation})
		patchObj = replacePaths(doc, patchObj)

		var err error
		newDoc, err = patchObj.Apply(doc)
		if err != nil {
			return doc, err
		}
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
	newParts := make([]string, len(parts))
	keys := []string{}

	for partIndex, part := range parts {
		if part == "" {
			newParts[partIndex] = part
			continue
		}

		if strings.HasPrefix(part, "[") && strings.HasSuffix(part, "]") {
			searchID := strings.TrimSuffix(strings.TrimPrefix(part, "["), "]")
			childPosition := 0
			if _, err := jsonparser.ArrayEach(doc, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				if err != nil {
					panic(err)
				}

				id, _, _, _ := jsonparser.Get(value, "id")
				if string(id) == searchID {
					key := fmt.Sprintf("%d", childPosition)
					keys = append(keys, key)
					newParts[partIndex] = key
				}

				childPosition++
			}, keys...); err != nil {
				keys = append(keys, part)
				newParts[partIndex] = part
			}
		} else {
			keys = append(keys, part)
			newParts[partIndex] = part
		}
	}

	return strings.Join(newParts, "/")
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
