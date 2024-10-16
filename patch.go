package pigeongo

import (
	"errors"
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
		patchObj, err = replacePaths(newDoc, patchObj, identifiers)
		if err != nil {
			return doc, err
		}

		newDoc, err = patchObj.Apply(newDoc)
		if err != nil {
			return doc, err
		}
	}

	return newDoc, nil
}

func replacePaths(doc []byte, patchObj jsonpatch.Patch, identifiers [][]string) (jsonpatch.Patch, error) {
	for _, patch := range patchObj {
		path, errPath := patch.Path()
		from, errFrom := patch.From()
		if errPath != nil && errFrom != nil {
			continue
		}

		if path != "" {
			path, err := replacePath(doc, path, identifiers)
			if err != nil {
				return nil, err
			}
			patch["path"] = rawMessage(fmt.Sprintf("\"%s\"", path))
		}

		if from != "" {
			from, err := replacePath(doc, from, identifiers)
			if err != nil {
				return nil, err
			}
			patch["from"] = rawMessage(fmt.Sprintf("\"%s\"", from))
		}
	}

	return patchObj, nil
}

func replacePath(doc []byte, path string, identifiers [][]string) (string, error) {
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
			found := false
			if _, err := jsonparser.ArrayEach(doc, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				if err != nil {
					panic(err)
				}

				if findID(value, identifiers) == searchID {
					keys = append(keys, fmt.Sprintf("[%d]", childPosition))
					newParts[partIndex] = fmt.Sprintf("%d", childPosition)
					found = true
				}

				childPosition++
			}, keys...); err != nil {
				keys = append(keys, part)
				newParts[partIndex] = part
			}

			if !found {
				return "", errors.New("id `" + searchID + "` not found")
			}
		} else {
			keys = append(keys, part)
			newParts[partIndex] = part
		}
	}

	return strings.Join(newParts, "/"), nil
}
