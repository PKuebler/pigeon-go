package pigeongo

import "strings"

func reverse(operations []Operation, identifiers [][]string) []Operation {
	reversedOperations := make([]Operation, len(operations))

	// reverse
	for i := len(operations) - 1; i >= 0; i-- {
		operation := operations[i]

		switch operation.Op {
		case "add":
			operation.Op = "remove"
			id := findID(*operation.Value, identifiers)
			if id != "" {
				parts := strings.Split(operation.Path, "/")
				parts[len(parts)-1] = "[" + id + "]"
				operation.Path = strings.Join(parts, "/")
			}
		case "remove":
			operation.Op = "add"

			parts := strings.Split(operation.Path, "/")
			if strings.HasPrefix(parts[len(parts)-1], "[") && strings.HasSuffix(parts[len(parts)-1], "]") {
				parts[len(parts)-1] = "0"
			}
			operation.Path = strings.Join(parts, "/")
		}

		// switch value and prev
		prev := operation.Prev
		value := operation.Value

		if prev == nil {
			operation.Value = nil
		} else {
			operation.Value = prev
		}

		if value == nil {
			operation.Prev = nil
		} else {
			operation.Prev = value
		}

		reversedOperations[len(operations)-1-i] = operation
	}

	return reversedOperations
}
