package pigeongo

import (
	"encoding/json"
	"fmt"
	"slices"
)

func validateDuplicateIdentifiers(doc []byte, identifiers [][]string) error {
	// Validate identifiers
	if len(identifiers) == 0 {
		return nil
	}

	var data any
	if err := json.Unmarshal(doc, &data); err != nil {
		return err
	}

	return walkValidateDuplicateIdentifiers(data, "", identifiers)
}

func walkValidateDuplicateIdentifiers(data any, currentPath string, identifiers [][]string) error {
	if len(identifiers) == 0 {
		return nil
	}

	switch v := data.(type) {
	case map[string]any:
		for key, value := range v {
			if err := walkValidateDuplicateIdentifiers(value, fmt.Sprintf("%s/%s", currentPath, key), identifiers); err != nil {
				return err
			}
		}
	case []any:
		foundIdentifiers := []string{}
		for i, item := range v {
			switch value := item.(type) {
			case map[string]any:
				id := getID(value, identifiers)
				if slices.Contains(foundIdentifiers, id) {
					if id == "" {
						id = "missing id"
					}
					return fmt.Errorf("duplicate identifier found: id `%s` at path %s/%d", id, currentPath, i)
				}
				foundIdentifiers = append(foundIdentifiers, id)

				if err := walkValidateDuplicateIdentifiers(value, fmt.Sprintf("%s/%s", currentPath, id), identifiers); err != nil {
					return err
				}
			case []any:
				if err := walkValidateDuplicateIdentifiers(value, fmt.Sprintf("%s/%d", currentPath, i), identifiers); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
