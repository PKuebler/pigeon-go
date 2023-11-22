package pigeongo

import "github.com/buger/jsonparser"

func findID(payload []byte, identifiers [][]string) string {
	for _, identifier := range identifiers {
		id, dataType, _, _ := jsonparser.Get(payload, identifier...)

		switch dataType {
		case jsonparser.NotExist:
			continue
		case jsonparser.Number, jsonparser.String:
			return string(id)
		case jsonparser.Null, jsonparser.Boolean, jsonparser.Object, jsonparser.Array, jsonparser.Unknown:
			return ""
		}
	}

	return ""
}
