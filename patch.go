package pigeongo

import (
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
	jsonpatch "gopkg.in/evanphx/json-patch.v5"
)

func patch(doc []byte, operations []Operation, identifiers [][]string) ([]byte, error) {
	newDoc := doc
	var err error

	for _, operation := range operations {
		patchObj := NewJsonpatchPatch([]Operation{operation})
		patchObj = replacePaths(newDoc, patchObj, identifiers)

		newDoc, err = patchObj.Apply(newDoc)
		if err != nil {
			return doc, err
		}
	}

	return newDoc, nil
}

func replacePaths(doc []byte, patchObj jsonpatch.Patch, identifiers [][]string) jsonpatch.Patch {
	for _, patch := range patchObj {
		path, errPath := patch.Path()
		from, errFrom := patch.From()
		if errPath != nil && errFrom != nil {
			continue
		}

		if path != "" {
			path = replacePath(doc, path, identifiers)
			patch["path"] = rawMessage(fmt.Sprintf("\"%s\"", path))
		}

		if from != "" {
			from = replacePath(doc, from, identifiers)
			patch["from"] = rawMessage(fmt.Sprintf("\"%s\"", from))
		}
	}

	return patchObj
}

func replacePath(doc []byte, path string, identifiers [][]string) string {
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

				if findID(value, identifiers) == searchID {
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
