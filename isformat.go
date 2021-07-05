package bnp

import (
	"encoding/json"
	"strconv"
)

// IsJSON checks if the given string is valid JSON
func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

// IsInt checks if the given string is an int
func IsInt(str string) bool {
	_, err := strconv.Atoi(str)
	return err != nil
}

// IsInt64 checks if the given string is an int
func IsInt64(str string) bool {
	_, err := strconv.ParseInt(str, 10, 64)
	return err != nil
}
