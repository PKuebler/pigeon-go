package pigeongo

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

func diff(left, right []byte, identifier [][]string) ([]Operation, error) {
	var l interface{}
	var r interface{}

	if err := json.Unmarshal(left, &l); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(right, &r); err != nil {
		return nil, err
	}

	ops := compare([]Operation{}, "", l, r)
	return ops, nil
}

func compare(ops []Operation, path string, left, right interface{}) []Operation {
	// Wenn beide nil sind, gibt es keine Änderung
	if left == nil && right == nil {
		return nil
	}

	// Wenn eines nil ist, gibt es eine Änderung
	if left == nil || right == nil {
		return append(ops, newChange(path, left, right))
	}

	leftKind := reflect.TypeOf(left).Kind()
	rightKind := reflect.TypeOf(right).Kind()

	// Unterschiedliche Typen - das ist eine Änderung
	if leftKind != rightKind {
		return append(ops, newChange(path, left, right))
	}

	switch leftKind {
	case reflect.Map:
		// Beide sind Maps, vergleichen wir die Keys
		return compareMaps(ops, path, left.(map[string]interface{}), right.(map[string]interface{}))
	case reflect.Slice:
		// Beide sind Arrays, vergleichen wir die Elemente
		return compareSlices(ops, path, left.([]interface{}), right.([]interface{}))
	default:
		// Primitive Werte vergleichen
		if !reflect.DeepEqual(left, right) {
			return append(ops, newChange(path, left, right))
		}
	}

	return ops
}

func compareMaps(ops []Operation, path string, left, right map[string]interface{}) []Operation {
	// Iteriere über das linke Objekt
	for key, leftVal := range left {
		rightVal, exists := right[key]
		newPath := path + "/" + key
		if exists {
			// Vergleiche die Werte, wenn der Key in beiden Objekten existiert
			ops = compare(ops, newPath, leftVal, rightVal)
		} else {
			// Key existiert nur im linken Objekt (gelöscht im rechten)
			ops = append(ops, newChange(newPath, leftVal, nil))
		}
	}

	// Iteriere über das rechte Objekt, um Keys zu finden, die nur im rechten Objekt existieren
	for key, rightVal := range right {
		if _, exists := left[key]; !exists {
			newPath := path + "/" + key
			ops = append(ops, newChange(newPath, nil, rightVal))
		}
	}

	return ops
}

func compareSlices(ops []Operation, path string, left, right []interface{}) []Operation {
	maxLen := len(left)
	if len(right) > maxLen {
		maxLen = len(right)
	}

	rightIndexMap := map[string]int{}
	for index, rightVal := range right {
		obj, ok := rightVal.(map[string]interface{})
		if !ok {
			continue
		}

		id := getID(obj)
		if id != "" {
			rightIndexMap[id] = index
		}
	}

	handledRight := map[int]bool{}
	for leftIndex, leftVal := range left {
		id := getID(leftVal)
		if rightIndex, exists := rightIndexMap[id]; exists {
			if leftIndex != rightIndex {
				oldPath := path + "/" + getArrayItemID(leftVal, right[rightIndex], leftIndex)
				newPath := fmt.Sprintf("%s/%d", path, rightIndex)
				ops = append(ops, addMove(oldPath, newPath))
			}
			handledRight[rightIndex] = true
			newPath := path + "/" + getArrayItemID(leftVal, right[rightIndex], leftIndex)
			ops = compare(ops, newPath, leftVal, right[rightIndex])
		} else {
			newPath := path + "/" + getArrayItemID(leftVal, nil, leftIndex)
			ops = compare(ops, newPath, leftVal, nil)
		}
	}

	adds := []Operation{}
	for rightIndex, rightVal := range right {
		if _, handled := handledRight[rightIndex]; !handled {
			newPath := fmt.Sprintf("%s/%d", path, rightIndex)
			adds = compare(adds, newPath, nil, rightVal)
		}
	}

	ops = append(adds, ops...)

	return ops
}

func getID(value interface{}) string {
	if m, ok := value.(map[string]interface{}); ok {
		if id, exists := m["id"]; exists {
			return formatID(id)
		}
	}
	return ""
}

// Funktion zum Ermitteln der ID eines Array-Elements
func getArrayItemID(left, right interface{}, index int) string {
	// Überprüfen, ob das linke Element eine ID hat
	if leftMap, ok := left.(map[string]interface{}); ok {
		if id, exists := leftMap["id"]; exists {
			return formatID(id)
		}
	}

	// Überprüfen, ob das rechte Element eine ID hat
	if rightMap, ok := right.(map[string]interface{}); ok {
		if id, exists := rightMap["id"]; exists {
			return formatID(id)
		}
	}

	// Wenn keine ID vorhanden ist, den Index verwenden
	return strconv.Itoa(index)
}

// Hilfsfunktion zum Formatieren der ID (mit eckigen Klammern)
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
