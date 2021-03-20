package bnp

import "encoding/json"

// IsJSON checks if the given string is valid JSON
func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
