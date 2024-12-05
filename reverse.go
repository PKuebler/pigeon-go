package pigeongo

import (
	"strconv"
	"strings"
)

func reverse(operations []Operation, identifiers [][]string) []Operation {
	reversedOperations := make([]Operation, len(operations))

	// reverse
	for i := len(operations) - 1; i >= 0; i-- {
		operation := operations[i]

		switch operation.Op {
		case "add":
			operation.Op = "remove"
			operation.Prev = nil

			// if value is a object
			id := findID(*operation.Value, identifiers)
			if id != "" {
				parts := strings.Split(operation.Path, "/")
				// if last part is an index position
				if _, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
					// replace /array/0 with /array/[objId]
					parts[len(parts)-1] = "[" + id + "]"
					operation.Path = strings.Join(parts, "/")
				} else if strings.HasPrefix(parts[len(parts)-1], "[") && strings.HasSuffix(parts[len(parts)-1], "]") {
					// replace /array/[insertBeforeThisId] with /array/[objId]
					parts[len(parts)-1] = "[" + id + "]"
					operation.Path = strings.Join(parts, "/")
				}
			}
		case "remove":
			operation.Op = "add"
			operation.Value = nil

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
