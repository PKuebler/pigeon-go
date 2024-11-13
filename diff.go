package pigeongo

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func diff(left, right []byte, identifiers [][]string) ([]Operation, error) {
	var l interface{}
	var r interface{}

	if err := json.Unmarshal(left, &l); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(right, &r); err != nil {
		return nil, err
	}

	ops := compare([]Operation{}, "", l, r, identifiers)
	return ops, nil
}

func compare(ops []Operation, path string, left, right interface{}, identifiers [][]string) []Operation {
	// if left and right nil, no changes
	if left == nil && right == nil {
		return nil
	}

	// if only one of them is nil, it is a add or remove
	if left == nil || right == nil {
		return append(ops, newChange(path, left, right))
	}

	leftKind := reflect.TypeOf(left).Kind()
	rightKind := reflect.TypeOf(right).Kind()

	// type mismatch is a change
	if leftKind != rightKind {
		return append(ops, newChange(path, left, right))
	}

	switch leftKind {
	case reflect.Map:
		// compare object keys
		return compareMaps(ops, path, left.(map[string]interface{}), right.(map[string]interface{}), identifiers)
	case reflect.Slice:
		// compare array values
		return compareSlices(ops, path, left.([]interface{}), right.([]interface{}), identifiers)
	default:
		// compare primitive values
		if !reflect.DeepEqual(left, right) {
			return append(ops, newChange(path, left, right))
		}
	}

	return ops
}

func compareMaps(ops []Operation, path string, left, right map[string]interface{}, identifiers [][]string) []Operation {
	// sort keys cosmetically only, as PigeonJS uses an alphabetical order.
	leftKeys := make([]string, 0, len(left))
	for key := range left {
		leftKeys = append(leftKeys, key)
	}
	sort.SliceStable(leftKeys, func(i int, j int) bool { return leftKeys[i] < leftKeys[j] })

	rightKeys := make([]string, 0, len(right))
	for key := range right {
		rightKeys = append(rightKeys, key)
	}
	sort.SliceStable(rightKeys, func(i int, j int) bool { return rightKeys[i] < rightKeys[j] })

	// iterates over the left object
	for _, key := range leftKeys {
		leftVal := left[key]
		rightVal, exists := right[key]
		newPath := path + "/" + key
		if exists {
			// compare values if the key exists in both objects
			ops = compare(ops, newPath, leftVal, rightVal, identifiers)
		} else {
			// key exists in the left object but not in the right one (removed to the right)
			ops = append(ops, newChange(newPath, leftVal, nil))
		}
	}

	// iterates over the right object, looking for keys that only exist in the right object
	for _, key := range rightKeys {
		if _, exists := left[key]; !exists {
			rightVal := right[key]
			newPath := path + "/" + key
			ops = append(ops, newChange(newPath, nil, rightVal))
		}
	}

	return ops
}

func isSlicePrimitive(slice []interface{}) bool {
	for _, item := range slice {
		switch item.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string, bool:
		default:
			return false
		}
	}

	return true
}

func comparePrimitiveSlices(ops []Operation, path string, left, right []interface{}) []Operation {
	if !isSlicePrimitive(right) || len(left) != len(right) {
		// replace all
		ops = append(ops, newChange(path, left, right))
		return ops
	}

	// compare values
	for i, val := range left {
		if val != right[i] {
			// stop an replace all
			ops = append(ops, newChange(path, left, right))
			return ops
		}
	}

	return ops
}

