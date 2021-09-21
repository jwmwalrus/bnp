package slice

import "reflect"

// HasValue checks if value is in the given slice
func HasValue(slic interface{}, val interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(slic).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(slic)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

// Contains checks if value is in the given slice
func Contains(slic interface{}, val interface{}) bool {
	exists, _ := HasValue(slic, val)
	return exists
}

// IndexOf returns the first index found for the given value, or -1
func IndexOf(slic interface{}, val interface{}) int {
	_, index := HasValue(slic, val)

	return index
}
