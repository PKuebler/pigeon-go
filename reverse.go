package pigeongo

import "strings"

func reverse(operations []Operation) []Operation {
	reversedOperations := make([]Operation, len(operations))

	// reverse
	for i := len(operations) - 1; i >= 0; i-- {
		operation := operations[i]

		switch operation.Op {
		case "add":
			operation.Op = "remove"
			id := findID(*operation.Value)
			if id != "" {
				parts := strings.Split(operation.Path, "/")
				parts[len(parts)-1] = "[" + id + "]"
				operation.Path = strings.Join(parts, "/")
			}
		case "remove":
			operation.Op = "add"
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