func compareSlices(ops []Operation, path string, left, right []interface{}, identifiers [][]string) []Operation {
	// is slice primitive, use only replace all operations
	if isSlicePrimitive(left) {
		ops = comparePrimitiveSlices(ops, path, left, right)
		return ops
	}

	// to compare id's we need to find the index of the object in the right slice by its id.
	rightIDIndexMap := map[string]int{}
	for index, rightVal := range right {
		switch value := rightVal.(type) {
		case map[string]interface{}:
			id := getID(value, identifiers)
			if id != "" {
				rightIDIndexMap[id] = index
			}
		default:
			continue
		}
	}

	handledRight := map[int]bool{}
	for leftIndex, leftVal := range left {
		switch leftVal.(type) {
		case map[string]interface{}:
			id := getID(leftVal, identifiers)
			if rightIndex, exists := rightIDIndexMap[id]; exists {
				// moved?
				if leftIndex != rightIndex {
					oldPath := path + "/" + getArrayItemID(leftVal, right[rightIndex], leftIndex, identifiers)
					newPath := fmt.Sprintf("%s/%d", path, rightIndex)
					ops = append(ops, addMove(oldPath, newPath))
				}
				handledRight[rightIndex] = true
				newPath := path + "/" + getArrayItemID(leftVal, right[rightIndex], leftIndex, identifiers)
				ops = compare(ops, newPath, leftVal, right[rightIndex], identifiers)
			} else {
				// remove
				newPath := path + "/" + getArrayItemID(leftVal, nil, leftIndex, identifiers)
				ops = compare(ops, newPath, leftVal, nil, identifiers)
			}
		}
	}

	for rightIndex, rightVal := range right {
		if _, handled := handledRight[rightIndex]; !handled {
			newPath := fmt.Sprintf("%s/%d", path, rightIndex)
			switch rightVal.(type) {
			case map[string]interface{}:
				op := newChange(newPath, nil, rightVal)

				// support add before id like `/array/[id]`
				if rightIndex < len(right)-1 {
					nextID := getID(right[rightIndex+1], identifiers)
					// use the id at the path, if it exists
					if nextID != "" {
						if _, handledNextID := handledRight[rightIndex+1]; handledNextID {
							pathParts := strings.Split(newPath, "/")
							op.Path = fmt.Sprintf("%s/%s", strings.Join(pathParts[:len(pathParts)-1], "/"), nextID)
						}
					}
				}

				ops = append(ops, op)
			}
		}
	}

	return ops
}

func getID(value interface{}, identifiers [][]string) string {
	if m, ok := value.(map[string]interface{}); ok {
		for _, identifier := range identifiers {
			layer := m
			for i, key := range identifier {
				value, exists := layer[key]
				if !exists {
					break
				}

				// last element
				if i == len(identifier)-1 {
					switch v := value.(type) {
					case float64:
						// is a int?
						if v == float64(int64(v)) {
							return formatID(fmt.Sprintf("%d", int64(v)))
						}
						return formatID(fmt.Sprintf("%f", v))
					case string:
						return formatID(v)
					default:
						return ""
					}
				}

				if m2, ok := value.(map[string]interface{}); !ok {
					// child is not a object
					break
				} else {
					layer = m2
				}
			}
		}
	}
	return ""
}

// getArrayItemID returns the ID of an array item.
func getArrayItemID(left, right interface{}, index int, identifiers [][]string) string {
	// check if the left object has an ID.
	if leftMap, ok := left.(map[string]interface{}); ok {
		id := getID(leftMap, identifiers)
		if id != "" {
			return id
		}
	}

	// check if the right object has an ID.
	if rightMap, ok := right.(map[string]interface{}); ok {
		id := getID(rightMap, identifiers)
		if id != "" {
			return id
		}
	}

	// use the index as ID.
	return strconv.Itoa(index)
}

// formatID formats an ID to a string with [<id>].
func formatID(id interface{}) string {
	switch v := id.(type) {
	case string:
		return "[" + v + "]"
	case float64:
		return "[" + strconv.FormatFloat(v, 'f', -1, 64) + "]"
	default:
		return "[unknown]"
	}
}

func newChange(path string, left, right interface{}) Operation {
	var prevValue, newValue *json.RawMessage
	var operation string

	if left != nil {
		leftBytes, _ := json.Marshal(left)
		prevValue = (*json.RawMessage)(&leftBytes)
	}

	if right != nil {
		rightBytes, _ := json.Marshal(right)
		newValue = (*json.RawMessage)(&rightBytes)
	}

	if left == nil {
		operation = "add"
	} else if right == nil {
		operation = "remove"
	} else {
		operation = "replace"
	}

	return Operation{
		Op:    operation,
		Path:  path,
		Prev:  prevValue,
		Value: newValue,
	}
}

func addMove(fromPath, toPath string) Operation {
	return Operation{
		Op:   "move",
		Path: toPath,
		From: fromPath,
	}
}
