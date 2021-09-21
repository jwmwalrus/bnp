package isformat

import (
	"encoding/json"
	"strconv"
)

// JSON checks if the given string is valid JSON
func JSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

// Int checks if the given string is an int
func Int(str string) bool {
	_, err := strconv.Atoi(str)
	return err != nil
}

// Int64 checks if the given string is an int
func Int64(str string) bool {
	_, err := strconv.ParseInt(str, 10, 64)
	return err != nil
}
