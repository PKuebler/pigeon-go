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

	return walkValidateDuplicateIdentifiers(data, identifiers)
}

func walkValidateDuplicateIdentifiers(data any, identifiers [][]string) error {
	if len(identifiers) == 0 {
		return nil
	}

	switch v := data.(type) {
	case map[string]any:
		for _, value := range v {
			if err := walkValidateDuplicateIdentifiers(value, identifiers); err != nil {
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
					return fmt.Errorf("duplicate identifier found: %s at index %d", id, i)
				}
				foundIdentifiers = append(foundIdentifiers, id)

				if err := walkValidateDuplicateIdentifiers(value, identifiers); err != nil {
					return err
				}
			case []any:
				if err := walkValidateDuplicateIdentifiers(value, identifiers); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
